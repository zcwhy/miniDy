package message

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/message"
	"net/http"
	"strconv"
)

type PostMessageActionFlow struct {
	fromUserId int64
	toUserId   int64
	actionType int32
	content    string
	c          *gin.Context
}

func PostMessageActionHandler(c *gin.Context) {
	err := (&PostMessageActionFlow{c: c}).Do()

	if err != nil {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.CommonResp{
		StatusCode: constant.SUCCESS,
		StatusMsg:  constant.SUCCESS_MESSAGE,
	})
}

func (p *PostMessageActionFlow) Do() error {
	if err := p.parseParam(); err != nil {
		return err
	}

	if err := p.doService(); err != nil {
		return err
	}

	return nil
}

func (p *PostMessageActionFlow) doService() error {
	return message.PostMessage(p.fromUserId, p.toUserId, p.actionType, p.content)
}

func (p *PostMessageActionFlow) parseParam() error {
	rawToUserid := p.c.Query("to_user_id")
	rawFromUserId, _ := p.c.Get("user_id")
	rawActionType := p.c.Query("action_type")
	content := p.c.Query("content")

	toUserId, err := strconv.ParseInt(rawToUserid, 10, 64)
	if err != nil {
		return err
	}
	p.toUserId = toUserId

	actionType, err := strconv.ParseInt(rawActionType, 10, 32)
	if err != nil {
		return err
	}
	p.actionType = int32(actionType)
	p.fromUserId = rawFromUserId.(int64)
	p.content = content

	return nil
}
