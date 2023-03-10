package message

import (
	"errors"
	"miniDy/model"
	"time"
)

type PostMessageActionService struct {
	fromUserId int64
	toUserId   int64
	actionType int32
	content    string
}

func PostMessage(fromUserId, toUserId int64, actionType int32, content string) error {
	return (&PostMessageActionService{
		fromUserId: fromUserId,
		toUserId:   toUserId,
		actionType: actionType,
		content:    content,
	}).Do()
}

func (p *PostMessageActionService) Do() error {
	if err := p.checkParam(); err != nil {
		return err
	}

	switch p.actionType {
	case 1:
		if err := p.updateData(); err != nil {
			return err
		}
	}

	return nil
}

func (p *PostMessageActionService) checkParam() error {
	userInfoDao := model.NewUserInfoDao()
	if userInfoDao.IsUserExistById(p.toUserId) == false {
		return errors.New("发送用户id错误")
	}

	if p.content == "" {
		return errors.New("发送内容为空")
	}

	return nil
}

func (p *PostMessageActionService) updateData() error {
	messageDao := model.NewMessageDao()
	message := &model.Message{
		FromUserId: p.fromUserId,
		ToUserId:   p.toUserId,
		Content:    p.content,
		CreateTime: time.Now().Unix(),
	}

	return messageDao.CreateMessage(message)
}
