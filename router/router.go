package router

import (
	"github.com/gin-gonic/gin"
	"miniDy/handlers/user_info"
	"miniDy/handlers/user_login"
	"miniDy/handlers/video"
	"miniDy/middleware"
	"miniDy/model"
)

func InitRouter() *gin.Engine {
	model.InitDB()

	r := gin.Default()

	r.Static("/static/", "./static")

	baseGroup := r.Group("/douyin")

	//basic apis
	baseGroup.GET("/feed", video.GetVideoFeedHandler)
	baseGroup.GET("/user", middleware.JWTMiddleWare, user_info.UserInfoHandler)           // 用户信息接口完成(xqy)
	baseGroup.GET("/publish/list", middleware.JWTMiddleWare, video.QueryVideoListHandler) // 发布列表接口完成(xqy)
	baseGroup.POST("/user/login/", user_login.UserLoginHandler)
	baseGroup.POST("/user/register/", user_login.UserRegisterHandler)
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare, video.PublishVideoHandler)

	//interaction apis
	baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare)
	baseGroup.GET("/favorite/list", middleware.CheckIdMiddleWare)
	baseGroup.POST("/comment/action", middleware.JWTMiddleWare)
	baseGroup.GET("/comment/list", middleware.JWTMiddleWare)

	//social apis
	baseGroup.POST("/relation/action/", middleware.JWTMiddleWare, user_info.PostFollowHandler) // 关注操作接口完成(xqy)
	baseGroup.GET("/relation/follow/list", middleware.CheckIdMiddleWare)
	baseGroup.GET("/favorite/follower/list", middleware.CheckIdMiddleWare)
	baseGroup.GET("/favorite/friend/list", middleware.CheckIdMiddleWare)
	//baseGroup.GET("/favorite/message/chat")
	//baseGroup.GET("/favorite/message/action")

	return r
}
