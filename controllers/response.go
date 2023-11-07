package controllers

import (
	"github.com/gin-gonic/gin"
)

// APIResponse defines a standard format for API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendResponse is a utility function to send a JSON response.
func SendResponse(c *gin.Context, httpStatus int, data interface{}, message string, err error) {
	// Prepare the response based on the presence of an error
	var res APIResponse
	if err != nil {
		res = APIResponse{Success: false, Error: err.Error()}
		c.JSON(httpStatus, res)
		return
	}

	res = APIResponse{Success: true, Data: data, Message: message}
	c.JSON(httpStatus, res)
}
