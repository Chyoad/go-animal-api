package response

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Struktur standar untuk response API.
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Succcess response.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// Error response.
func Error(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	var errorMsg string
	if errDetail != nil {
		if e, ok := errDetail.(error); ok {
			errorMsg = e.Error()
		} else if s, ok := errDetail.(string); ok {
			errorMsg = s
		}
	}

	fullMessage := message
	// Jika errDetail (string dari error) tidak kosong dan belum ada dalam message utama, gabungkan.
	if errorMsg != "" {
		if message != "" && !strings.Contains(message, errorMsg) {
			fullMessage = message + ": " + errorMsg
		} else if message == "" {
			fullMessage = errorMsg // Jika message utama kosong, gunakan errorMsg sebagai fullMessage
		}
	}
	
	// Jika errDetail adalah string (errorMsg), dan sudah digabungkan ke fullMessage, bisa set Error ke nil atau tetap errorMsg.
	var errorField interface{}
	if errorMsg != "" {
		errorField = errorMsg
	} else if errDetail != nil && errorMsg == "" { // Jika errDetail bukan error.Error() atau string
		errorField = errDetail
	}


	c.JSON(statusCode, APIResponse{
		Status:  statusCode,
		Message: fullMessage,
		Error:   errorField,
	})
}