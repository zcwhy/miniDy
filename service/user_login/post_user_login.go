package user_login

import (
	"errors"
	"miniDy/middleware"
	"miniDy/model"
	"miniDy/util"
)

type UserLoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type PostUserLoginFlow struct {
	username string
	password string

	UserId int64
	Token  string
}

func PostUserLogin(username, password string) (*UserLoginResponse, error) {
	return newPostUserLogin(username, password).do()
}

func newPostUserLogin(username, password string) *PostUserLoginFlow {
	return &PostUserLoginFlow{username: username, password: password}
}

func (u *PostUserLoginFlow) do() (*UserLoginResponse, error) {
	if err := u.checkParam(); err != nil {
		return nil, err
	}

	if err := u.prepareData(); err != nil {
		return nil, err
	}

	return &UserLoginResponse{UserId: u.UserId, Token: u.Token}, nil
}

func (u *PostUserLoginFlow) checkParam() error {
	if u.username == "" {
		return errors.New("用户名为空")
	}

	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *PostUserLoginFlow) prepareData() error {
	userLoginDao := model.NewLoginDao()
	userLogin, err := userLoginDao.QueryUserByName(u.username)

	//查询的错误处理
	if err != nil {
		return err
	}

	if userLogin.Username == "" {
		return errors.New("用户名不存在")
	}

	if userLogin.Password != util.SHA256(u.password) {
		return errors.New("用户名或密码错误")
	}

	//拿到id 拿到token
	u.UserId = userLogin.UserInfoId
	u.Token, err = middleware.ReleaseToken(userLogin)

	if err != nil {
		return err
	}

	return nil
}
