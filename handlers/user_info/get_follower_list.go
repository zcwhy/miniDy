package user_info

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/user_info"
	"net/http"
	"strconv"
)

type GetFollowerListResponse struct {
	response.CommonResp
	*user_info.UserFollowerList
}

func GetFollowerListHandler(c *gin.Context) {
	strUserId := c.Query("user_id")
	userId, err := strconv.ParseInt(strUserId, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, GetFollowerListResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  "用户id解析错误",
			},
		})
		return
	}

	followerList, err := user_info.GetFollowerList(userId)

	if err != nil {
		c.JSON(http.StatusOK, GetFollowerListResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GetFollowerListResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		UserFollowerList: followerList,
	})
}
