package model

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type UserInfo struct {
	Id              int64       `json:"id" gorm:"id,omitempty"`                             // 用户id
	Name            string      `json:"name" gorm:"name,omitempty"`                         // 用户名称
	FollowCount     int64       `json:"follow_count" gorm:"follow_count,omitempty"`         // 关注总数
	FollowerCount   int64       `json:"follower_count" gorm:"follower_count,omitempty"`     // 粉丝总数
	IsFollow        bool        `json:"is_follow" gorm:"is_follow,omitempty"`               // true-已关注，false-未关注
	Avatar          string      `json:"avatar" gorm:"avatar,omitempty"`                     //用户头像
	BackgroundImage string      `json:"background_image" gorm:"background_image,omitempty"` //用户个人页顶部大图
	Signature       string      `json:"signature" gorm:"signature,omitempty"`               //个人简介
	TotalFavorited  int64       `json:"total_favorited" gorm:"total_favorited,omitempty"`   //获赞数量
	WorkCount       int64       `json:"work_count" gorm:"work_count,omitempty"`             //作品数量
	FavoriteCount   int64       `json:"favorite_count" gorm:"favorite_count,omitempty"`     //点赞数量
	User            *UserLogin  `json:"-"`                                                  //用户与账号密码之间的一对一
	Videos          []*Video    `json:"-"`                                                  //用户与投稿视频的一对多
	Follows         []*UserInfo `json:"-" gorm:"many2many:user_relations;"`                 //用户之间的多对多
	FavorVideos     []*Video    `json:"-" gorm:"many2many:user_favor_videos;"`              //用户与点赞视频之间的多对多
	Comments        []*Comment  `json:"-"`                                                  //用户与评论的一对多
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
	if err := DB.Create(info).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserInfoDAO) QueryUserInfoById(userId int64, userInfo *UserInfo) error {
	if userInfo == nil {
		return errors.New("空指针错误")
	}
	DB.Where("id=?", userId).First(userInfo)
	if userInfo.Id == 0 {
		return errors.New("用户不存在")
	}
	return nil
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

// AddUserWorkCount 在用户发布作品时，用户的作品数量 + 1
func (s *UserInfoDAO) AddUserWorkCount(userId int64) error {
	return DB.Model(&UserInfo{}).Where("id = ?", userId).
		Update("work_count", gorm.Expr(" work_count + ?", 1)).Error
}

// AddUserTotalFavorite 在用户发布作品被点赞时，用户的获赞总量 + 1
func (s *UserInfoDAO) AddUserTotalFavorite(userId int64) error {
	return DB.Model(&UserInfo{}).Where("id = ?", userId).
		Update("total_favorited", gorm.Expr(" total_favorited + ?", 1)).Error
}

// SubUserTotalFavorite 在用户发布的视频被取消点赞时，用户的获赞总量 - 1
func (s *UserInfoDAO) SubUserTotalFavorite(userId int64) error {
	return DB.Model(&UserInfo{}).Where("id = ?", userId).
		Update("total_favorited", gorm.Expr(" total_favorited - ?", 1)).Error
}

// AddUserFavoriteCount 在用户点赞视频时，用户喜欢的作品的数量 + 1
func (s *UserInfoDAO) AddUserFavoriteCount(userId int64) error {
	return DB.Model(&UserInfo{}).Where("id = ?", userId).
		Update("favorite_count", gorm.Expr(" favorite_count + ?", 1)).Error
}

// SubUserFavoriteCount 在用户取消点赞视频时，用户喜欢的作品的数量 - 1
func (s *UserInfoDAO) SubUserFavoriteCount(userId int64) error {
	return DB.Model(&UserInfo{}).Where("id = ?", userId).
		Update("favorite_count", gorm.Expr(" favorite_count - ?", 1)).Error
}

func (s *UserInfoDAO) IsFollowExist(userId int64, followId int64) bool {
	var userinfo UserInfo
	exist := DB.Raw("SELECT r.* from user_relations r WHERE r.user_info_id = ? AND r.follow_id = ?", userId, followId).Scan(userinfo).RowsAffected
	//log.Printf("########**%#v", exist)
	return exist == 1
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

func (s *UserInfoDAO) QueryFriendListById(id int64, followerList *[]*UserInfo) error {
	if followerList == nil {
		errors.New("QueryFollowerListById followList 空指针")
	}
	return DB.Raw("SELECT u.* FROM user_relations r1, user_relations r2, user_infos u  WHERE r1.user_info_id = r2.follow_id "+
		"AND  r1.follow_id = r2.user_info_id "+
		"AND r1.user_info_id = ? AND r1.follow_id = u.id", id).Scan(followerList).Error
}
