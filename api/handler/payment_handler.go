package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice-with-grpc/api/helper"
	"microservice-with-grpc/api/request"
	"microservice-with-grpc/api/response"
	"microservice-with-grpc/api/validator"
	pb "microservice-with-grpc/gen/payment/v1"
)

type PaymentHandler struct {
	Client pb.PaymentClient
}

func NewPaymentHandler(client pb.PaymentClient) *PaymentHandler {
	return &PaymentHandler{Client: client}
}

func (h *PaymentHandler) Qris(ctx *gin.Context) {
	body := &request.QrisPayment{}
	err := ctx.ShouldBindJSON(body)
	if err != nil {
		log.Printf("[handler error] error binding request body to json: %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed to process qris payment",
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
	grpcRequest := helper.BuildQrisPaymentGrpcRequest(body)
	grpcResponse, err := h.Client.Qris(ctx, grpcRequest)
	log.Printf("[payment grpc response] %v", grpcResponse)
	if err != nil {
		log.Printf("[handler error] qris payment error from customer grpc service: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, &response.Response{
			ResponseCode:    http.StatusServiceUnavailable,
			ResponseMessage: "Failed to process qris payment",
			Error:           "Payment service unavailable. Please try again later",
		})
		return
	}
	ctx.JSON(http.StatusOK, &response.Response{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Qris payment successful",
	})
}
