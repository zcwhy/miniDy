package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
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

var (
	ErrIvdPtr        = errors.New("空指针错误")
	ErrEmptyUserList = errors.New("用户列表为空")
)

func NewUserInfoDao() *UserInfoDAO {
	UserInfoOnce.Do(func() {
		UserInfoDao = new(UserInfoDAO)
	})

	return UserInfoDao
}

func (s *UserInfoDAO) UserRegister(info *UserInfo) error {
	_ = DB.Create(info).Error
	return nil
}

// <<<<<<< HEAD
//
//	func (s *UserInfoDAO) QueryUserInfoById(id int64, info *UserInfo) error {
//		if info == nil {
//			errors.New("QueryUserInfoById UserInfo 空指针")
//		}
//		return DB.Model(&UserInfo{}).Where("id = ?", id).Find(info).Error
func (s *UserInfoDAO) QueryUserInfoById(userId int64, userInfo *UserInfo) error {
	if userInfo == nil {
		return ErrIvdPtr
	}
	DB.Where("id=?", userId).Select([]string{"id", "name", "follow_count", "follower_count", "is_follow"}).First(userInfo)
	if userInfo.Id == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func (s *UserInfoDAO) IsUserExistById(userId int64) bool {
	var userinfo UserInfo
	if err := DB.Where("id=?", userId).Select("id").First(&userinfo).Error; err != nil {
		log.Println(err)
	}
	if userinfo.Id == 0 {
		return false
	}
	return true
}

func (s *UserInfoDAO) AddUserFollow(userId int64, followId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count+1 WHERE id = ?", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO user_relations (`user_info_id`,`follow_id`) VALUES (?,?)", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *UserInfoDAO) CancelUserFollow(userId int64, followId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM user_relations WHERE user_info_id=? AND follow_id=?", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *UserInfoDAO) IsFollowExist(userId int64, followId int64) bool {
	var userinfo UserInfo
	exist := DB.Raw("SELECT r.* from user_relations r WHERE r.user_info_id = ? AND r.follow_id = ?", userId, followId).Scan(userinfo).RowsAffected
	//log.Printf("########**%#v", exist)
	if exist == 1 {
		return true
	}
	return false
}
