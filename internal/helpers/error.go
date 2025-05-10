package helpers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code    int    `json:"code"`
	Rc      string `json:"rc"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func SetFitur(ctx *fiber.Ctx, fitur string) {
	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "fitur", fitur))
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, errors any, message ...string) *Error {
	err := &Error{
		Code:    code,
		Rc:      "99",
		Message: "General Error",
		Errors:  errors,
	}

	if len(message) > 0 {
		err.Rc = message[0]
	}

	if len(message) > 1 {
		err.Message = message[1]
	}
	return err
}

var JWTNotValid = NewError(400, nil, "BX", "JWT Token Not Valid")
var ValidateError = NewError(422, nil, "VE", "Validation Error")
var InvalidAccessMidOrMpan = NewError(422, nil, "R0", "Akses Tidak diizinkan")
var DNA = NewError(200, nil, "DNA", "Data Tidak Diizinkan")
var QRPM = NewError(200, nil, "QRPM", "Jumlah Pembayaran QRIS tidak bisa melebihi 10.000.000")
var FailedTrx = NewError(200, nil, "FTX", "Gagal Membatalkan Transaksi, Coba Lagi")
var DataNotFound = NewError(200, nil, "DNF", "Data Not Found")
var DataNotValid = NewError(200, nil, "DNV", "Data Not Valid")
var DataAllProductNotFound = NewError(200, nil, "DAPNF", "Data All Product Not Found")
var DataAllTransactionNotFound = NewError(200, nil, "DATNF", "Data All Transaction Not Found")
var ServiceThirdPartyError = NewError(200, nil, "99", "Service ThirdParty Error")
var DatabaseError = NewError(200, "DBF", "99", "Database Error")
var TRPM = NewError(200, nil, "TRPM", "Jumlah Pembayaran Tunai tidak bisa melebihi 100.000.000")

// insert Product Bulk Via Excel
var MimeTypeNotAllowed = NewError(200, nil, "UB00", "Mime Type Not Allowed")
var FormatBase64Invalid = NewError(200, nil, "UB01", "Format Base64 Request Invalid")
