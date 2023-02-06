package video

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/video"
	"net/http"
)

func PublishVideoHandler(c *gin.Context) {
	var req video.PublishVideoRequest

	rawId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  "未解析到user_id"})
		return
	}

	req.UserId = rawId.(int64)
	req.Title = c.PostForm("title")

	form, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  err.Error()},
		)
		return
	}
	req.Video = form
	req.Context = c

	err = video.PublishVideo(&req)
	if err != nil {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, response.CommonResp{
		StatusCode: constant.SUCCESS,
		StatusMsg:  constant.SUCCESS_MESSAGE},
	)
}
