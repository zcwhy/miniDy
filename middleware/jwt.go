package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"miniDy/model"
	"miniDy/model/response"
	"net/http"
	"time"
)

var jwtKey = []byte("zcwhy")

type Claims struct {
	UserId int64
	jwt.RegisteredClaims
}

func ReleaseToken(user *model.UserLogin) (string, error) {
	expirationTime := time.Now().Add(15 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.UserInfoId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			Subject:   "user token",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ParseToken(tokenString string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if token != nil {
		//类型断言
		if claims, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return claims, true
			} else {
				return claims, false
			}
		}
	}
	return nil, false
}

func JWTMiddleWare(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		token = c.PostForm("token")
	}

	//用户不存在
	if token == "" {
		c.JSON(http.StatusOK, response.CommonResp{StatusCode: 401, StatusMsg: "用户不存在"})
		c.Abort() //阻止执行
		return
	}

	//验证token
	claims, ok := ParseToken(token)
	if !ok {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: 403,
			StatusMsg:  "token不正确",
		})
		c.Abort() //阻止执行
		return
	}

	if time.Now().Unix() > claims.ExpiresAt.Time.Unix() {
		c.JSON(http.StatusOK, response.CommonResp{
			StatusCode: 402,
			StatusMsg:  "token过期",
		})
		c.Abort() //阻止执行
		return
	}

	c.Set("user_id", claims.UserId)
	c.Next()
}
