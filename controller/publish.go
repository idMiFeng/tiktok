package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	username := strings.TrimSuffix(token, service.SALT)
	user, _ := model.GetUserByName(username)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")
	UserID, _ := strconv.ParseInt(user_id, 10, 64)
	videos, _ := model.GetVideoByUserId(UserID)
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"video_list":  videos,
	})
}
