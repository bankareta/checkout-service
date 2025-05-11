package usecase

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/helpers"
	master_discount "checkout-service/internal/master_discount"
	master_products "checkout-service/internal/master_products"
	"checkout-service/internal/model"
	tr_transaction "checkout-service/internal/tr_transaction"
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type transactionUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	TransactionRepository    tr_transaction.TransactionRepository
	MasterProductsRepository master_products.MasterProductsRepository
	MasterDiscountRepository master_discount.MasterDiscountRepository
	mapper                   tr_transaction.TransactionMapper
}

func NewTransactionUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	transactionRepository tr_transaction.TransactionRepository,
	masterProductsRepository master_products.MasterProductsRepository,
	masterDiscountRepository master_discount.MasterDiscountRepository,
	mapper tr_transaction.TransactionMapper) tr_transaction.TransactionUseCase {
	return &transactionUseCase{
		DB:                       db,
		Log:                      logger,
		Validate:                 validate,
		TransactionRepository:    transactionRepository,
		MasterProductsRepository: masterProductsRepository,
		MasterDiscountRepository: masterDiscountRepository,
		mapper:                   mapper,
	}
}

func (c transactionUseCase) ScanProduct(ctx context.Context, request model.ScanProductRequest) (model.ScanProductResponse, error) {
	var res model.ScanProductResponse

	resProduct, err := c.MasterProductsRepository.ScanProduct(request.SKU)
	if err != nil {
		c.Log.Warnf("Failed Get Product : %+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, helpers.NewError(404, nil, "01", "Failed get Product, try again.")
		}
		return res, helpers.NewError(500, nil, "500", err.Error())
	}

	respMap := c.mapper.MapScanProductsResp(res, *resProduct)
	return respMap, nil
}

func (c transactionUseCase) Checkout(ctx context.Context, request model.CheckoutRequest) (model.CheckoutResponse, error) {
	var res model.CheckoutResponse
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err.Error())
		return res, helpers.ValidateError
	}

	freeProductIDs, err := c.MasterDiscountRepository.GetFreeProductIDs()
	if err != nil {
		c.Log.Errorf("failed to fetch free product discounts: %v", err)
		return res, helpers.NewError(500, nil, "01", "Failed checkout, try again.")
	}

	// Remap request items
	request.Items = adjustItemsForFreeProduct(request.Items, freeProductIDs)

	if len(request.Items) > 0 {
		resRaw, err := c.TransactionRepository.ExecuteTransactionWithResult(ctx, func(tx *gorm.DB) (interface{}, error) {
			return c.ProcessCheckout(tx, request)
		})
		if err != nil {
			return res, err
		}
		resp, ok := resRaw.(*model.CheckoutResponse)
		if !ok {
			return res, errors.New("invalid type assertion")
		}
		return *resp, nil
	}
	// return if request len item
	return res, helpers.NewError(404, nil, "01", "Failed checkout, try again.")
}

func adjustItemsForFreeProduct(items []model.CheckoutItem, freeProductIDs map[int]bool) []model.CheckoutItem {
	var adjusted []model.CheckoutItem

	for _, item := range items {
		if freeProductIDs[item.ProductID] {
			if item.Qty > 1 {
				item.Qty -= 1
				adjusted = append(adjusted, item)
			}
		} else {
			adjusted = append(adjusted, item)
		}
	}
	return adjusted
}

