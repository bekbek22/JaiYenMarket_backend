package utils

import "github.com/gofiber/fiber/v2"

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, status int, data interface{}, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return c.Status(status).JSON(SuccessResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, status int, message string, err error) error {
	var errDetail any
	if err != nil {
		errDetail = err.Error()
	}
	return c.Status(status).JSON(ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   errDetail,
	})
}
