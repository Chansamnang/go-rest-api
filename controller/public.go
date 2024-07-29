package controller

import "github.com/gin-gonic/gin"

func PublicRegisterHandler(g *gin.RouterGroup) {
	public := g.Group("/public")
	{
		public.POST("login", UserLogin)
	}
}
