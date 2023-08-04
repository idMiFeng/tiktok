package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"net/http"
	"time"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latest_time := c.Query("latest_time")
	Videos, _ := model.GetVideoByTime(latest_time)
	var nextTime int64
	if len(Videos) > 0 {
		nextTime = Videos[len(Videos)-1].CreatedTime.Unix()
	} else {
		// 如果视频列表为空，则设置下次请求的 NextTime 为当前时间
		nextTime = time.Now().Unix()
	}
	serverURL := "http://192.168.200.108:8080/"
	for i := range Videos {
		Videos[i].PlayUrl = serverURL + Videos[i].PlayUrl
		Videos[i].CoverUrl = serverURL + Videos[i].CoverUrl
		anthor, _ := model.GetUserById(Videos[i].UserID)
		Videos[i].Author = anthor
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"next_time":   nextTime,
		"video_list":  Videos,
	})
}
