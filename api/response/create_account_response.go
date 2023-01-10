package response

type CreateAccount struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Error           string `json:"error,omitempty"`
}
