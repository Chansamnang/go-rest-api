package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api/internal/config"
	"go-rest-api/models"
	"go-rest-api/models/requests"
	"go-rest-api/models/response"
	"go-rest-api/pkg"
	userService "go-rest-api/services"
	"go.uber.org/zap"
	"net/http"
)

func UserRegisterHandlers(g *gin.RouterGroup) {
	user := g.Group("/user")
	{
		user.POST("", CreateUser)
		user.GET("", GetAllUsers)
		user.GET("info", UserGetInfo)
	}
}

func UserLogin(c *gin.Context) {
	var userRequest requests.UserLoginRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		pkg.ApiResponse(c, http.StatusBadRequest, err.Error(), "")
		return
	}

	user, err := userService.GetUserByUsername(userRequest.Username)
	if err != nil || user == nil {
		pkg.ApiResponse(c, http.StatusForbidden, "invalid_user", "")
		return
	}
	validPassword := pkg.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if !validPassword {
		pkg.ApiResponse(c, http.StatusUnauthorized, "invalid_user", "")
		return
	}

	token, err := pkg.GenerateJwtToken(user)
	if err != nil {
		config.Logger.Error("Error generating token", zap.Error(err))
		return
	}

	//err = config.Cache.Set(models.UserToken+fmt.Sprintf("%v", user.ID), token, 24*60*60)
	err = config.Cache.Set(models.UserToken+fmt.Sprintf("%v", user.ID), base64.StdEncoding.EncodeToString([]byte(token)), 24*60*60)
	if err != nil {
		config.Logger.Error("Error generating token", zap.Error(err))
		return
	}

	pkg.ApiResponse(c, http.StatusOK, pkg.Success, response.UserLoginResponse{
		UserId:   user.ID,
		Username: user.Username,
		Token:    token,
	})
}

func UserGetInfo(c *gin.Context) {
	user, err := pkg.CurrentUser(c)
	if err != nil {
		pkg.ApiResponse(c, http.StatusBadRequest, pkg.Error, err.Error())
		return
	}
	pkg.ApiResponse(c, http.StatusOK, pkg.Success, response.UserInfoResponse{}.ToFormat(user))
}

func CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		pkg.ApiResponse(c, http.StatusBadRequest, pkg.Error, err.Error())
		return
	}

	if err := userService.CreateUser(user); err != nil {
		pkg.ApiResponse(c, http.StatusBadRequest, pkg.Error, err.Error())
		return
	}

	pkg.ApiResponse(c, http.StatusCreated, pkg.Created, "")
}

func GetAllUsers(c *gin.Context) {
	users, err := userService.GetAllUser()
	if err != nil {
		pkg.ApiResponse(c, http.StatusBadRequest, pkg.Error, err.Error())
		return
	}

	pkg.ApiResponse(c, http.StatusOK, pkg.Success, users)
}
