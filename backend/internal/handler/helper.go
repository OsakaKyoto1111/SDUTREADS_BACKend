package handler

import "github.com/labstack/echo/v4"

// GetUserIDFromContext - безопасно получает user_id из echo.Context
func GetUserIDFromContext(c echo.Context) (uint, bool) {
	val := c.Get("user_id")
	if val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case uint:
		return v, true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint(v), true
	default:
		return 0, false
	}
}
