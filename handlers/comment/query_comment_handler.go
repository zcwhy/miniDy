package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model"
	"miniDy/model/response"
	"miniDy/service/comment"
	"miniDy/util"
	"net/http"
)

type QueryCommentListResponse struct {
	response.CommonResp
	Response *[]*model.Comment
}

type ProxyQueryCommentListHandler struct {
	*gin.Context
	videoId int64
}

func NewProxyQueryCommentListHandler(c *gin.Context) *ProxyQueryCommentListHandler {
	return &ProxyQueryCommentListHandler{Context: c}
}

func QueryCommentListHandler(c *gin.Context) {
	NewProxyQueryCommentListHandler(c).Do()
}

func (p *ProxyQueryCommentListHandler) Do() {
	//解析query
	if err := p.parser(); err != nil {
		p.retError(err)
		return
	}
	//调用service层QueryCommentList
	commentListRes, err := comment.QueryCommentList(p.videoId)
	if err != nil {
		p.retError(err)
		return
	}
	p.retOK(commentListRes)
}

func (p *ProxyQueryCommentListHandler) parser() error {
	//
	var err error
	p.videoId, err = util.StringToInt64(p.DefaultQuery("video_id", "0"))
	if err != nil {
		return errors.New("解析视频ID出错")
	}
	if p.videoId == 0 {
		return errors.New("视频不存在")
	}
	return nil
}

func (p *ProxyQueryCommentListHandler) retError(err error) {
	p.JSON(http.StatusOK, QueryCommentListResponse{
		CommonResp: response.CommonResp{StatusCode: constant.FAILURE, StatusMsg: err.Error()},
		Response:   nil,
	})
}

func (p *ProxyQueryCommentListHandler) retOK(commentList *comment.CommentList) {
	p.JSON(http.StatusOK, QueryCommentListResponse{
		CommonResp: response.CommonResp{StatusCode: constant.SUCCESS},
		Response:   &commentList.Comments,
	})
}
