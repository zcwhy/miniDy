package model

import (
	"errors"
	"sync"
)

type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	Name          string      `json:"name" gorm:"name,omitempty"`
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *UserLogin  `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                     //用户与评论的一对多
}

type UserInfoDAO struct {
}

var (
	UserInfoDao  *UserInfoDAO
	UserInfoOnce sync.Once
)

func NewUserInfoDao() *UserInfoDAO {
	UserInfoOnce.Do(func() {
		UserInfoDao = new(UserInfoDAO)
	})

	return UserInfoDao
}

func (s *UserInfoDAO) UserRegister(info *UserInfo) error {
	return DB.Create(info).Error
}

func (s *UserInfoDAO) QueryUserInfoById(id int64, info *UserInfo) error {
	if info == nil {
		errors.New("QueryUserInfoById UserInfo 空指针")
	}
	return DB.Model(&UserInfo{}).Where("id = ?", id).Find(info).Error
}

func (s *UserInfoDAO) IsUserExistById(userId int64) bool {
	return DB.Find(&UserInfo{}, userId).RowsAffected == 1
}

func (s *UserInfoDAO) QueryFollowListById(id int64, followList *[]*UserInfo) error {
	if followList == nil {
		errors.New("QueryFollowListById followList 空指针")
	}
	return DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.user_info_id = ? AND r.follow_id = u.id", id).Scan(followList).Error
}

func (s *UserInfoDAO) QueryFollowerListById(id int64, followerList *[]*UserInfo) error {
	if followerList == nil {
		errors.New("QueryFollowerListById followList 空指针")
	}
	return DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.follow_id = ? AND r.user_info_id = u.id", id).Scan(followerList).Error
}
