package model

import (
	"base-golang/internal/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
)

type WebResponse[T any] struct {
	ResponseCode   string        `json:"responseCode"`
	ResponseDesc   string        `json:"responseDesc"`
	ResponseData   T             `json:"responseData"`
	Paging         *PageMetadata `json:"paging,omitempty"`
	ResponseErrors any           `json:"responseErrors,omitempty"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func ResponseSuccess[T any](ctx *fiber.Ctx, desc string, data T, paging *PageMetadata) (res error) {

	response := WebResponse[T]{
		ResponseCode: "00",
		ResponseDesc: desc,
		ResponseData: data,
		Paging:       paging,
	}
	var thirdParty []helpers.ThirdPartyLog
	thirdPartyExist := ctx.Context().Value("thirdParty")
	if thirdPartyExist != nil {
		thirdParty = thirdPartyExist.([]helpers.ThirdPartyLog)
	}
	param := helpers.LogResponseParam{
		ResponseCode:   response.ResponseCode,
		ResponseBody:   response,
		ResponseHeader: ctx.GetRespHeaders(),
		Timestamp:      time.Now().String(),
		ThirdParty:     thirdParty,
	}
	helpers.LogResponse(param)
	res = ctx.JSON(response)
	return
}

func ResponseError[T any](ctx *fiber.Ctx, http_code int, code, desc string, errors any) (res error) {
	response := WebResponse[T]{
		ResponseCode:   code,
		ResponseDesc:   desc,
		ResponseErrors: errors,
	}

	var thirdParty []helpers.ThirdPartyLog
	thirdPartyExist := ctx.Context().Value("thirdParty")
	if thirdPartyExist != nil {
		thirdParty = thirdPartyExist.([]helpers.ThirdPartyLog)
	}
	param := helpers.LogResponseParam{
		ResponseCode:   response.ResponseCode,
		ResponseBody:   response,
		ResponseHeader: ctx.GetRespHeaders(),
		Timestamp:      time.Now().String(),
		ThirdParty:     thirdParty,
	}
	helpers.LogResponse(param)
	res = ctx.Status(http_code).JSON(response)
	return
}

func ResponseSuccessCustom[T any](ctx *fiber.Ctx, desc string, data T, paging *PageMetadata, code string) (res error) {

	response := WebResponse[T]{
		ResponseCode: code,
		ResponseDesc: desc,
		ResponseData: data,
		Paging:       paging,
	}
	var thirdParty []helpers.ThirdPartyLog
	thirdPartyExist := ctx.Context().Value("thirdParty")
	if thirdPartyExist != nil {
		thirdParty = thirdPartyExist.([]helpers.ThirdPartyLog)
	}
	param := helpers.LogResponseParam{
		ResponseCode:   response.ResponseCode,
		ResponseBody:   response,
		ResponseHeader: ctx.GetRespHeaders(),
		Timestamp:      time.Now().String(),
		ThirdParty:     thirdParty,
	}
	helpers.LogResponse(param)
	res = ctx.JSON(response)
	return
}
