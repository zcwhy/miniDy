package video

import (
	"errors"
	"miniDy/model"
)

type List struct {
	Video []*model.Video `json:"video_list,omitempty"`
}

type QueryVideoListByUserIdFlow struct {
	userId int64
	videos []*model.Video

	videoList *List
}

func NewQueryVideoListByUserIdFlow(userId int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{userId: userId}
}

func (q *QueryVideoListByUserIdFlow) checkNum() error {
	if !model.NewUserInfoDao().IsUserExistById(q.userId) { // 调用下层的models层查询数据库中是否存在
		return errors.New("用户不存在")
	}
	return nil
}

func (q *QueryVideoListByUserIdFlow) packDate() error {
	err := model.NewVideoDao().QueryVideoListByUserId(q.userId, &q.videos)
	if err != nil {
		return err
	}
	//作者信息查询
	var userInfo model.UserInfo
	err = model.NewUserInfoDao().QueryUserInfoById(q.userId, &userInfo) // 调用下层的models层根据userid查询用户信息
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段)
	for i := range q.videos {
		q.videos[i].Author = userInfo
		q.videos[i].IsFavorite = model.NewVideoDao().IsUserFavorVideoExist(q.userId, q.videos[i].Id) // 根据userId和videoId进行查询
	}

	q.videoList = &List{Video: q.videos}

	return nil
}

func (q *QueryVideoListByUserIdFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packDate(); err != nil {
		return nil, err
	}

	return q.videoList, nil
}

func QueryVideoListByUserId(userId int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(userId).Do()
}
