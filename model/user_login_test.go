package model

import (
	"testing"
)

func TestUserLoginDao_IsUserExist(t *testing.T) {
	userLoginDao := NewLoginDao()
	exist := userLoginDao.IsUserExist("zcwhy333")

	if exist == true {
		t.Error("user:zcwhy not exist but got true")
	}
}
