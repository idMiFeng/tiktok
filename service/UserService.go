package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/idMiFeng/tiktok/model"
	"log"
	"strconv"
)

// SALT 盐值
const SALT = "BYTEDANCE"

// GetTokenName 根据用户Id生成token
func GetTokenName(Id int64) string {
	token := strconv.FormatInt(Id, 10) + SALT
	return token
}

// Encryption md5加盐加密
func Encryption(password string) string {
	password += SALT
	hash := md5.New()
	hash.Write([]byte(password))
	hash_password := hex.EncodeToString(hash.Sum(nil))
	return hash_password
}

// RegisterService 注册服务
func RegisterService(username string, password string) (int64, error) {
	log.Println(username, "---", password)

	// 查表，是否存在id
	user, err := model.GetUserByName(username)
	if err != nil {
		return 0, errors.New("数据库查询错误")
	}
	if username == user.Username {
		return 0, errors.New("用户名已经存在")
	}
	// 插入
	user, _ = model.InsertUser(username, password)
	return user.Id, err
}

// LoginService 登录服务
func LoginService(username string, password string) (int64, error) {

	user, err := model.GetUserByName(username)
	if err != nil {
		return 0, errors.New("用户名不存在")
	}

	if password != user.Password {
		return 0, errors.New("密码不正确")
	}

	return user.Id, nil
}

// UserService 用户查询服务
func UserService(Id int64) (model.UserRegister, error) {
	user, err := model.GetUserById(Id)
	if err != nil {
		return model.UserRegister{}, errors.New("用户不存在")
	}
	return user, nil
}
