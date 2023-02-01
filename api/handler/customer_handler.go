package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice-with-grpc/api/helper"
	"microservice-with-grpc/api/request"
	"microservice-with-grpc/api/response"
	"microservice-with-grpc/api/validator"
	pb "microservice-with-grpc/gen/customer/v1"
)

type CustomerHandler struct {
	Client pb.CustomerClient
}

func NewCustomerHandler(client pb.CustomerClient) *CustomerHandler {
	return &CustomerHandler{
		Client: client,
	}
}

func (h *CustomerHandler) AccountCreation(ctx *gin.Context) {
	body := &request.CreateAccount{}
	err := ctx.ShouldBindJSON(body)
	if err != nil {
		log.Printf("[handler error] error binding request body to json: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed create new account",
			Error:           "Request body is not valid",
		})
		return
	}
	valid, errMsg := validator.ValidateRequestBody(body)
	if !valid {
		log.Printf("[handler error] request body validation error: %v", errMsg)
		ctx.JSON(http.StatusBadRequest, &response.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed create new account",
			Error:           errMsg,
		})
		return
	}
	grpcRequest := helper.BuildCustomerGrpcRequest(body)
	grpcResponse, err := h.Client.AccountCreation(ctx, grpcRequest)
	log.Printf("[customer grpc response] %v", grpcResponse)
	if err != nil {
		log.Printf("[handler error] error create account from customer grpc service: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, &response.Response{
			ResponseCode:    http.StatusServiceUnavailable,
			ResponseMessage: "Failed create new account",
			Error:           "Account service unavailable. Please try again later",
		})
		return
	}
	ctx.JSON(http.StatusOK, &response.Response{
		ResponseCode:    http.StatusOK,
		ResponseMessage: grpcResponse.GetMessage(),
	})
}

func (h *CustomerHandler) AccountInquiry(ctx *gin.Context) {
	accountNumber := ctx.Query("accountNumber")
	if accountNumber == "" {
		log.Println("[handle error] missing query parameter: accountNumber")
		ctx.JSON(http.StatusBadRequest, response.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed inquiry account",
			Error:           "Missing query parameter: accountNumber",
		})
		return
	}
	grpcResponse, err := h.Client.AccountInquiry(ctx, &pb.InquiryRequest{AccountNumber: accountNumber})
	log.Printf("[customer grpc response] %v", grpcResponse)
	if err != nil {
		log.Printf("[handler error] error inquiry account from customer grpc service: %v", err)
		ctx.JSON(http.StatusNotFound, &response.Response{
			ResponseCode:    http.StatusNotFound,
			ResponseMessage: "Failed inquiry account",
			Error:           "Account not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, &response.Response{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success inquiry account",
		Data:            grpcResponse,
	})
}
