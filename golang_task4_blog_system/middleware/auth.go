package middleware

import (
	"golang_task4_blog_system/database"
	"golang_task4_blog_system/models"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		// 这里留空，我们会动态验证
	})
}

func BasicAuth() gin.HandlerFunc {
	// 从数据库加载所有用户到内存
	// accounts := gin.Accounts{}
	// var user models.User
	// database.DB.Find(&user)
	// accounts[user.Username] = user.Password
	return gin.BasicAuth(gin.Accounts{ //临时测试
		"user_name":"wilson",
		"Password":"wilson1234",
	})
}

// GetCurrentUser 获取当前登录用户
func GetCurrentUser(c *gin.Context) *models.User {
	username, exists := c.Get(gin.AuthUserKey)
	if !exists {
		return nil
	}
	
	var user models.User
	if err := database.DB.Where("username = ?", username.(string)).First(&user).Error; err != nil {
		return nil
	}
	
	return &user
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	user := GetCurrentUser(c)
	if user == nil {
		return 0
	}
	return user.ID
}

// IsPostAuthor 检查当前用户是否是文章作者
func IsPostAuthor(c *gin.Context, postID uint) bool {
	currentUserID := GetCurrentUserID(c)
	if currentUserID == 0 {
		return false
	}
	
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		return false
	}
	
	return post.UserID == currentUserID
}
