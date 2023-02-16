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
	if err := video.PostFavor(p.videoId, p.actionType); err != nil {
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
	rawVideoId, exist := p.Get("video_id")
	if exist != true {
		return errors.New("videoID not exist")
	}
	p.videoId, err = util.StringToInt64(fmt.Sprint(rawVideoId))
	if err != nil {
		return errors.New("parse video_id error")
	}
	p.actionType, err = util.StringToInt64(p.Query("action_type"))
	if err != nil {
		return errors.New("parse action_type error")
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
