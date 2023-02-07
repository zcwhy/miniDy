package model

import "testing"

func TestVideoDao_IsUserFavorVideoExist(t *testing.T) {
	videoDao := NewVideoDao()
	exist := videoDao.IsUserFavorVideoExist(1, 1)
	if exist == true {
		t.Error("favor not exist but got true")
	}
}
