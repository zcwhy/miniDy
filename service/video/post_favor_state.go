package video

import (
	"fmt"
	"miniDy/constant"
	"miniDy/model"
)

type PostFavorFlow struct {
	userId     int64
	videoId    int64
	actionType int64
}

func NewPostFavorFlow(userId, videoId, actionType int64) *PostFavorFlow {
	return &PostFavorFlow{userId: userId, videoId: videoId, actionType: actionType}
}

func PostFavor(userId, videoId, actionType int64) error {
	return NewPostFavorFlow(userId, videoId, actionType).Do()
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
		return fmt.Errorf("未定义用户行为 %v", p.actionType)
	}

	return nil
}

func (p *PostFavorFlow) UpFavor() error {
	videoDao := model.NewVideoDao()
	userInfoDao := model.NewUserInfoDao()

	if err := videoDao.UpFavorByVideoId(p.userId, p.videoId); err != nil {
		return err
	}

	//点赞时，自己的喜欢数量 + 1
	if err := userInfoDao.AddUserFavoriteCount(p.userId); err != nil {
		return err
	}

	//被点赞视频的作者的总获赞数量 + 1
	user := &model.UserInfo{}

	if err := videoDao.QueryAuthorByVideoId(p.videoId, user); err != nil {
		return err
	}
	fmt.Println(p.videoId)

	if err := userInfoDao.AddUserTotalFavorite(user.Id); err != nil {
		return err
	}

	return nil
}

func (p *PostFavorFlow) DownFavor() error {
	videoDao := model.NewVideoDao()
	userInfoDao := model.NewUserInfoDao()
	if err := videoDao.DownFavorByVideoId(p.userId, p.videoId); err != nil {
		return err
	}

	//取消点赞时，自己的喜欢数量 - 1
	if err := userInfoDao.SubUserFavoriteCount(p.userId); err != nil {
		return err
	}

	//被取消点赞视频的作者的总获赞数量 - 1
	user := &model.UserInfo{}

	if err := videoDao.QueryAuthorByVideoId(p.videoId, user); err != nil {
		return err
	}

	if err := userInfoDao.SubUserTotalFavorite(user.Id); err != nil {
		return err
	}

	return nil
}
