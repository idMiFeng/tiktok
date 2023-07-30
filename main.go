package main

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
)

func main() {
	//连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close()

	//建表
	dao.DB.AutoMigrate(&model.UserRegister{})

	go service.RunMessageServer()

	r := gin.Default()

	//注册路由
	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
