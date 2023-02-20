package user_info

import (
	"errors"
	"miniDy/model"
)

type UserFollowList struct {
	FollowList []*model.UserInfo `json:"user_list"`
}

func GetFollowList(userId int64) (*UserFollowList, error) {
	userInfoDao := model.NewUserInfoDao()

	if exist := userInfoDao.IsUserExistById(userId); !exist {
		return nil, errors.New("用户不存在")
	}

	var followList []*model.UserInfo

	if err := userInfoDao.QueryFollowListById(userId, &followList); err != nil {
		return nil, err
	}

	for _, user := range followList {
		user.IsFollow = true
	}

	return &UserFollowList{FollowList: followList}, nil
}
