package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/video"
	"miniDy/util"
	"net/http"
)

type PostFavorResponse struct {
	response.CommonResp
}

type ProxyPostFavorHandler struct {
	*gin.Context

	userId     int64
	videoId    int64
	actionType int64
}

func NewProxyFavorHandler(c *gin.Context) *ProxyPostFavorHandler {
	return &ProxyPostFavorHandler{Context: c}
}

func (p *ProxyPostFavorHandler) Do() {
	if err := p.parser(); err != nil {
		p.retError(err)
		return
	}
	if err := video.PostFavor(p.userId, p.videoId, p.actionType); err != nil {
		p.retError(err)
		return
	}
	p.retOk()
}

func PostFavorHandler(c *gin.Context) {
	NewProxyFavorHandler(c).Do()
}

func (p *ProxyPostFavorHandler) parser() error {
	var err error
	rawVideoId := p.DefaultQuery("video_id", "0")
	p.videoId, err = util.StringToInt64(fmt.Sprint(rawVideoId))
	if err != nil {
		return errors.New("解析视频ID出错")
	}

	rawUserId := p.DefaultQuery("user_id", "0")
	p.userId, err = util.StringToInt64(fmt.Sprint(rawUserId))
	if err != nil {
		return errors.New("解析用户ID出错")
	}

	p.actionType, err = util.StringToInt64(p.DefaultQuery("action_type", "0"))
	if err != nil {
		return errors.New("解析用户操作出错")
	}
	return nil
}

func (p *ProxyPostFavorHandler) retOk() {
	p.JSON(http.StatusOK, PostFavorResponse{
		CommonResp: response.CommonResp{StatusCode: constant.SUCCESS},
	})
}

func (p *ProxyPostFavorHandler) retError(err error) {
	p.JSON(http.StatusOK, PostFavorResponse{
		CommonResp: response.CommonResp{StatusCode: constant.FAILURE, StatusMsg: err.Error()},
	})
}
