package comment

import (
	"errors"
	"miniDy/model"
	"miniDy/util"
)

/*
从handler处拿到收集的参数
1、构造流对象， 执行流对象的do方法
2、参数校验
3、收集dao层查询的数据，封装成返回对象返回给handler层
*/
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

func (p *QueryCommentListFlow) Do() (*CommentList, error) {
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

func (p *QueryCommentListFlow) check() error {
	//if !models.NewVideoDAO().IsVideoExistById(q.videoId) {
	//	return fmt.Errorf("视频%d不存在或已经被删除", q.videoId)
	//}
	return nil
}

func (p *QueryCommentListFlow) prepare() error {
	return model.NewCommentDAO().QueryCommentListByVideoId(p.videoId, &p.commentList.Comments)
}

func (p *QueryCommentListFlow) pack() error {
	if p.commentList.Comments == nil {
		return errors.New("the video have not comment")
	}
	for _, v := range p.commentList.Comments {
		util.DateFormat(v)
	}
	p.Response = &p.commentList
	return nil
}
