package model

import (
	"github.com/idMiFeng/tiktok/dao"
	"time"
)

type Video struct {
	Id            int64     `json:"id,omitempty"`
	Author        User      `json:"author"`
	PlayUrl       string    `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
	Tittle        string    `json:"tittle ,omitempty"`
	UserID        int64     `gorm:"foreignKey:UserID"`      // 外键字段，关联到 UserRegister 表的 Id
	CreatedTime   time.Time `json:"created_time,omitempty"` // 添加一个时间字段用于存储投稿时间戳
}

// 根据用户Id获取视频列表
func GetVideoByUserId(userId int64) ([]Video, error) {
	var videos []Video
	err := dao.DB.Where("user_id = ?", userId).Order("created_time desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoByTime(latest_time string) ([]Video, error) {
	latestTimeObj, _ := time.Parse(time.RFC3339, latest_time)
	var videos []Video
	err := dao.DB.Order("created_time desc").Where("created_time <= ?", latestTimeObj).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
