package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/idMiFeng/tiktok/dao"
	"github.com/idMiFeng/tiktok/model"
	"github.com/idMiFeng/tiktok/service"
	"net/http"
	"strconv"
	"strings"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	id := strings.TrimSuffix(token, service.SALT)
	user_id, _ := strconv.ParseInt(id, 10, 64)
	to_user_ID := c.Query("to_user_id")
	to_user_id, _ := strconv.ParseInt(to_user_ID, 10, 64)
	action_type := c.Query("action_type")
	switch action_type {
	case "1":
		//更新用户关系表
		Follow := model.Follow{
			UserID:    user_id,
			To_userId: to_user_id,
		}
		dao.DB.Create(&Follow)
		//用户自己关注总数加一,对方用户粉丝数加一
		user, _ := model.GetUserById(user_id)
		user.FollowCount++
		dao.DB.Save(&user)
		to_user, _ := model.GetUserById(to_user_id)
		to_user.FollowerCount++
		to_user.IsFollow = true
		dao.DB.Save(&to_user)
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
		})
		return
	case "2":
		Follow := model.Follow{
			UserID:    user_id,
			To_userId: to_user_id,
		}
		dao.DB.Where(&Follow).Delete(&Follow)
		user, _ := model.GetUserById(user_id)
		user.FollowCount--
		dao.DB.Save(&user)
		to_user, _ := model.GetUserById(to_user_id)
		to_user.FollowerCount--
		to_user.IsFollow = false
		dao.DB.Save(&to_user)
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
		})
		return

	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId := c.Query("user_id")
	user_id, _ := strconv.ParseInt(userId, 10, 64)
	//GetTuidsByUid函数得到关注的人的id
	toUserIDs, _ := model.GetTuidsByUid(user_id)
	var users []model.User
	_ = dao.DB.Where("id IN (?)", toUserIDs).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"user_list":   users,
	})
	return
}

// FollowerList all users have same follower list
// 通过查Follow表实现
func FollowerList(c *gin.Context) {
	userId := c.Query("user_id")
	user_id, _ := strconv.ParseInt(userId, 10, 64)
	//用GetUidsByTuid函数得到粉丝的id
	toUserIDs, _ := model.GetUidsByTuid(user_id)
	var users []model.User
	_ = dao.DB.Where("id IN (?)", toUserIDs).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"user_list":   users,
	})
	return
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	/* id := c.Query("user_id")
	   user_id, _ := strconv.ParseInt(id, 10, 64)

	   // 查询该用户关注的用户列表
	   var followingUsers []model.Follow
	   if err := dao.DB.Where("user_id = ?", user_id).Find(&followingUsers).Error; err != nil {
	       c.JSON(http.StatusOK, gin.H{
	           "status_code": 1,
	           "status_msg":  err.Error(),
	           "user_list":   nil,
	       })
	       return
	   }

	   // 查询关注该用户的用户列表
	   var followers []model.Follow
	   if err := dao.DB.Where("to_user_id = ?", user_id).Find(&followers).Error; err != nil {
	       c.JSON(http.StatusOK, gin.H{
	           "status_code": 1,
	           "status_msg":  err.Error(),
	           "user_list":   nil,
	       })
	       return
	   }

	   // 找出互相关注的用户
	   var friends []model.User
	   for _, followingUser := range followingUsers {
	       for _, follower := range followers {
	           if followingUser.To_userId == follower.UserID {
	               // 双方互相关注，将该用户加入好友列表
	               friend, err := model.GetUserById(followingUser.To_userId)
	               if err != nil {
	                   c.JSON(http.StatusOK, gin.H{
	                       "status_code": 1,
	                       "status_msg":  err.Error(),
	                       "user_list":   nil,
	                   })
	                   return
	               }
	               friends = append(friends, friend)
	               break
	           }
	       }
	   }

	   c.JSON(http.StatusOK, gin.H{
	       "status_code": 0,
	       "status_msg":  "",
	       "user_list":   friends,
	   })
	*/
	//为了测试方便设置为返回本用户
	id := c.Query("user_id")
	user_id, _ := strconv.ParseInt(id, 10, 64)
	user, _ := model.GetUserById(user_id)
	var users []model.User
	users = append(users, user)
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"user_list":   users,
	})
}
