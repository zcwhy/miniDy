package message

import (
	"errors"
	"miniDy/model"
)

type GetMessageRecordsService struct {
	fromUserId int64
	toUserId   int64
	lastTime   int64

	messageList []*model.Message
}

func GetMessageRecords(userId, toUserId int64, lastTime int64) ([]*model.Message, error) {
	service := &GetMessageRecordsService{
		fromUserId: userId,
		toUserId:   toUserId,
		lastTime:   lastTime,
	}
	err := service.Do()
	return service.messageList, err
}

func (s *GetMessageRecordsService) Do() error {
	userInfoDao := model.NewUserInfoDao()
	messageDao := model.NewMessageDao()

	if exist := userInfoDao.IsUserExistById(s.toUserId); exist == false {
		return errors.New("发送用户id错误")
	}

	s.messageList = make([]*model.Message, 0)
	if err := messageDao.QueryMessages(s.fromUserId, s.toUserId, s.lastTime, &s.messageList); err != nil {
		return err
	}

	return nil
}
