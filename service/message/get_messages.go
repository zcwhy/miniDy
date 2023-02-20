package message

import (
	"errors"
	"miniDy/constant"
	"miniDy/model"
	"time"
)

type GetMessageRecordsService struct {
	userId   int64
	toUserId int64
	lastTime time.Time

	messageList []*model.Message
}

func GetMessageRecords(userId, toUserId int64, lastTime int64) ([]*model.Message, error) {
	service := &GetMessageRecordsService{
		userId:   userId,
		toUserId: toUserId,
		lastTime: time.Unix(lastTime/1000, 0),
	}
	err := service.Do()

	return service.messageList, err
}

func (s GetMessageRecordsService) Do() error {
	userInfoDao := model.NewUserInfoDao()
	messageDao := model.NewMessageDao()

	if exist := userInfoDao.IsUserExistById(s.toUserId); exist == false {
		return errors.New("发送用户id错误")
	}

	s.messageList = make([]*model.Message, constant.MAX_MESSAGE_NUMBER)
	if err := messageDao.QueryMessages(s.userId, s.toUserId, s.lastTime, &s.messageList); err != nil {
		return err
	}
	return nil
}
