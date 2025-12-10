package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func respondJSON(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, echo.Map{
		"status": "ok",
		"data":   data,
	})
}

func respondError(c echo.Context, code int, msg string) error {
	return c.JSON(code, echo.Map{
		"status":  "error",
		"message": msg,
	})
}

func bindJSON(c echo.Context, dst interface{}) error {
	if err := c.Bind(dst); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	return nil
}

func parseIDParam(c echo.Context, name string) (uint, error) {
	idStr := c.Param(name)
	u64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return uint(u64), nil
}

func GetUserIDFromContext(c echo.Context) (uint, bool) {
	v := c.Get("user_id")
	if v == nil {
		return 0, false
	}

	switch t := v.(type) {
	case uint:
		return t, true
	case int:
		return uint(t), true
	case float64:
		return uint(t), true
	default:
		return 0, false
	}
}

func requireAuth(c echo.Context) (uint, bool) {
	id, ok := GetUserIDFromContext(c)
	if !ok || id == 0 {
		_ = respondError(c, http.StatusUnauthorized, "unauthorized")
		return 0, false
	}
	return id, true
}
