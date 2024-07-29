package pkg

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/lang"
)

func ApiResponse(c *gin.Context, code int, msg string, i interface{}) {
	locale := lang.GetLang(c)
	statusCode := generateStatusCode(code)

	c.JSON(statusCode, response{
		Code:    code,
		Message: lang.T(locale, msg),
		Data:    i,
	})
}
