package model

import "github.com/idMiFeng/tiktok/dao"

type UserRegister struct {
	Id       int64  `gorm:"primary_key;auto_increment"`
	Username string `gorm:"column:name;index"` //用Username字段查询次数多添加索引
	Password string `gorm:"column:password"`
}

type User struct {
	Id            int64        `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Name          string       `json:"name,omitempty" gorm:"column:name;index" gorm:"column:name;index"`
	FollowCount   int64        `json:"follow_count,omitempty"`
	FollowerCount int64        `json:"follower_count,omitempty"`
	IsFollow      bool         `json:"is_follow,omitempty"`
	UserID        int64        // 外键字段，关联到 UserRegister 表的主键 Id
	UserRegister  UserRegister `gorm:"ForeignKey:UserID"`
}

// InsertUser 插入用户,同时往User表和UserRegister表写入数据
func InsertUser(username string, password string) (UserRegister, error) {
	userRegister := UserRegister{
		Username: username,
		Password: password,
	}
	result := dao.DB.Create(&userRegister)
	if result.Error != nil {
		return userRegister, result.Error
	}

	// Use the generated Id to insert data into User table
	userInfo := User{
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		UserID:        userRegister.Id, // Use the generated Id as the UserID
	}
	result = dao.DB.Create(&userInfo)
	if result.Error != nil {
		return userRegister, result.Error
	}

	return userRegister, nil
}

// GetUserByName 根据用户名查询用户
func GetUserByName(username string) (UserRegister, error) {
	user := UserRegister{}
	if err := dao.DB.Where("name = ?", username).Limit(1).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetUserInfoByName 根据用户名查找User表
func GetUserInfoByName(username string) (User, error) {
	userinfo := User{}
	if err := dao.DB.Where("name = ?", username).Limit(1).Find(&userinfo).Error; err != nil {
		return userinfo, err
	}
	return userinfo, nil
}

// GetUserById 根据Id查询查找User表
func GetUserById(Id int64) (User, error) {
	user := User{}
	if err := dao.DB.Where("user_id = ?", Id).Limit(1).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
