package video

import (
	"fmt"
	"miniDy/constant"
	"miniDy/model"
)

type PostFavorFlow struct {
	videoId    int64
	actionType int64
}

func NewPostFavorFlow(videoId, actionType int64) *PostFavorFlow {
	return &PostFavorFlow{videoId: videoId, actionType: actionType}
}

func PostFavor(videoId, actionType int64) error {
	return NewPostFavorFlow(videoId, actionType).Do()
}

func (p *PostFavorFlow) Do() error {
	if !model.NewVideoDao().IsVideoExistById(p.videoId) {
		return fmt.Errorf("视频%d不存在或已经被删除", p.videoId)
	}

	switch p.actionType {
	case constant.UPCLICK:
		if err := p.UpFavor(); err != nil {
			return err
		}
	case constant.DOWNCLICK:
		if err := p.DownFavor(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("undefined action_type %v", p.actionType)
	}

	return nil
}

func (p *PostFavorFlow) UpFavor() error {
	if err := model.NewVideoDao().UpFavorByVideoId(p.videoId); err != nil {
		return err
	}
	return nil
}

func (p *PostFavorFlow) DownFavor() error {
	if err := model.NewVideoDao().DownFavorByVideoId(p.videoId); err != nil {
		return err
	}
	return nil
}
