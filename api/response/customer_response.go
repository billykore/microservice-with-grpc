package response

type Customer struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Error           string `json:"error,omitempty"`
	Data            any    `json:"data,omitempty"`
}