func (c transactionUseCase) ProcessCheckout(tx *gorm.DB, request model.CheckoutRequest) (*model.CheckoutResponse, error) {
	var result model.CheckoutResponse

	// declare first id transaction
	trans, err := c.TransactionRepository.AddTransaction(tx, entity.Transaction{
		Status: 0,
	})
	if err != nil {
		c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
			return &result, helpers.NewError(400, nil, "01", "Failed checkout, try again.")
		}
		return &result, helpers.NewError(500, nil, "500", err.Error())
	}

	// set total all
	var totalAllPrice float64
	itemResp := make([]model.CheckoutItemResponse, 0)

	for _, t := range request.Items {
		var totalPrice float64

		// search products
		resProduct, err := c.MasterProductsRepository.DetailProduct(uint(t.ProductID))
		if err != nil {
			c.Log.Warnf("Failed checkout : %+v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &result, helpers.NewError(404, nil, "01", "Failed checkout, try again.")
			}
			return &result, helpers.NewError(500, nil, "500", err.Error())
		}

		// check discount products
		resDiscount, err := c.MasterDiscountRepository.GetDiscountProduct(t.ProductID)
		if err != nil {
			c.Log.Warnf("Failed checkout : %+v", err)
			return &result, helpers.NewError(500, nil, "500", err.Error())
		}

		// process discount if any
		if len(resDiscount) > 0 {
			// Process the discounts
			err = c.processDiscounts(tx, resDiscount, *resProduct, t, trans, &totalPrice, &itemResp)
			if err != nil {
				return &result, err
			}
		} else {
			// No discount, just add the product
			err = c.addProductToTransaction(tx, *resProduct, t.Qty, trans, &totalPrice, &itemResp)
			if err != nil {
				return &result, err
			}
		}

		// accumulate total
		totalAllPrice += totalPrice
	}

	// mapping result
	result.ID = int(trans.ID)
	result.CheckoutDate = trans.CreatedAt
	result.PriceTotal = c.mapper.RoundFloat(totalAllPrice, 2)
	result.TotalFormatted = c.mapper.FormatPrice(totalAllPrice)
	result.Items = itemResp

	return &result, nil
}

func (c transactionUseCase) processDiscounts(tx *gorm.DB, discounts []entity.ProductsDiscount, product entity.Products, item model.CheckoutItem, trans entity.Transaction, totalPrice *float64, itemResp *[]model.CheckoutItemResponse) error {
	for _, discount := range discounts {
		resDetailDiscount, err := c.MasterDiscountRepository.GetDiscount(discount.DiscountsId)
		if err != nil {
			c.Log.Warnf("Failed checkout : %+v", err)
			return helpers.NewError(500, nil, "500", err.Error())
		}

		// Handle specific discount types
		switch resDetailDiscount.Type {
		case 1:
			// Free product
			err := c.discFreeProducts(tx, *resDetailDiscount, product, item, trans, totalPrice, itemResp)
			if err != nil {
				return err
			}
		case 2:
			// Buy X get Y
			err := c.discBuyXGetY(tx, *resDetailDiscount, product, item, trans, totalPrice, itemResp)
			if err != nil {
				return err
			}
		case 3:
			// Percentage or nominal discount
			err := c.discPercentAmount(tx, *resDetailDiscount, product, item, trans, totalPrice, itemResp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c transactionUseCase) discFreeProducts(tx *gorm.DB, resDetailDiscount entity.Discount, resProduct entity.Products, item model.CheckoutItem, trans entity.Transaction, totalPrice *float64, itemResp *[]model.CheckoutItemResponse) error {
	var productResp []model.ProductCheckoutResponse
	var discountDescriptions []string
	subTotal := resProduct.Price * float64(item.Qty)
	// set first product
	err := c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
		TransactionID: int(trans.ID),
		ProductID:     int(resProduct.ID),
		ProductName:   &resProduct.Name,
		SKU:           &resProduct.SKU,
		Qty:           item.Qty,
		Price:         resProduct.Price,
		Status:        0,
		TotalPrice:    subTotal,
	})
	if err != nil {
		c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
			return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
		}
		return helpers.NewError(500, nil, "500", err.Error())
	}

	productResp = append(productResp, model.ProductCheckoutResponse{
		ProductID:      int(resProduct.ID),
		ProductName:    resProduct.Name,
		Qty:            item.Qty,
		Price:          resProduct.Price,
		PriceFormatted: c.mapper.FormatPrice(resProduct.Price),
		PriceTotal:     subTotal,
		TotalFormatted: c.mapper.FormatPrice(subTotal),
	})

	*totalPrice += subTotal
	// check required qty
	if item.Qty >= resDetailDiscount.RequiredQty {
		// get free products
		freeProduct, err := c.MasterProductsRepository.DetailProduct(uint(resDetailDiscount.FreeIDProduct))
		if err != nil {
			c.Log.Warnf("Failed checkout : %+v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return helpers.NewError(404, nil, "01", "Failed checkout, try again.")
			}
			return helpers.NewError(500, nil, "500", err.Error())
		}

		// insert free products to detail transaction
		err = c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
			TransactionID: int(trans.ID),
			ProductID:     int(freeProduct.ID),
			ProductName:   &freeProduct.Name,
			SKU:           &freeProduct.SKU,
			Qty:           1,
			Price:         0,
			Status:        0,
			TotalPrice:    0,
		})
		if err != nil {
			c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
				return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
			}
			return helpers.NewError(500, nil, "500", err.Error())
		}

		// map response items
		productResp = append(productResp, model.ProductCheckoutResponse{
			ProductID:      int(freeProduct.ID),
			ProductName:    freeProduct.Name,
			Qty:            1,
			Price:          0,
			PriceFormatted: c.mapper.FormatPrice(0),
			PriceTotal:     0,
			TotalFormatted: c.mapper.FormatPrice(0),
		})

		desc := fmt.Sprintf("Each sale of a %s comes with a free %s", resProduct.Name, freeProduct.Name)
		discountDescriptions = append(discountDescriptions, desc)
	}
	*itemResp = append(*itemResp, model.CheckoutItemResponse{
		Products:            productResp,
		DiscountDescription: discountDescriptions,
		PriceTotal:          *totalPrice,
		TotalFormatted:      c.mapper.FormatPrice(*totalPrice),
	})
	return nil
}

