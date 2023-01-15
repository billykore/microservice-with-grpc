package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice-with-grpc/api/helper"
	"microservice-with-grpc/api/request"
	"microservice-with-grpc/api/response"
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
		log.Printf("[handler error] error binding request body to json. %v", err)
		ctx.JSON(http.StatusBadRequest, &response.Customer{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed create new account",
			Error:           "Request body is not valid",
		})
		return
	}
	grpcRequest := helper.BuildGrpcRequest(body)
	grpcResponse, err := h.Client.AccountCreation(ctx, grpcRequest)
	if err != nil {
		log.Printf("[grpc response] %v", grpcResponse)
		log.Printf("[handler error] error create account from grpc service. %v", err)
		ctx.JSON(http.StatusServiceUnavailable, &response.Customer{
			ResponseCode:    http.StatusServiceUnavailable,
			ResponseMessage: "Failed create new account",
			Error:           "Account service unavailable. Please try again later",
		})
	} else {
		log.Printf("[grpc response] %v", grpcResponse)
		ctx.JSON(http.StatusOK, &response.Customer{
			ResponseCode:    http.StatusOK,
			ResponseMessage: grpcResponse.GetMessage(),
		})
	}
}
