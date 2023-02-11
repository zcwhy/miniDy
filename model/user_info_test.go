package model

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitDB()
	code := m.Run()
	os.Exit(code)
}

func TestUserInfoDAO_UserRegister(t *testing.T) {
	userLoginDao := NewLoginDao()
	userInfoDao := NewUserInfoDao()

	userLogin := &UserLogin{Username: "zcwhy", Password: "123456"}
	user := &UserInfo{User: userLogin, Name: userLogin.Username}

	exist := userLoginDao.IsUserExist("zcwhy")

	if exist == true {
		t.Error("user:zcwhy not exist but got true")
	}

	err := userInfoDao.UserRegister(user)
	if err != nil {
		return
	}

	exist = userLoginDao.IsUserExist("zcwhy")

	if exist == false {
		t.Error("user:zcwhy should exist but got false")
	}
}

func TestUserInfoDAO_QueryUserInfoById(t *testing.T) {
	userInfoDao := NewUserInfoDao()
	userInfo := &UserInfo{}

	_ = userInfoDao.QueryUserInfoById(7, userInfo)

	if userInfo.Id != 7 {
		t.Errorf("expected userId is 7 but got %d", userInfo.Id)
	}
}
