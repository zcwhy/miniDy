package user_info

import (
	"errors"
	"miniDy/constant"
	"miniDy/model"
)

type FriendUser struct {
	User    *model.UserInfo `json:",inline"`
	Message string
	MsgType int64 `json:"msg_type"`
}

type UserFriendList struct {
	FriendList []*FriendUser `json:"user_list"`
}

func GetFriendList(userId int64) (*UserFriendList, error) {
	userInfoDao := model.NewUserInfoDao()

	if exist := userInfoDao.IsUserExistById(userId); !exist {
		return nil, errors.New("用户不存在")
	}

	var friends []*model.UserInfo

	if err := userInfoDao.QueryFriendListById(userId, &friends); err != nil {
		return nil, err
	}

	var friendList []*FriendUser
	messageDao := model.NewMessageDao()

	for _, friend := range friends {
		friend.IsFollow = true

		message := &model.Message{}
		err := messageDao.QueryLatestMessageById(userId, friend.Id, message)

		if err != nil {
			return nil, err
		}

		friendUser := &FriendUser{User: friend, Message: message.Content}
		if message.FromUserId != userId {
			friendUser.MsgType = constant.TO_MESSAGE
		}

		friendList = append(friendList, friendUser)
	}

	return &UserFriendList{FriendList: friendList}, nil
}
