package model

import (
	"fmt"
	"testing"
)

func TestMessageDAO_QueryMessages(t *testing.T) {
	messageDao := NewMessageDao()
	messageList := make([]*Message, 0)

	err := messageDao.QueryMessages(8, 10, 1677065980, &messageList)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(messageList) != 0 {
		t.Errorf("error")
	}
	fmt.Println(len(messageList))
}
