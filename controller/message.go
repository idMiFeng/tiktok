package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	UserId := strings.TrimSuffix(token, service.SALT)
	user_id, _ := strconv.ParseInt(UserId, 10, 64)
	toUserId := c.Query("to_user_id")
	to_user_id, _ := strconv.ParseInt(toUserId, 10, 64)
	// 格式化为 "yyyy-MM-dd HH:MM:ss" 格式的字符串
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	// 将字符串转换为time.Time类型
	parsedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", currentTime, time.Local)
	timestamp := parsedTime.Unix()
	content := c.Query("content")
	_, err := model.InsertMessage(user_id, to_user_id, content, timestamp)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "false",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "success",
		})
		return
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	UserId := strings.TrimSuffix(token, service.SALT)
	user_id, _ := strconv.ParseInt(UserId, 10, 64)
	toUserId := c.Query("to_user_id")
	to_user_id, _ := strconv.ParseInt(toUserId, 10, 64)
	pre_msg_time_string := c.Query("pre_msg_time")
	pre_msg_time, _ := strconv.ParseInt(pre_msg_time_string, 10, 64)
	var messages []model.Message
	messages, _ = model.GetMessagesByTime(user_id, to_user_id, pre_msg_time)
	c.JSON(http.StatusOK, gin.H{
		"status_code":  0,
		"status_msg":   "",
		"message_list": messages,
	})
}
