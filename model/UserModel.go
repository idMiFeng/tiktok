package model

import "github.com/idMiFeng/tiktok/dao"

type UserRegister struct {
	Id       int64  `gorm:"primary_key;auto_increment"`
	Username string `gorm:"column:name;index"` //用Username字段查询次数多添加索引
	Password string `gorm:"column:password"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	UserID        int64  `gorm:"foreignkey:UserID"` // 外键字段，关联到 UserRegister 表的主键 Id
}

// InsertUser 插入用户
func InsertUser(username string, password string) (UserRegister, error) {
	user := UserRegister{
		Username: username,
		Password: password,
	}
	result := dao.DB.Create(&user)
	return user, result.Error
	// return true, nil
}

// GetUserByName 根据用户名查询用户
func GetUserByName(username string) (UserRegister, error) {
	user := UserRegister{}
	if err := dao.DB.Where("username = ?", username).Limit(1).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetUserById 根据Id查询用户
func GetUserById(Id int64) (UserRegister, error) {
	user := UserRegister{}
	if err := dao.DB.Where("Id = ?", Id).Limit(1).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
