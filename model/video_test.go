package model

import (
	"fmt"
	"testing"
	"time"
)

func TestVideoDAO_QueryVideosByTime(t *testing.T) {
	VideoDao := NewVideoDao()

	var timeStamp int64
	timeStamp = 1675688656245

	lastTime := time.Unix(timeStamp/int64(1000), 0)
	var videos []*Video

	err := VideoDao.QueryVideosByTime(lastTime, &videos)

	fmt.Println(len(videos))

	if err != nil {
		fmt.Println(err)
	}

	for _, i := range videos {
		fmt.Println(i.Title)
	}
}

func TestVideoDao_IsUserFavorVideoExist(t *testing.T) {
	videoDao := NewVideoDao()
	exist := videoDao.IsUserFavorVideoExist(1, 1)
	if exist == true {
		t.Error("favor not exist but got true")
	}
}
