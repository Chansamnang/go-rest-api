package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-rest-api/models"
	"go-rest-api/services"
	"os"
	"strings"
	"time"
)

var TokenTTL = time.Hour * 24 * 7

type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var privateKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJwtToken(user *models.User) (string, error) {
	claims := Claims{
		UserId:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(privateKey)
}

func CurrentUser(c *gin.Context) (*models.User, error) {
	authToken := GetTokenFromRequest(c)
	claims, err := ParseJwtClaimToken(authToken)
	if err != nil {
		return nil, err
	}
	userId := claims.UserId

	user, err := services.GetUserById(userId)
	if err != nil {
		return &models.User{}, err
	}
	user.Password = ""
	return user, nil
}

func ParseJwtClaimToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return privateKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func GetTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
