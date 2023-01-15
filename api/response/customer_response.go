package response

type Customer struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Error           string `json:"error,omitempty"`
	Data            string `json:"data,omitempty"`
}
