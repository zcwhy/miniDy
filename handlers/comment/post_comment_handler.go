package comment

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/comment"
	"miniDy/util"
	"net/http"
)

type PostCommentResponse struct {
	response.CommonResp
	*comment.Response
}

type ProxyPostCommentHandler struct {
	*gin.Context

	videoId     int64
	userId      int64
	commentId   int64
	actionType  int64
	commentText string
}

func NewProxyPostCommentHandler(c *gin.Context) *ProxyPostCommentHandler {
	return &ProxyPostCommentHandler{Context: c}
}

func PostCommentHandler(c *gin.Context) {
	NewProxyPostCommentHandler(c).Do()
}

func (p *ProxyPostCommentHandler) Do() {
	//解析query
	if err := p.parser(); err != nil {
		p.retError(err)
		return
	}
	//调用service层PostComment
	commentRes, err := comment.PostComment(p.userId, p.videoId, p.commentId, p.actionType, p.commentText)

	if err != nil {
		p.retError(err)
		return
	}

	p.retOK(commentRes)
}

func (p *ProxyPostCommentHandler) parser() error {
	//
	var err error = nil
	p.videoId, err = util.StringToInt64(p.DefaultQuery("video_id", "0"))
	if err != nil {
		return errors.New("解析视频ID出错")
	}
	if p.videoId == 0 {
		return errors.New("视频不存在")
	}

	p.userId, err = util.StringToInt64(p.DefaultQuery("user_id", "0"))
	if err != nil {
		return errors.New("解析用户ID出错")
	}
	p.actionType, err = util.StringToInt64(p.Query("action_type"))
	if err != nil {
		return errors.New("解析用户操作出错")
	}
	//
	switch p.actionType {
	case constant.CREATE:
		p.commentText = p.Query("comment_text")
	case constant.DELETE:
		p.commentId, err = util.StringToInt64(p.Query("comment_id"))
		if err != nil {
			return errors.New("解析评论ID出错")
		}
	default:
		return fmt.Errorf("未定义用户行为 %v", p.actionType)
	}
	return nil
}

func (p *ProxyPostCommentHandler) retError(err error) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResp: response.CommonResp{StatusCode: constant.FAILURE, StatusMsg: err.Error()},
		Response:   nil,
	})
}

func (p *ProxyPostCommentHandler) retOK(comment *comment.Response) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResp: response.CommonResp{StatusCode: constant.SUCCESS},
		Response:   comment,
	})
}
