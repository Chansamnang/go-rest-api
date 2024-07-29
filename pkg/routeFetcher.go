package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api/internal/config"
	"go-rest-api/models"
	"go-rest-api/models/repositories"
	"reflect"
	"runtime"
	"strings"
)

var superadmin = "superadmin"

type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

func RouteFetcher(engine *gin.Engine) {
	routes := engine.Routes()
	var allPermission []models.Permission
	for _, route := range routes {

		name := getShortHandlerName(route.HandlerFunc)
		exitPermission, err := repositories.PermissionRepository.FindByName(name)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error query finding existing permission %v", err))
		}
		if exitPermission == nil {
			newPermission := models.Permission{
				Name:   name,
				Method: route.Method,
				Route:  route.Path,
			}
			allPermission = append(allPermission, newPermission)
		}
	}

	role, err := repositories.RoleRepository.FindByName(superadmin)
	if role == nil {
		defaultRole := models.Role{
			Name:        superadmin,
			Status:      true,
			Description: superadmin,
			Permissions: allPermission,
		}
		err = repositories.RoleRepository.SaveRole(&defaultRole)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error creating default superadmin role %v", err))
			return
		}

		role, err = repositories.RoleRepository.FindByName(superadmin)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error finding role by name after creation: %v", err))
			return
		}
	} else {
		role.Permissions = allPermission
		err = repositories.RoleRepository.SaveRole(role)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error updating role by name after creation: %v", err))
			return
		}
	}

	user, err := repositories.UserRepository.GetUserByUsername(superadmin)
	if err != nil {
		config.Logger.Error(fmt.Sprintf("Error finding role by name after creation: %v", err))
		return
	}

	if user == nil {
		defaultUser := models.User{
			Username: superadmin,
			Password: "123456",
			RoleID:   role.ID,
		}
		err = repositories.UserRepository.CreateUser(&defaultUser)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error creating default superadmin account %v", err))
			return
		}
		user, err = repositories.UserRepository.GetUserByUsername(superadmin)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error finding role by name after creation: %v", err))
			return
		}
	}
}

func getShortHandlerName(handler gin.HandlerFunc) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}
