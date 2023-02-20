package message

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model"
	"miniDy/model/response"
	"miniDy/service/message"
	"net/http"
	"strconv"
)

var context *gin.Context

type ChattingRecordsResponse struct {
	response.CommonResp
	MessageList []*model.Message `json:"message_list"`
}

func GetChattingRecordsHandler(c *gin.Context) {
	context = c
	rawUserId, _ := c.Get("user_id")
	rawToUserId := c.Query("to_user_id")
	rawTime := c.Query("pre_msg_time")

	userId := rawUserId.(int64)
	toUserId, err := strconv.ParseInt(rawToUserId, 10, 64)
	if err != nil {
		sendErrorResponse(err)
		return
	}

	timeStamp, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		sendErrorResponse(err)
		return
	}

	messageList, err := message.GetMessageRecords(userId, toUserId, timeStamp)
	if err != nil {
		sendErrorResponse(err)
		return
	}

	c.JSON(http.StatusOK, ChattingRecordsResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  "操作成功",
		},
		MessageList: messageList,
	})
}

func sendErrorResponse(err error) {
	context.JSON(http.StatusOK, ChattingRecordsResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  err.Error(),
		},
	})
}
