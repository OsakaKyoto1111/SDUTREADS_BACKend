package utils

import (
	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, status int, data interface{}) error {
	payload := echo.Map{"data": data}
	if data == nil {
		payload = echo.Map{"data": echo.Map{}}
	}
	return c.JSON(status, payload)
}

func Error(c echo.Context, status int, err error) error {
	return c.JSON(status, echo.Map{"error": err.Error()})
}