func (c transactionUseCase) discBuyXGetY(tx *gorm.DB, resDetailDiscount entity.Discount, resProduct entity.Products, item model.CheckoutItem, trans entity.Transaction, totalPrice *float64, itemResp *[]model.CheckoutItemResponse) error {
	var productResp []model.ProductCheckoutResponse
	var discountDescriptions []string

	required := resDetailDiscount.RequiredQty
	final := resDetailDiscount.FinalQty

	if required > 0 && final > 0 && item.Qty >= required {
		times := item.Qty / required
		remain := item.Qty % required
		payQty := times*final + remain

		subTotal := resProduct.Price * float64(payQty)

		err := c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
			TransactionID: int(trans.ID),
			ProductID:     int(resProduct.ID),
			ProductName:   &resProduct.Name,
			SKU:           &resProduct.SKU,
			Qty:           item.Qty,
			Price:         resProduct.Price,
			Status:        0,
			TotalPrice:    subTotal,
		})
		if err != nil {
			c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
				return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
			}
			return helpers.NewError(500, nil, "500", err.Error())
		}

		productResp = append(productResp, model.ProductCheckoutResponse{
			ProductID:      int(resProduct.ID),
			ProductName:    resProduct.Name,
			Qty:            item.Qty,
			Price:          resProduct.Price,
			PriceFormatted: c.mapper.FormatPrice(resProduct.Price),
			PriceTotal:     subTotal,
			TotalFormatted: c.mapper.FormatPrice(subTotal),
		})

		*totalPrice += subTotal

		desc := fmt.Sprintf("Buy %d %s for the price of %d", required, resProduct.Name, final)
		discountDescriptions = append(discountDescriptions, desc)
	} else {
		subTotal := resProduct.Price * float64(item.Qty)

		err := c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
			TransactionID: int(trans.ID),
			ProductID:     int(resProduct.ID),
			ProductName:   &resProduct.Name,
			SKU:           &resProduct.SKU,
			Qty:           item.Qty,
			Price:         resProduct.Price,
			Status:        0,
			TotalPrice:    subTotal,
		})
		if err != nil {
			c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
				return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
			}
			return helpers.NewError(500, nil, "500", err.Error())
		}
		productResp = append(productResp, model.ProductCheckoutResponse{
			ProductID:      int(resProduct.ID),
			ProductName:    resProduct.Name,
			Qty:            item.Qty,
			Price:          resProduct.Price,
			PriceFormatted: c.mapper.FormatPrice(resProduct.Price),
			PriceTotal:     subTotal,
			TotalFormatted: c.mapper.FormatPrice(subTotal),
		})
		*totalPrice += subTotal
	}
	*itemResp = append(*itemResp, model.CheckoutItemResponse{
		Products:            productResp,
		DiscountDescription: discountDescriptions,
		PriceTotal:          *totalPrice,
		TotalFormatted:      c.mapper.FormatPrice(*totalPrice),
	})
	return nil
}

