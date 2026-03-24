package response

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
}

func Success(c fiber.Ctx, message string, data interface{}, statusCode int) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:     "success",
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	})
}

func Error(c fiber.Ctx, message string, statusCode int) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:     "error",
		Message:    message,
		Data:       nil,
		StatusCode: statusCode,
	})
}

func ValidationError(c fiber.Ctx, message string, errors interface{}) error {
	return c.Status(fiber.ErrUnprocessableEntity.Code).JSON(APIResponse{
		Status:  "error",
		Message: message,
		Data: fiber.Map{
			"errors": errors,
		},
		StatusCode: fiber.ErrUnprocessableEntity.Code,
	})
}
