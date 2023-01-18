package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"microservice-with-grpc/api/response"
	pb "microservice-with-grpc/gen/auth/v1"
)

type AuthHandler struct {
	Client pb.AuthClient
}

func NewAuthHandler(client pb.AuthClient) *AuthHandler {
	return &AuthHandler{Client: client}
}

func (h *AuthHandler) GetToken(ctx *gin.Context) {
	grantType := ctx.PostForm("grantType")
	if grantType != "password" {
		log.Printf("[handler error] invalid grantType: %v", grantType)
		ctx.JSON(http.StatusBadRequest, &response.Response{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "Failed get token",
			Error:           "Unsupported grant type",
		})
		return
	}
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	grpcRequest := &pb.TokenRequest{
		Username:  username,
		Password:  password,
		GrantType: grantType,
	}
	grpcResponse, err := h.Client.GetToken(ctx, grpcRequest)
	log.Printf("[grpc response] %v", grpcResponse)
	if err != nil {
		log.Printf("[handler error] error get token from grpc service: %v", err)
		ctx.JSON(http.StatusUnauthorized, &response.Response{
			ResponseCode:    http.StatusUnauthorized,
			ResponseMessage: "Failed get token",
			Error:           "Invalid username or password",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Response{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "Success get token",
		Data:            grpcResponse,
	})
}
