package video

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"net/http"
)

type QueryFavorListResponse struct {
	response.CommonResp
}

type ProxyQueryFavorListHandler struct {
	*gin.Context

	videoId    int64
	actionType int64
}

func NewProxyQueryFavorListHandler(c *gin.Context) *ProxyQueryFavorListHandler {
	return &ProxyQueryFavorListHandler{Context: c}
}

func (p *ProxyQueryFavorListHandler) Do() {

}

func QueryFavorListHandler(c *gin.Context) {
	NewProxyQueryFavorListHandler(c).Do()
}

func (p *ProxyQueryFavorListHandler) parser() error {

	return nil
}

func (p *ProxyQueryFavorListHandler) retOk() {
	p.JSON(http.StatusOK, PostFavorResponse{
		CommonResp: response.CommonResp{StatusCode: constant.SUCCESS},
	})
}

func (p *ProxyQueryFavorListHandler) retError(err error) {
	p.JSON(http.StatusOK, PostFavorResponse{
		CommonResp: response.CommonResp{StatusCode: constant.FAILURE, StatusMsg: err.Error()},
	})
}
