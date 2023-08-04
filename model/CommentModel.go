package model

import (
	"github.com/idMiFeng/tiktok/dao"
	"log"
	"time"
)

type Comment struct {
	ID         int64  `json:"id" gorm:"primary_key;auto_increment"` // 评论id
	UserId     int64  //这条评论的用户的ID
	Content    string `json:"content"`     // 评论内容
	User       User   `json:"user"`        // 评论用户的具体信息
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
	VideoId    int64  `gorm:"column:video_id;index"`
	Video      Video  `gorm:"ForeignKey:VideoId"`
}

// 插入评论
func InsertComment(userId int64, videoId int64, content string) (Comment, error) {
	user, _ := GetUserById(userId)
	log.Println(user)
	currentTime := time.Now()
	formattedTime := currentTime.Format("01-02")
	comment := Comment{
		UserId:     userId,
		Content:    content,
		User:       user,
		CreateDate: formattedTime,
		VideoId:    videoId,
	}
	result := dao.DB.Create(&comment)
	if result.Error != nil {
		return comment, result.Error
	}
	return comment, nil
}

// 删除评论
func CancelComment(comment_id int64) error {
	comment := Comment{ID: comment_id}
	result := dao.DB.Delete(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 根据video_id返回comment列表
func GetCmtsByVid(video_id int64) ([]Comment, error) {
	var comments []Comment
	result := dao.DB.Where("video_id = ?", video_id).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}
