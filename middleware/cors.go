package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization, Token, did, device,language,device_type")
		c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,OPTIONS,DELETE")
		c.Header("Access-Control-Expose-Headers", "Value-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Value-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
