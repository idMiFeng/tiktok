package model

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/idMiFeng/tiktok/dao"
	"log"
	"strconv"
)

type UserRegister struct {
	Id       int64  `gorm:"primary_key;auto_increment"`
	Username string `gorm:"column:name;index"` //用Username字段查询次数多添加索引
	Password string `gorm:"column:password"`
}

type User struct {
	Id               int64        `json:"id,omitempty" gorm:"primary_key;auto_increment"`
	Name             string       `json:"name,omitempty" gorm:"column:name;index" gorm:"column:name;index"`
	FollowCount      int64        `json:"follow_count,omitempty"`
	FollowerCount    int64        `json:"follower_count,omitempty"`
	Avatar           string       `json:"avatar,omitempty"`
	Background_image string       `json:"background_image,omitempty"`
	IsFollow         bool         `json:"is_follow,omitempty"`
	Work_count       int64        `json:"work_count,omitempty"`
	Total_favorited  int64        `json:"total_favorited,omitempty"`
	UserID           int64        `gorm:"column:user_id;index"` // 外键字段，关联到 UserRegister 表的主键 Id
	UserRegister     UserRegister `gorm:"ForeignKey:UserID"`
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
		Avatar:        "https://profile-avatar.csdnimg.cn/b74263567d8a4bf78e3c0b649a622a54_weixin_63603830.jpg!1",
		UserID:        userRegister.Id, // Use the generated Id as the UserID
	}
	result = dao.DB.Create(&userInfo)
	if result.Error != nil {
		return userRegister, result.Error
	}
	//将用户信息缓存到redis数据库
	err := CacheUser(dao.Rdb, userInfo)
	if err != nil {
		log.Println("用户信息缓存失败", err)
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

// 将用户信息缓存到redis
func CacheUser(redisClient *redis.Client, user User) error {
	ctx := context.Background()

	// 将用户信息序列化成JSON字符串
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 设置key为UserID，value为用户信息的JSON字符串
	key := strconv.FormatInt(user.UserID, 10)
	err = redisClient.Set(ctx, key, userData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// 从redis缓存中获得用户信息
func GetUserByRedis(redisClient *redis.Client, userID int64) (User, error) {
	ctx := context.Background()
	key := strconv.FormatInt(userID, 10)
	userData, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Redis缓存中不存在，从数据库中获取用户信息
		user, err := GetUserById(userID)
		if err != nil {
			return User{}, err
		}

		// 将用户信息缓存到Redis
		err = CacheUser(redisClient, user)
		if err != nil {
			return user, err
		}

		return user, nil
	} else if err != nil {
		return User{}, err
	}

	// 从Redis缓存中反序列化用户信息
	var user User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
