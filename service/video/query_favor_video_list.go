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

func QueryFavorList(userId int64) (*model.FavorListResponse, error) {
	return NewQueryFavorListFlow(userId).Do()
}

func (q *QueryFavorListFlow) Do() (*model.FavorListResponse, error) {
	videoDao := model.NewVideoDao()
	//FavorVideoList := &model.FavorListResponse{}
	FavorVideoList, err := videoDao.QueryFavorListByUserId(q.userId)
	if err != nil {
		return nil, err
	}
	if FavorVideoList != nil {
		for _, v := range FavorVideoList.FavorVideoList {
			err := videoDao.QueryAuthorByVideoId(v.Id, &v.Author)
			if err != nil {
				return nil, err
			}
		}
	}
	return FavorVideoList, nil
}
