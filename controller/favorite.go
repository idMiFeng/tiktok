package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {

	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	id, _ := strconv.ParseInt(video_id, 10, 64)
	Type, _ := strconv.ParseInt(action_type, 10, 64)
	var video model.Video
	//查表获得该video
	if err := dao.DB.Where("id = ?", id).First(&video).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Video not found",
		})
		return
	}
	if Type == 1 {
		video.IsFavorite = true
		video.FavoriteCount++
	} else {
		video.IsFavorite = false
		video.FavoriteCount--
	}
	//修改内容后更新数据库
	if err := dao.DB.Save(&video).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Failed to update video",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
		})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
