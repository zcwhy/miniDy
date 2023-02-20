package user_info

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/user_info"
	"net/http"
	"strconv"
)

type GetFriendLisListResponse struct {
	response.CommonResp
	*user_info.UserFriendList
}

func GetFriendListHandler(c *gin.Context) {
	strUserId := c.Query("user_id")
	userId, err := strconv.ParseInt(strUserId, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, GetFriendLisListResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  "用户id解析错误",
			},
		})
		return
	}

	friendList, err := user_info.GetFriendList(userId)

	if err != nil {
		c.JSON(http.StatusOK, GetFriendLisListResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GetFriendLisListResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		UserFriendList: friendList,
	})
}
