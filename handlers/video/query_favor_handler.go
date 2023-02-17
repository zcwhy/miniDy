package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model"
	"miniDy/model/response"
	"miniDy/service/video"
	"miniDy/util"
	"net/http"
)

type QueryFavorListResponse struct {
	response.CommonResp
	Response *model.FavorListResponse
}

type ProxyQueryFavorListHandler struct {
	*gin.Context

	userId int64
}

func NewProxyQueryFavorListHandler(c *gin.Context) *ProxyQueryFavorListHandler {
	return &ProxyQueryFavorListHandler{Context: c}
}

func (p *ProxyQueryFavorListHandler) Do() {
	if err := p.parser(); err != nil {
		p.retError(err)
		return
	}
	favorVideoListRes, err := video.QueryFavorList(p.userId)
	if err != nil {
		p.retError(err)
		return
	}

	p.retOk(favorVideoListRes)
}

func QueryFavorListHandler(c *gin.Context) {
	NewProxyQueryFavorListHandler(c).Do()
}

func (p *ProxyQueryFavorListHandler) parser() error {
	rawUserId, ok := p.Get("user_id")
	if !ok {
		return errors.New("用户不存在")
	}
	var err error
	if p.userId, err = util.StringToInt64(fmt.Sprint(rawUserId)); err != nil {
		return err
	}
	return nil
}

func (p *ProxyQueryFavorListHandler) retOk(videoList *model.FavorListResponse) {
	p.JSON(http.StatusOK, QueryFavorListResponse{
		CommonResp: response.CommonResp{StatusCode: constant.SUCCESS},
		Response:   videoList,
	})
}

func (p *ProxyQueryFavorListHandler) retError(err error) {
	p.JSON(http.StatusOK, QueryFavorListResponse{
		CommonResp: response.CommonResp{StatusCode: constant.FAILURE, StatusMsg: err.Error()},
	})
}
