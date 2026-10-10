package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func OK(c *gin.Context, status int, msg string, data any) {
	c.JSON(status, Response{Success: true, Message: msg, Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{Success: false, Message: msg})
}
