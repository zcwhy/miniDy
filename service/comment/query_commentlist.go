package comment

import (
	"errors"
	"fmt"
	"miniDy/model"
	"miniDy/util"
)

type CommentList struct {
	Comments []*model.Comment
}

type QueryCommentListFlow struct {
	videoId int64

	commentList CommentList
	Response    *CommentList
}

func QueryCommentList(videoId int64) (*CommentList, error) {
	return NewQueryCommentListFlow(videoId).Do()
}

func NewQueryCommentListFlow(videoId int64) *QueryCommentListFlow {
	return &QueryCommentListFlow{videoId: videoId}
}

func (q *QueryCommentListFlow) Do() (*CommentList, error) {
	if err := q.check(); err != nil {
		return nil, err
	}
	if err := q.prepare(); err != nil {
		return nil, err
	}
	if err := q.pack(); err != nil {
		return nil, err
	}
	return q.Response, nil
}

func (q *QueryCommentListFlow) check() error {
	if !model.NewVideoDao().IsVideoExistById(q.videoId) {
		return fmt.Errorf("视频%d不存在或已经被删除", q.videoId)
	}
	return nil
}

func (q *QueryCommentListFlow) prepare() error {
	return model.NewCommentDAO().QueryCommentListByVideoId(q.videoId, &q.commentList.Comments)
}

func (q *QueryCommentListFlow) pack() error {
	if q.commentList.Comments == nil {
		return errors.New("the video have not comment")
	}
	for _, v := range q.commentList.Comments {
		util.DateFormat(v)
	}
	q.Response = &q.commentList
	return nil
}
