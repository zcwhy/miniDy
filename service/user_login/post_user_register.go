package user_login

import (
	"errors"
	"miniDy/constant"
	"miniDy/middleware"
	"miniDy/model"
	"miniDy/util"
)

/*
	从handler处拿到收集的参数
	1、构造流对象， 执行流对象的do方法
	2、参数校验
	3、收集dao层查询的数据，封装成返回对象返回给handler层
*/

type UserRegisterResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type PostUserRegisterFlow struct {
	username string
	password string

	userid int64
	token  string
}

func PostUserRegister(username string, password string) (*UserRegisterResponse, error) {
	return NewPostUserRegister(username, password).Do()
}

// NewPostUserRegister 构造一个流对象
func NewPostUserRegister(username, password string) *PostUserRegisterFlow {
	return &PostUserRegisterFlow{username: username, password: password}
}

func (q *PostUserRegisterFlow) Do() (*UserRegisterResponse, error) {
	if err := q.checkParam(); err != nil {
		return nil, err
	}
	if err := q.updateData(); err != nil {
		return nil, err
	}

	return &UserRegisterResponse{Token: q.token, UserId: q.userid}, nil
}

func (q *PostUserRegisterFlow) checkParam() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}

	if q.password == "" {
		return errors.New("密码为空")
	}

	if (len(q.password) > constant.MAX_PASSWORD_LEN) ||
		(len(q.username) > constant.MAX_USERNAME_LEN) {
		return errors.New("用户名或密码长度超出限制")
	}

	return nil
}

func (q *PostUserRegisterFlow) updateData() error {
	userLoginDao := model.NewLoginDao()
	userExist := userLoginDao.IsUserExist(q.username)

	if userExist == true {
		return errors.New("用户名已存在")
	}

	//构造结构体时对密码进行加密，checkParam的时候不能加密，加密后的长度会超过32
	shaPassword := util.SHA256(q.password)
	userLogin := model.UserLogin{Username: q.username, Password: shaPassword}
	userinfo := model.UserInfo{User: &userLogin, Name: q.username}

	//更新操作，由于userLogin属于userInfo，故更新userInfo
	userInfoDao := model.NewUserInfoDao()
	err := userInfoDao.UserRegister(&userinfo)

	if err != nil {
		return err
	}

	token, err := middleware.ReleaseToken(&userLogin)
	if err != nil {
		return err
	}
	q.token = token
	q.userid = userinfo.Id

	return nil
}
