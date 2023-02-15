package router

import (
	"github.com/gin-gonic/gin"
	"miniDy/handlers/comment"
	"miniDy/handlers/user_login"
	"miniDy/model"
)

func InitRouter() *gin.Engine {
	model.InitDB()

	r := gin.Default()
	baseGroup := r.Group("/douyin")

	//basic apis
	baseGroup.GET("/feed")
	baseGroup.GET("/user")
	baseGroup.GET("/publish/list")
	baseGroup.POST("/user/login")
	baseGroup.POST("/user/register", user_login.UserRegisterHandler)
	baseGroup.POST("/publish/action")

	//interaction apis
	baseGroup.POST("/favorite/action")
	baseGroup.GET("/favorite/list")
	baseGroup.POST("/comment/action", comment.PostCommentHandler)
	baseGroup.GET("/comment/list", comment.QueryCommentListHandler)

	//social apis
	baseGroup.POST("/relation/action")
	baseGroup.GET("/relation/follow/list")
	baseGroup.GET("/favorite/follower/list")
	baseGroup.GET("/favorite/friend/list")
	baseGroup.GET("/favorite/message/chat")
	baseGroup.GET("/favorite/message/action")

	return r
}
