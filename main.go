package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
)

func main() {

	// 初始化Redis客户端连接
	if err := dao.InitRedisClient(); err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	//连接mysql数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close()

	//建表
	dao.DB.AutoMigrate(&model.UserRegister{}, &model.User{}, &model.Video{}, &model.Comment{}, &model.Follow{}, &model.Message{})

	go service.RunMessageServer()

	r := gin.Default()

	//注册路由
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
