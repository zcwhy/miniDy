package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
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

type userFavorVideo struct {
	userInfoId int64
	videoId    int64
}

type FavorListResponse struct {
	FavorVideoList []*Video `json:"video_list"`
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

func (v *VideoDAO) IsVideoExistById(videoId int64) bool {
	newVideo := Video{}
	if err := DB.First(&newVideo, videoId); err != nil {
		log.Println(err)
	}
	if newVideo.Id == 0 {
		return false
	}
	return true
}

func (v *VideoDAO) UpFavorByVideoId(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Model(&Video{}).Where("id = ?", videoId).
			Update("favorite_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			return err
		}

		ufv := &userFavorVideo{userInfoId: userId, videoId: videoId}
		if v.IsUserFavorVideoExist(userId, videoId) == false {
			err = tx.Table("user_favor_videos").Create(ufv).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (v *VideoDAO) DownFavorByVideoId(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Model(&Video{}).Where("id = ?", videoId).
			Update("favorite_count", gorm.Expr("comment_count + ?", -1)).Error
		if err != nil {
			return err
		}
		ufv := &userFavorVideo{userInfoId: userId, videoId: videoId}
		if v.IsUserFavorVideoExist(userId, videoId) == true {
			err = tx.Table("user_favor_videos").Delete(ufv).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (v *VideoDAO) QueryFavorListByUserId(userId int64) (*FavorListResponse, error) {
	var list []*Video
	var fvList []*userFavorVideo
	var err error
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&userFavorVideo{}).Where("user_info_id = ?", userId).Find(&fvList).Error; err != nil {
			return err
		}
		for _, v := range fvList {
			var favorVideo *Video
			if err = tx.Model(&Video{}).First(favorVideo, v.videoId).Error; err != nil {
				return err
			}
			list = append(list, favorVideo)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &FavorListResponse{FavorVideoList: list}, nil
}

func (v *VideoDAO) QueryAuthorByVideoId(videoId int64, author *UserInfo) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		newVideo := Video{}
		if err := tx.Model(&Video{}).First(&newVideo, videoId).Error; err != nil {
			return err
		}
		if err := tx.Model(&UserInfo{}).First(author, newVideo.UserInfoId).Error; err != nil {
			return err
		}
		return nil
	})
}
