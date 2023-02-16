package comment

import (
	"errors"
	"fmt"
	"miniDy/constant"
	"miniDy/model"
	"miniDy/util"
	"time"
)

type Response struct {
	MyComment *model.Comment `json:"comment"`
}

type PostCommentFlow struct {
	userId      int64
	videoId     int64
	commentId   int64
	actionType  int64
	commentText string

	comment *model.Comment

	*Response
}

func PostComment(userId, videoId, commentId, actionType int64, commentText string) (*Response, error) {
	return NewPostCommentFlow(userId, videoId, commentId, actionType, commentText).Do()
}

func NewPostCommentFlow(userId, videoId, commentId, actionType int64, commentText string) *PostCommentFlow {
	return &PostCommentFlow{userId: userId, videoId: videoId, commentId: commentId, actionType: actionType, commentText: commentText}
}

func (p *PostCommentFlow) Do() (*Response, error) {
	if err := p.check(); err != nil {
		return nil, err
	}
	if err := p.prepare(); err != nil {
		return nil, err
	}
	if err := p.pack(); err != nil {
		return nil, err
	}
	return p.Response, nil
}

func (p *PostCommentFlow) CreateComment() error {

	p.comment = &model.Comment{
		UserInfoId: p.userId,
		VideoId:    p.videoId,
		Content:    p.commentText,
		CreatedAt:  time.Time{},
	}
	commentDAO := model.NewCommentDAO()
	if err := commentDAO.AddComment(p.comment); err != nil {
		return err
	}
	return nil
}

func (p *PostCommentFlow) DeleteComment() error {
	if err := model.NewCommentDAO().DeleteComment(p.comment); err != nil {
		return err
	}
	return nil
}

func (p *PostCommentFlow) check() error {
	if !model.NewUserInfoDao().IsUserExistById(p.userId) {
		return fmt.Errorf("用户%d处于登出状态", p.userId)
	}
	if !model.NewVideoDao().IsVideoExistById(p.videoId) {
		return fmt.Errorf("视频%d不存在或已经被删除", p.videoId)
	}
	return nil
}

func (p *PostCommentFlow) prepare() error {
	switch p.actionType {
	case constant.CREATE:
		return p.CreateComment()
	case constant.DELETE:
		return p.DeleteComment()
	default:
		return fmt.Errorf("undefined action %v", p.actionType)
	}
}

func (p *PostCommentFlow) pack() error {
	if p.comment == nil {
		return errors.New("comment is null-ptr")
	}
	userInfo := model.UserInfo{}
	if err := model.NewUserInfoDao().QueryUserInfoById(p.comment.UserInfoId, &userInfo); err != nil {
		return err
	}
	p.comment.User = userInfo
	util.DateFormat(p.comment)
	p.Response = &Response{MyComment: p.comment}
	return nil
}
