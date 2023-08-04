package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	actionType := c.Query("action_type")
	switch actionType {
	case "1": //发布评论
		//根据token获得用户ID
		token := c.Query("token")
		id := strings.TrimSuffix(token, service.SALT)
		user_id, _ := strconv.ParseInt(id, 10, 64)
		comment_text := c.Query("comment_text")
		videoId := c.Query("video_id")
		video_id, _ := strconv.ParseInt(videoId, 10, 64)
		comment, _ := model.InsertComment(user_id, video_id, comment_text)
		log.Println(comment)
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
			"comment":     comment,
		})
		return
	case "2": //删除评论
		commentId := c.Query("comment_id")
		comment_id, _ := strconv.ParseInt(commentId, 10, 64)
		_ = model.CancelComment(comment_id)
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
		})
		return
	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")
	video_id, _ := strconv.ParseInt(videoId, 10, 64)
	comments, _ := model.GetCmtsByVid(video_id)
	for i := range comments {
		user, _ := model.GetUserById(comments[i].UserId)
		comments[i].User = user
	}
	c.JSON(http.StatusOK, gin.H{"status_code": 0,
		"status_msg":   "",
		"comment_list": comments,
	})
}
