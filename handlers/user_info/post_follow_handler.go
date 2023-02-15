package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/user_info"
	"net/http"
	"strconv"
)

// ProxyPostFollow 代理结构体
type ProxyPostFollow struct {
	c *gin.Context
	// 解析c的参数存入结构体以便使用
	userId     int64
	followId   int64
	actionType int
}

// NewProxyPostFollow 新建代理实例
func NewProxyPostFollow(c *gin.Context) *ProxyPostFollow {
	return &ProxyPostFollow{c: c}
}

// PostFollowPrepareNum 解析c的参数到代理结构体
func (p *ProxyPostFollow) PostFollowPrepareNum() error {
	// 根据上层中间件设置的user_id获得user_id
	rawUserId, _ := p.c.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析错误")
	}
	p.userId = userId
	// 根据URL中获得follow_id
	rawToUserId := p.c.Query("to_user_id")
	toUserId, err := strconv.ParseInt(rawToUserId, 10, 64)
	if err != nil {
		return errors.New("toUserId解析错误")
	}
	p.followId = toUserId
	// 根据URL中获得actionType
	rawFollowType := p.c.Query("action_type")
	followType, err := strconv.ParseInt(rawFollowType, 10, 32)
	if err != nil {
		return errors.New("actionType解析错误")
	}
	p.actionType = int(followType)
	return nil
}

func (p *ProxyPostFollow) PostFollowOk() {
	p.c.JSON(http.StatusOK, response.CommonResp{
		StatusCode: constant.SUCCESS,
		StatusMsg:  constant.SUCCESS_MESSAGE,
	})
}

func (p *ProxyPostFollow) PostFollowError(msg string) {
	p.c.JSON(http.StatusOK, response.CommonResp{
		StatusCode: constant.FAILURE,
		StatusMsg:  msg,
	})
}

func (p *ProxyPostFollow) DoPostFollow() error {
	err := p.PostFollowPrepareNum()
	if err != nil {
		return err
	}

	err = user_info.PostFollow(p.userId, p.followId, p.actionType)
	if err != nil {
		if errors.Is(err, user_info.ErrIvdAct) || errors.Is(err, user_info.ErrIvdFolUsr) || errors.Is(err, user_info.ErrIvdDel) {
			return err
		} else {
			return errors.New("请勿重复关注")
		}
	}
	return nil
}

func PostFollowHandler(c *gin.Context) {
	p := NewProxyPostFollow(c)

	err := p.DoPostFollow()
	if err != nil {
		p.PostFollowError(err.Error())
		return
	}
	p.PostFollowOk()
}
