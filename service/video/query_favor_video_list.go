package video

import (
	"miniDy/model"
)

type QueryFavorListFlow struct {
	userId int64
}

func NewQueryFavorListFlow(userId int64) *QueryFavorListFlow {
	return &QueryFavorListFlow{userId}
}

func QueryFavorList(userId int64) (*[]*model.Video, error) {
	return NewQueryFavorListFlow(userId).Do()
}

func (q *QueryFavorListFlow) Do() (*[]*model.Video, error) {
	FavorVideoList, err := model.NewVideoDao().QueryFavorListByUserId(q.userId)
	if err != nil {
		return nil, err
	}
	if FavorVideoList != nil {
		for _, v := range *FavorVideoList {
			user := model.UserInfo{}
			if err = model.NewVideoDao().QueryAuthorByVideoId(v.Id, &user); err != nil {
				return nil, err
			}
			v.Author = user
		}
	}
	return FavorVideoList, nil
}
