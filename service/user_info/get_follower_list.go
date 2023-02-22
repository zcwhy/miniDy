package user_info

import (
	"errors"
	"miniDy/model"
)

type UserFollowerList struct {
	FollowList []*model.UserInfo `json:"user_list"`
}

func GetFollowerList(userId int64) (*UserFollowerList, error) {
	userInfoDao := model.NewUserInfoDao()

	if exist := userInfoDao.IsUserExistById(userId); !exist {
		return nil, errors.New("用户不存在")
	}

	var followerList []*model.UserInfo

	if err := userInfoDao.QueryFollowerListById(userId, &followerList); err != nil {
		return nil, err
	}

	for _, follower := range followerList {
		isFollow := userInfoDao.IsFollowExist(userId, follower.Id)
		follower.IsFollow = isFollow == true
	}

	return &UserFollowerList{FollowList: followerList}, nil
}
