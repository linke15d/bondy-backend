package middleware

import (
	"github.com/linke15d/bondy-backend/pkg/i18n"

	"github.com/gin-gonic/gin"
)

func Lang() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.ParseLang(c.GetHeader("Accept-Language"))
		c.Set("lang", lang)
		c.Next()
	}
}
