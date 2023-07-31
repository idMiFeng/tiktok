package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	data, err := c.FormFile("data") //用于获取 POST 请求中上传的文件数据。
	title := c.PostForm("title")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	username := strings.TrimSuffix(token, service.SALT)
	user, _ := model.GetUserInfoByName(username)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 创建Video对象并保存到数据库
	video := model.Video{
		Author:        user,
		PlayUrl:       saveFile,   // 这里假设保存的是视频文件路径
		CoverUrl:      "",         // 视频封面的URL
		FavoriteCount: 0,          // 默认点赞数为0
		CommentCount:  0,          // 默认评论数为0
		IsFavorite:    false,      // 默认未点赞
		Title:         title,      // 设置视频标题
		UserID:        user.Id,    // 外键字段关联到 UserRegister 表的 Id
		CreatedTime:   time.Now(), // 设置投稿时间为当前时间
	}

	if err := dao.DB.Create(&video).Error; err != nil {
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
