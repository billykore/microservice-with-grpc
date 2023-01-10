package router

import (
	"github.com/gin-gonic/gin"
	"microservice-with-grpc/api/handler"
)

func New(h handler.Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/account/create", h.Customer.AccountCreation)
	}

	return r
}
