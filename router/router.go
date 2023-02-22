package router

import (
	"github.com/gin-gonic/gin"
	"miniDy/handlers/comment"
	"miniDy/handlers/message"
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
	baseGroup.GET("/user", middleware.JWTMiddleWare, user_info.UserInfoHandler)            // 用户信息接口完成(xqy)
	baseGroup.GET("/publish/list/", middleware.JWTMiddleWare, video.QueryVideoListHandler) // 发布列表接口完成(xqy)
	baseGroup.POST("/user/login/", user_login.UserLoginHandler)
	baseGroup.POST("/user/register/", user_login.UserRegisterHandler)
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare, video.PublishVideoHandler)

	//interaction apis
	baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare, video.PostFavorHandler)
	baseGroup.GET("/favorite/list/", middleware.CheckIdMiddleWare, video.QueryFavorListHandler)
	baseGroup.POST("/comment/action/", middleware.JWTMiddleWare, comment.PostCommentHandler)
	baseGroup.GET("/comment/list/", middleware.JWTMiddleWare, comment.QueryCommentListHandler)

	//social apis
	baseGroup.POST("/relation/action/", middleware.JWTMiddleWare, user_info.PostFollowHandler)
	baseGroup.GET("/relation/follow/list/", middleware.CheckIdMiddleWare, user_info.GetFollowListHandler)
	baseGroup.GET("/relation/follower/list/", middleware.CheckIdMiddleWare, user_info.GetFollowerListHandler)
	baseGroup.GET("/relation/friend/list/", middleware.CheckIdMiddleWare, user_info.GetFriendListHandler)
	baseGroup.GET("/message/chat/", middleware.JWTMiddleWare, message.GetChattingRecordsHandler)
	baseGroup.POST("/message/action/", middleware.JWTMiddleWare, message.PostMessageActionHandler)

	return r
}
