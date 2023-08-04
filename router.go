package main

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/controller"
	"github.com/idMiFeng/tiktok/middle"
	"net/http"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources设置浏览器能通过url访问本地文件
	r.Static("/public", "./public")
	r.LoadHTMLGlob("templates/*")

	// home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middle.VerifyToken, controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middle.VerifyTokenQuery, controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", middle.VerifyTokenQuery, controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
