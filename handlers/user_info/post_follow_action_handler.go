package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type FollowActionFlow struct {
	c *gin.Context

	userId     int64
	toUserId   int64
	actionType int64
}

func PostFollowActionHandler(c *gin.Context) {
	NewFollowActionFlow(c).Do()
}

func NewFollowActionFlow(c *gin.Context) *FollowActionFlow {
	return &FollowActionFlow{c: c}
}

func (p *FollowActionFlow) Do() {

}

func (p *FollowActionFlow) prepareNum() error {
	rawUserId, _ := p.c.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId

	//解析需要关注的id
	followId := p.c.Query("to_user_id")
	parseInt, err := strconv.ParseInt(followId, 10, 64)
	if err != nil {
		return err
	}
	p.toUserId = parseInt

	//解析action_type
	actionType := p.c.Query("action_type")
	parseInt, err = strconv.ParseInt(actionType, 10, 32)
	if err != nil {
		return err
	}
	p.actionType = parseInt

	return nil
}
