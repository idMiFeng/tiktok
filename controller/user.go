package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// *********处理注册的函数************
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 判断是否为空
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "账号或密码为空",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	// 校验参数长度
	if len(password) > 32 || len(username) > 32 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "长度应小于32位",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	password = service.Encryption(password)

	//调用service层注册服务
	Id, err := service.RegisterService(username, password)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败",
			"user_id":     0,
			"token":       "",
		})
		return
	} else {
		//注册token
		token := service.GetTokenName(Id)

		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "注册成功",
			"user_id":     Id,
			"token":       token,
		})
	}
}

// *********处理登录的函数************
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//判断是否为空
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "账号或密码为空",
			"user_id":     0,
			"token":       "",
		})
		return
	}
	password = service.Encryption(password)
	fmt.Println(username, password)
	Id, err := service.LoginService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "登录失败",
			"user_id":     0,
			"token":       "",
		})
	} else {
		// 颁发token
		token := service.GetTokenName(Id)
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "登录成功",
			"user_id":     Id,
			"token":       token,
		})
	}

}

// *********查询用户信息的函数************
func UserInfo(c *gin.Context) {
	Id := c.Query("user_id")
	UserId, _ := strconv.ParseInt(Id, 10, 64)
	user, err := service.UserService(UserId)
	var total_favorited int64
	var work_count int64
	//查询点赞数量和作品数量并更新
	_ = dao.DB.Model(&Video{}).Where("is_favorite = ? AND user_id = ?", true, UserId).Count(&total_favorited).Error
	_ = dao.DB.Model(&Video{}).Where("user_id = ?", UserId).Count(&work_count).Error
	user.Total_favorited = total_favorited
	user.Work_count = work_count
	_ = dao.DB.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Id不匹配",
			"user":        nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "查询成功",
			"user":        user,
		})
	}
}
