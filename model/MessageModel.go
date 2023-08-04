package model

import (
	"github.com/idMiFeng/tiktok/dao"
)

type Message struct {
	ID           int64  `json:"id" gorm:"primary_key;auto_increment"`
	From_user_id int64  `json:"from_user_id"`
	ToUserId     int64  `json:"to_user_id"`
	Content      string `json:"content"`
	CreateTime   int64  `json:"create_time" gorm:"column:create_time"` // 创建时间
}

func InsertMessage(UserId int64, ToUserId int64, Content string, CreateTime int64) (Message, error) {
	Message := Message{
		From_user_id: UserId,
		ToUserId:     ToUserId,
		Content:      Content,
		CreateTime:   CreateTime,
	}
	result := dao.DB.Create(&Message)
	if result.Error != nil {
		return Message, result.Error
	}

	return Message, nil
}

func GetMessagesByTime(UserId int64, ToUserId int64, Pre_msg_time int64) ([]Message, error) {
	var messages []Message
	err := dao.DB.Where("from_user_id = ? AND to_user_id = ? AND create_time>?", UserId, ToUserId, Pre_msg_time).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
