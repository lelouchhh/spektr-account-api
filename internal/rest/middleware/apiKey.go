package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// APIKeyMiddleware проверяет наличие и корректность API ключа в заголовке.
func APIKey(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем API ключ из заголовка запроса
			key := c.Request().Header.Get("X-API-Key")
			if key == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing API Key")
			}

			// Проверяем, совпадает ли переданный API ключ с ожидаемым
			if key != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
			}

			// Если ключ правильный, продолжаем обработку
			return next(c)
		}
	}
}
