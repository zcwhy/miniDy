package video

import (
	"miniDy/model"
	"time"
)

type GetVideoFeedResponse struct {
	NextTime  int64          `json:"next_time"`
	VideoList []*model.Video `json:"video_list"`
}

type GetVideoFeedService struct {
	latestTime time.Time

	nextTime  int64
	videoList []*model.Video
}

func GetVideoFeed(latestTime int64) (*GetVideoFeedResponse, error) {
	return NewGetVideoFeedService(latestTime).Do()
}

func NewGetVideoFeedService(latestTime int64) *GetVideoFeedService {
	lt := time.Unix(latestTime, 0)

	return &GetVideoFeedService{latestTime: lt}
}

func (s *GetVideoFeedService) Do() (*GetVideoFeedResponse, error) {
	s.videoList = []*model.Video{}
	videoDao := model.NewVideoDao()

	err := videoDao.QueryVideosByTime(s.latestTime, &s.videoList)

	if err != nil {
		return nil, err
	}

	//如果返回列表没有数据， next_time设置为现在，否则设置为返回的第一条视频的创建时间(数据是以升序返回的)
	nextTime := time.Now().Unix()
	if len(s.videoList) != 0 {
		nextTime = s.videoList[0].CreatedAt.Unix()
	}

	//根据userInfoId填充作者信息
	for _, video := range s.videoList {
		userInfoId := video.UserInfoId
		userInfoDao := model.NewUserInfoDao()

		video.Author = model.UserInfo{}
		err := userInfoDao.QueryUserInfoById(userInfoId, &video.Author)

		if err != nil {
			return nil, err
		}
	}

	return &GetVideoFeedResponse{nextTime, s.videoList}, nil
}
