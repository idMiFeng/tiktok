package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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
	userID := strings.TrimSuffix(token, service.SALT)
	ID, _ := strconv.ParseInt(userID, 10, 64)
	user, _ := model.GetUserById(ID)
	user.Work_count++
	dao.DB.Save(&user)
	finalName := fmt.Sprintf("%d_%s", user.UserID, filename)
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
		PlayUrl:       saveFile,    // 这里保存的是视频文件路径
		CoverUrl:      "",          // 视频封面的URL
		FavoriteCount: 0,           // 默认点赞数为0
		CommentCount:  0,           // 默认评论数为0
		IsFavorite:    false,       // 默认未点赞
		Title:         title,       // 设置视频标题
		UserID:        user.UserID, // 外键字段关联到 UserRegister 表的 Id
		CreatedTime:   time.Now(),  // 设置投稿时间为当前时间
	}
	//截取视频第一帧作为封面
	if err := generateVideoCover(saveFile, &video); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
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

// 使用Golang中的 exec 包来执行 FFmpeg 命令对视频进行截图
func generateVideoCover(videoPath string, video *model.Video) error {
	// 设置FFmpeg命令及参数
	ffmpegCmd := "ffmpeg"
	//用当前时间戳作为随机字符串用于生成截图保存地址
	currentTime := time.Now()
	timestamp := currentTime.Unix()
	outputFile := "./public/" + strconv.FormatInt(timestamp, 10) + "cover.jpg"

	// 创建命令对象
	cmd := exec.Command(ffmpegCmd, "-i", videoPath, "-ss", "00:00:01", "-vframes", "1", outputFile)

	// 执行命令并等待执行完成
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行FFmpeg命令失败: %s", err)
	}

	// 检查输出文件是否存在
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return fmt.Errorf("截图文件不存在")
	}
	// 将截图的URL赋值给CoverUrl字段
	coverURL := filepath.Join(outputFile)
	video.CoverUrl = coverURL
	return nil
}
