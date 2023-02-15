package model

import (
	"errors"
	"miniDy/constant"
	"sync"
	"time"
)

type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

type VideoDAO struct {
}

var (
	VideoDao     *VideoDAO
	VideoDaoOnce sync.Once
)

func NewVideoDao() *VideoDAO {
	VideoDaoOnce.Do(func() {
		VideoDao = new(VideoDAO)
	})
	return VideoDao
}

func (*VideoDAO) CountUserVideoById(userId int64) (int64, error) {
	result := DB.Model(&Video{}).Where("user_info_id=?", userId)
	return result.RowsAffected, result.Error
}

func (*VideoDAO) CreateVideo(userId int64, playUrl string, title string) error {
	return DB.Create(&Video{
		UserInfoId: userId,
		PlayUrl:    playUrl,
		Title:      title,
	}).Error
}

func (*VideoDAO) QueryVideosByTime(latestTime time.Time, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideosByTime videoList 空指针")
	}
	return DB.Model(&Video{}).Where("created_at >= ?", latestTime).
		Order("created_at ASC").
		Limit(constant.MAX_VIDEO_NUMBER).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).
		Find(videoList).Error
}

func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return DB.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title"}).
		Find(videoList).Error
}

func (v *VideoDAO) IsUserFavorVideoExist(userId int64, videoId int64) bool {
	userFavorVideo := &Video{}
	exist := DB.Raw("SELECT f.* from user_favor_videos f WHERE f.user_info_id = ? AND f.video_id = ?", userId, videoId).Scan(userFavorVideo).RowsAffected
	if exist == 1 {
		return true
	}
	return false
}
