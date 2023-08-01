package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"log"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	latest_time := c.Query("latest_time")
	log.Println(latest_time)
	Videos, _ := model.GetVideoByTime(latest_time)
	var nextTime int64
	if len(Videos) > 0 {
		nextTime = Videos[len(Videos)-1].CreatedTime.Unix()
	} else {
		// 如果视频列表为空，则设置下次请求的 NextTime 为当前时间
		nextTime = time.Now().Unix()
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"video_list":  Videos,
		"next_time":   nextTime,
	})
}
