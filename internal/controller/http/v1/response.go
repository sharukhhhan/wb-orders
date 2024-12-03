package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Content interface{} `json:"content,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c echo.Context, statusCode int, err error) error {
	var response ErrorResponse
	if statusCode != http.StatusInternalServerError {
		response.Error = err.Error()
	} else {
		response.Error = http.StatusText(http.StatusInternalServerError)
	}

	errJSON := c.JSON(statusCode, response)
	if errJSON != nil {
		return fmt.Errorf("error while returning json: %w", errJSON)
	}
	return err
}

func newSuccessResponse(c echo.Context, message string, content interface{}) error {
	response := SuccessResponse{
		Message: message,
		Content: content,
	}
	return c.JSON(http.StatusOK, response)
}
