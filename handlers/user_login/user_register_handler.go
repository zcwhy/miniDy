package user_login

import (
	"github.com/gin-gonic/gin"
	"miniDy/constant"
	"miniDy/model/response"
	"miniDy/service/user_login"
	"net/http"
)

type UserRegisterResponse struct {
	response.CommonResp
	*user_login.UserRegisterResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	registerResponse, err := user_login.PostUserRegister(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResp: response.CommonResp{
				StatusCode: constant.FAILURE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResp: response.CommonResp{
			StatusCode: constant.SUCCESS,
			StatusMsg:  constant.SUCCESS_MESSAGE,
		},
		UserRegisterResponse: registerResponse,
	})
}
