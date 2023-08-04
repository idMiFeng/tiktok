package model

import (
	"github.com/idMiFeng/tiktok/dao"
)

// 用户与关注用户关系维护
type Follow struct {
	Id        int64 `gorm:"primary_key;auto_increment"`
	UserID    int64 //用户ID
	To_userId int64 //此用户关注的人的ID
}

// 根据用户ID查询 Follow 表，并返回 To_userId 字段的值列表
func GetTuidsByUid(userID int64) ([]int64, error) {
	var follows []Follow
	result := dao.DB.Where("user_id = ?", userID).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}

	var toUserIDs []int64
	for _, follow := range follows {
		toUserIDs = append(toUserIDs, follow.To_userId)
	}

	return toUserIDs, nil
}

// 把用户ID作为被关注者字段来查follow表得到关注者(查一个用户的粉丝)
func GetUidsByTuid(userID int64) ([]int64, error) {
	var follows []Follow
	result := dao.DB.Where("to_user_id = ?", userID).Find(&follows)
	if result.Error != nil {
		return nil, result.Error
	}

	var userIDs []int64
	for _, follow := range follows {
		userIDs = append(userIDs, follow.UserID)
	}

	return userIDs, nil
}
