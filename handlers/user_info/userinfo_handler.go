package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model"
	"miniDy/model/response"
	"net/http"
)

// 返回相应结构体
type UserResponse struct {
	response.CommonResp
	User *model.UserInfo `json:"user"`
}

// UserInfo代理结构体
type ProxyUserInfo struct {
	c *gin.Context
}

// 新建代理实实例
func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

// 前端返回错误
func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  msg,
		},
	})
}

// 前端返回正确
func (p *ProxyUserInfo) UserInfoOk(user *model.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		User: user,
	})
}

// 由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("解析userId失败")
	}

	//var unserInfo model.NewUserInfoDAO()

	var userInfo model.UserInfo
	err := model.NewUserInfoDao().QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

// 用户信息接口Handler
func UserInfoHandler(c *gin.Context) {
	p := NewProxyUserInfo(c) // 新建代理实实例

	// 得到上层中间件根据token解析的userId*******
	rawId, ok := c.Get("user_id")
	if !ok {
		p.UserInfoError("解析userId出错")
		return
	}
	err := p.DoQueryUserInfoByUserId(rawId)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}
