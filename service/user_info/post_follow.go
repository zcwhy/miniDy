package user_info

import (
	"errors"
	"miniDy/model"
)

type PostFollowFlow struct {
	userId     int64
	followId   int64
	actionType int
}

const (
	FOLLOW = 1
	CANCEL = 2
)

var (
	ErrIvdAct    = errors.New("未定义操作")
	ErrIvdFolUsr = errors.New("关注用户不存在")
	ErrIvdDel    = errors.New("还未关注无法取消")
)

func NewPostFollowFlow(userId int64, followId int64, actionType int) *PostFollowFlow {
	return &PostFollowFlow{userId: userId, followId: followId, actionType: actionType}
}

func (p *PostFollowFlow) checkNum() error {
	isUserExist := model.NewUserInfoDao().IsUserExistById(p.followId)
	if !isUserExist {
		return ErrIvdFolUsr
	}

	if p.actionType != FOLLOW && p.actionType != CANCEL {
		//log.Printf("*********++**%#v,********%#v", p.userId, p.followId)
		return ErrIvdAct
	}

	if p.userId == p.followId {
		return ErrIvdAct
	}

	isFollowExist := model.NewUserInfoDao().IsFollowExist(p.userId, p.followId)
	if p.actionType == CANCEL && !isFollowExist {
		return ErrIvdDel
	}

	return nil
}

func (p *PostFollowFlow) publish() error {
	var err error
	switch p.actionType {
	case FOLLOW:
		err = model.NewUserInfoDao().AddUserFollow(p.userId, p.followId)
	case CANCEL:
		err = model.NewUserInfoDao().CancelUserFollow(p.userId, p.followId)
	default:
		return ErrIvdAct
	}
	return err
}

func (p *PostFollowFlow) Do() error {
	if err := p.checkNum(); err != nil {
		return err
	}
	if err := p.publish(); err != nil {
		return err
	}
	return nil
}

func PostFollow(userId int64, followId int64, actionType int) error {
	return NewPostFollowFlow(userId, followId, actionType).Do()
}
