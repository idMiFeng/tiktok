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
		//用户获赞数加一
		user, _ := model.GetUserById(video.UserID)
		user.Total_favorited++
		dao.DB.Save(&user)
		dao.DB.Model(&model.Video{}).Where("id = ?", id).Update(video)

	} else {
		video.IsFavorite = false
		if video.FavoriteCount > 0 {
			video.FavoriteCount--
			//用户获赞减一
			user, _ := model.GetUserById(video.UserID)
			user.Total_favorited--
			dao.DB.Save(&user)
		}
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
	//id := c.Query("user_id")
	//user_id, _ := strconv.ParseInt(id, 10, 64)
	var videos []Video
	dao.DB.Where("is_favorite = ?", true).Find(&videos)
	serverURL := "http://192.168.200.108:8080/"
	for i := range videos {
		videos[i].CoverUrl = serverURL + videos[i].CoverUrl
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"video_list":  videos,
	})
}
