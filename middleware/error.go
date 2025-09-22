package middleware

import (
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware обрабатывает ошибки и формирует единообразный ответ.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Gin error: %v", err.Err)
				if hub := sentry.GetHubFromContext(c.Request.Context()); hub != nil {
					hub.CaptureException(err.Err)
				}
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": c.Errors[0].Error(),
			})
			return
		}
	}
}
