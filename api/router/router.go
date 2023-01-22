package router

import (
	"github.com/gin-gonic/gin"
	"microservice-with-grpc/api/handler"
)

func New(h handler.Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		// Customer routes
		v1.POST("/account/create", h.Customer.AccountCreation)
		v1.GET("/account/inquiry", h.Customer.AccountInquiry)
		// Auth routes
		v1.POST("/auth/token", h.Auth.GetToken)
		// Payment routes
		v1.POST("/payment/qris", h.Payment.Qris)
	}
	return r
}
