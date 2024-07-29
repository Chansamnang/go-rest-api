package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api/internal/config"
	"go-rest-api/models"
	"go-rest-api/models/repositories"
	"go-rest-api/pkg"
	"go.uber.org/zap"
	"net/http"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := pkg.GetTokenFromRequest(c)

		if len(tokenString) <= 7 {
			pkg.ApiResponse(c, http.StatusUnauthorized, "invalid_token", "")
			c.Abort()
			return
		}

		claims, err := pkg.ParseJwtClaimToken(tokenString)
		if err != nil {
			pkg.ApiResponse(c, http.StatusUnauthorized, "invalid_token", "")
			c.Abort()
			return
		}

		encodeToken, err := config.Cache.Get(models.UserToken + fmt.Sprintf("%d", claims.UserId))
		if err != nil {
			pkg.ApiResponse(c, http.StatusUnauthorized, "invalid_token", "")
			c.Abort()
			return
		}

		decodeString, err := base64.StdEncoding.DecodeString(encodeToken)
		if err != nil {
			config.Logger.Error("decode token error", zap.Error(err))
			c.Abort()
			return
		}

		if string(decodeString) != tokenString {
			pkg.ApiResponse(c, http.StatusUnauthorized, "invalid_token", "")
			c.Abort()
			return
		}

		//checking permission
		permission, err := repositories.PermissionRepository.FindByRouteAndMethod(c.Request.RequestURI, c.Request.Method)
		if err != nil {
			config.Logger.Error("[Query Permission] error", zap.Error(err))
			c.Abort()
			return
		}

		fmt.Println(permission)

		user, err := repositories.UserRepository.GetUserById(claims.UserId)
		if err != nil {
			config.Logger.Error("[Query User] error", zap.Error(err))
			c.Abort()
			return
		}

		if !hasPermission(user.Role.Permissions, permission) {
			pkg.ApiResponse(c, http.StatusForbidden, "permission_denied", "")
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(allPerm []models.Permission, requestPerm *models.Permission) bool {
	for _, perm := range allPerm {
		if perm.ID == requestPerm.ID {
			return true
		}
	}
	return false
}
