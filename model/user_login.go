package model

import (
	"sync"
)

type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

type UserLoginDAO struct {
}

// 保证只执行一次
var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})

	return userLoginDao
}

func (s *UserLoginDAO) IsUserExist(username string) bool {
	user := &UserLogin{}
	exist := DB.Where("username = ?", username).Find(user).RowsAffected
	if exist == 1 {
		return true
	}
	return false
}
