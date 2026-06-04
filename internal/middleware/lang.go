package middleware

import "github.com/gin-gonic/gin"

func LangMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if lang != "ru" && lang != "en" {
			lang = "ru"
		}
		c.Set("lang", lang)
		c.Next()
	}
}
