package middle

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"strings"
)

func VerifyToken(c *gin.Context) {
	token := c.PostForm("token")
	if len(token) == 0 {
		//请求处理中止
		c.Abort()
		//返回json
		c.JSON(http.StatusOK,
			gin.H{
				"status_code": 1,
				"status_msg":  "token为空",
			})
		return
	}
	flag, err := ParseToken(token)
	if flag != true {
		c.JSON(http.StatusOK,
			gin.H{
				"status_code": 1,
				"status_msg":  err,
			})
	} else {
		//解析正确
		c.Next()
	}
}

func VerifyTokenQuery(c *gin.Context) {
	token := c.Query("token")
	if len(token) == 0 {
		//请求处理中止
		c.Abort()
		//返回json
		c.JSON(http.StatusOK,
			gin.H{
				"status_code": 1,
				"status_msg":  "token为空",
			})
		return
	}
	flag, err := ParseToken(token)
	if flag != true {
		c.JSON(http.StatusOK,
			gin.H{
				"status_code": 1,
				"status_msg":  err,
			})
	} else {
		//解析正确
		c.Next()
	}
}

// 解析token，判断用户是否登录
func ParseToken(token string) (bool, error) {
	username := strings.TrimSuffix(token, service.SALT)
	//根据token解析出来的用户名查表，找到则说明用户已登录
	_, err := model.GetUserByName(username)
	if err != nil {
		return false, errors.New("用户不存在")
	}
	return true, nil

}
