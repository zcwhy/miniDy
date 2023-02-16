package model

import (
	"errors"
	"sync"
	"time"
)

type Message struct {
	Id         int64
	ToUserId   int64 `json:"to_user_id"`
	FromUserId int64 `json:"from_user_id"`
	Content    string
	CreateTime time.Time `json:"create_time"`
}

type MessageDAO struct {
}

var (
	messageDao  *MessageDAO
	messageOnce sync.Once
)

func NewMessageDao() *MessageDAO {
	messageOnce.Do(func() {
		messageDao = &MessageDAO{}
	})

	return messageDao
}

func (m *MessageDAO) CreateMessage(message *Message) error {
	if message == nil {
		return errors.New("CreateMessage message 空指针")
	}

	return DB.Create(message).Error
}
