package user_login

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/user_login"
	"net/http"
)

type UserLoginResponse struct {
	response.CommonResp
	*user_login.UserLoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	resp, err := user_login.PostUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		UserLoginResponse: resp,
	})
}