func (c transactionUseCase) discPercentAmount(tx *gorm.DB, resDetailDiscount entity.Discount, resProduct entity.Products, item model.CheckoutItem, trans entity.Transaction, totalPrice *float64, itemResp *[]model.CheckoutItemResponse) error {
	var productResp []model.ProductCheckoutResponse
	var discountDescriptions []string

	required := resDetailDiscount.RequiredQty
	amount := resDetailDiscount.Amount
	isPercentage := resDetailDiscount.IsPercentage
	if item.Qty >= required {
		subTotal := resProduct.Price * float64(item.Qty)

		if isPercentage == 1 {
			discountAmount := (subTotal * float64(amount)) / 100
			subTotal -= discountAmount
		} else {
			subTotal -= float64(amount)
		}

		err := c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
			TransactionID: int(trans.ID),
			ProductID:     int(resProduct.ID),
			ProductName:   &resProduct.Name,
			SKU:           &resProduct.SKU,
			Qty:           item.Qty,
			Price:         resProduct.Price,
			Status:        0,
			TotalPrice:    subTotal,
		})
		if err != nil {
			c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
				return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
			}
			return helpers.NewError(500, nil, "500", err.Error())
		}

		productResp = append(productResp, model.ProductCheckoutResponse{
			ProductID:      int(resProduct.ID),
			ProductName:    resProduct.Name,
			Qty:            item.Qty,
			Price:          resProduct.Price,
			PriceFormatted: c.mapper.FormatPrice(resProduct.Price),
			PriceTotal:     subTotal,
			TotalFormatted: c.mapper.FormatPrice(subTotal),
		})

		*totalPrice += subTotal
		if isPercentage == 1 {
			desc := fmt.Sprintf("Buy %d %s, get %.0f%% off", required, resProduct.Name, float64(amount))
			discountDescriptions = append(discountDescriptions, desc)
		} else {
			desc := fmt.Sprintf("Buy %d %s, get $%.2f off", required, resProduct.Name, float64(amount))
			discountDescriptions = append(discountDescriptions, desc)
		}
	}
	*itemResp = append(*itemResp, model.CheckoutItemResponse{
		Products:            productResp,
		DiscountDescription: discountDescriptions,
		PriceTotal:          *totalPrice,
		TotalFormatted:      c.mapper.FormatPrice(*totalPrice),
	})
	return nil
}

func (c transactionUseCase) addProductToTransaction(tx *gorm.DB, product entity.Products, qty int, trans entity.Transaction, totalPrice *float64, itemResp *[]model.CheckoutItemResponse) error {
	subTotal := product.Price * float64(qty)

	err := c.TransactionRepository.AddDetailTransaction(tx, entity.TransactionDetail{
		TransactionID: int(trans.ID),
		ProductID:     int(product.ID),
		ProductName:   &product.Name,
		SKU:           &product.SKU,
		Qty:           qty,
		Price:         product.Price,
		Status:        0,
		TotalPrice:    subTotal,
	})
	if err != nil {
		c.Log.Warnf("[0] Failed to insert transaction : %+v", err)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && (mysqlErr.Number == 1364 || mysqlErr.Number == 1054) {
			return helpers.NewError(400, nil, "01", "Failed checkout, try again.")
		}
		return helpers.NewError(500, nil, "500", err.Error())
	}

	// Add product response to itemResp
	*itemResp = append(*itemResp, model.CheckoutItemResponse{
		Products: []model.ProductCheckoutResponse{
			{
				ProductID:      int(product.ID),
				ProductName:    product.Name,
				Qty:            qty,
				Price:          product.Price,
				PriceFormatted: c.mapper.FormatPrice(product.Price),
				PriceTotal:     subTotal,
				TotalFormatted: c.mapper.FormatPrice(subTotal),
			},
		},
		DiscountDescription: []string{},
		PriceTotal:          subTotal,
		TotalFormatted:      c.mapper.FormatPrice(subTotal),
	})

	*totalPrice += subTotal
	return nil
}
