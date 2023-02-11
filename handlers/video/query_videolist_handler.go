package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/video"
	"net/http"
)

// ProxyQueryVideoList 代理结构体
type ProxyQueryVideoList struct {
	c *gin.Context
}

// NewProxyQueryVideoList 创建代理结构体
func NewProxyQueryVideoList(c *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{c: c}
}

type ListResponse struct {
	response.CommonResp
	*video.List
}

// QueryVideoListOk 前端返回Ok
func (p *ProxyQueryVideoList) QueryVideoListOk(videoList *video.List) {
	p.c.JSON(http.StatusOK, ListResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		List: videoList,
	})
}

// QueryVideoListError 给前端返回Error
func (p *ProxyQueryVideoList) QueryVideoListError(msg string) {
	p.c.JSON(http.StatusOK, ListResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.FAILURE,
			StatusMsg:  msg,
		},
	})
}

func (p *ProxyQueryVideoList) DoQueryVideoListByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}

	videoList, err := video.QueryVideoListByUserId(userId)
	if err != nil {
		return err
	}
	p.QueryVideoListOk(videoList)
	return nil
}

func QueryVideoListHandler(c *gin.Context) {
	p := NewProxyQueryVideoList(c)
	rawId, ok := c.Get("user_id")
	if !ok {
		p.QueryVideoListError("解析userId出错")
		return
	}
	err := p.DoQueryVideoListByUserId(rawId)
	if err != nil {
		p.QueryVideoListError(err.Error())
	}
}
