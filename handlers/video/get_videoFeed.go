package video

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/middleware"
	"miniDy/model/response"
	"miniDy/service/video"
	"net/http"
	"strconv"
	"time"
)

type GetVideoFeedResponse struct {
	response.CommonResp
	*video.GetVideoFeedResponse
}

func GetVideoFeedHandler(c *gin.Context) {
	rawLatestTime := c.Query("latest_time")
	token := c.Query("token")

	latestTime := time.Now().Unix()
	if rawLatestTime != "" {
		var err error
		latestTime, err = strconv.ParseInt(rawLatestTime, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, GetVideoFeedResponse{
				CommonResp: response.CommonResp{
					StatusCode: constant.FAILURE,
					StatusMsg:  "时间解析错误",
				},
			})
			return
		}
	}

	if token != "" {
		valid := middleware.TokenVerify(token)

		if !valid {
			c.JSON(http.StatusOK, GetVideoFeedResponse{
				CommonResp: response.CommonResp{
					StatusCode: constant.FAILURE,
					StatusMsg:  "token错误",
				},
			})
			return
		}
	}

	//前端传来的时间戳以毫秒为单位， go只能用以秒为单位的时间戳
	resp, err := video.GetVideoFeed(latestTime / 1000)

	if err != nil {
		c.JSON(http.StatusOK, GetVideoFeedResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GetVideoFeedResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  "操作成功",
		},
		GetVideoFeedResponse: resp,
	})
}
