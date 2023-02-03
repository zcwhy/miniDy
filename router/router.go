package router

import (
	"github.com/gin-gonic/gin"
	"miniDy/handlers/user_login"
	"miniDy/middleware"
	"miniDy/model"
)

func InitRouter() *gin.Engine {
	model.InitDB()

	r := gin.Default()

	r.Static("/static", "./static")

	baseGroup := r.Group("/douyin")

	//basic apis
	baseGroup.GET("/feed")
	baseGroup.GET("/user", middleware.JWTMiddleWare)
	baseGroup.GET("/publish/list", middleware.CheckIdMiddleWare)
	baseGroup.POST("/user/login")
	baseGroup.POST("/user/register", user_login.UserRegisterHandler)
	baseGroup.POST("/publish/action", middleware.JWTMiddleWare)

	//interaction apis
	baseGroup.POST("/favorite/action", middleware.JWTMiddleWare)
	baseGroup.GET("/favorite/list", middleware.CheckIdMiddleWare)
	baseGroup.POST("/comment/action", middleware.JWTMiddleWare)
	baseGroup.GET("/comment/list", middleware.JWTMiddleWare)

	//social apis
	baseGroup.POST("/relation/action", middleware.JWTMiddleWare)
	baseGroup.GET("/relation/follow/list", middleware.CheckIdMiddleWare)
	baseGroup.GET("/favorite/follower/list", middleware.CheckIdMiddleWare)
	baseGroup.GET("/favorite/friend/list", middleware.CheckIdMiddleWare)
	//baseGroup.GET("/favorite/message/chat")
	//baseGroup.GET("/favorite/message/action")

	return r
}
