package core

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func APINoMethodFound() gin.H {
	return APIBadResponse("invalid_action")
}

func APIBadResponse(error string) gin.H {
	return gin.H{"error": error}
}

type HTTPErrorResponse struct {
	Code    string
	Message string
	Params  []interface{}
}

func (her *HTTPErrorResponse) Error() string {
	ret, _ := json.Marshal(her)
	return string(ret)
}

func NewHTTPErrorResponse(code string, message string, params ...interface{}) *HTTPErrorResponse {
	return &HTTPErrorResponse{Code: code, Message: message, Params: params}
}

func APIBadResponseWithCode(code string, error string, params ...string) gin.H {
	return gin.H{"Code": code, "Message": error, "Params": params}
}

func APISuccessResp() gin.H {
	return gin.H{"status": true}
}
