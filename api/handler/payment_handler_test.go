package handler

import (
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	pb "microservice-with-grpc/gen/payment/v1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func paymentRoute(h *PaymentHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/payment/qris", h.Qris)
	return r
}

func paymentGrpcClientConn() (*grpc.ClientConn, func()) {
	paymentConn, err := grpc.Dial("172.22.0.1:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) //docker
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	closer := func() {
		err = paymentConn.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
	}
	return paymentConn, closer
}

func TestPaymentHandler_Qris(t *testing.T) {
	type args struct {
		request string
	}

	type expectation struct {
		statusCode int
		response   gin.H
	}

	// test cases
	tests := map[string]struct {
		args     args
		expected expectation
	}{
		"success": {
			args: args{
				request: `{"merchantId":"M-001","trxNumber":"0000001","accountSource":"001001000001300","accountDestination":"001001000002300","amount":"500000"}`,
			},
			expected: expectation{
				statusCode: 200,
				response: gin.H{
					"responseCode":    200,
					"responseMessage": "Qris payment successful",
				},
			},
		},
		"request_body_not_valid": {
			args: args{
				request: "",
			},
			expected: expectation{
				statusCode: 400,
				response: gin.H{
					"responseCode":    400,
					"responseMessage": "Failed to process qris payment",
					"error":           "Request body is not valid",
				},
			},
		},
		"payment_service_unavailable": {
			args: args{
				request: `{"merchantId":"M-001","trxNumber":"0000001","accountSource":"001001000001300","accountDestination":"001001000002300","amount":"500000"}`,
			},
			expected: expectation{
				statusCode: 503,
				response: gin.H{
					"responseCode":    503,
					"responseMessage": "Failed to process qris payment",
					"error":           "Payment service unavailable. Please try again later",
				},
			},
		},
	}

	paymentGrpcConn, closer := paymentGrpcClientConn()
	defer closer()
	paymentClient := pb.NewPaymentClient(paymentGrpcConn)
	payment := NewPaymentHandler(paymentClient)
	route := paymentRoute(payment)

	// run the test
	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/payment/qris",
				strings.NewReader(test.args.request))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			route.ServeHTTP(rec, req)
			response := rec.Result()
			assert.Equal(t, test.expected.statusCode, response.StatusCode)

			body, _ := io.ReadAll(response.Body)
			var responseBody map[string]interface{}
			err := json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, test.expected.response["responseCode"].(int), responseBody["responseCode"].(int))
			assert.Equal(t, test.expected.response["responseMessage"].(string), responseBody["responseCode"].(string))
			assert.Equal(t, test.expected.response["error"].(string), responseBody["error"].(string))
		})
	}
}
