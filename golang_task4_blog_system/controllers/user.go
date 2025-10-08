package controllers

import (
	"golang_task4_blog_system/database"
	"golang_task4_blog_system/models"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		c.JSON(500, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建用户
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(200, gin.H{"message": "注册成功"})
}

// 登录（返回成功信息，实际认证通过 BasicAuth）
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 验证用户
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "用户不存在"})
		return
	}

	if !user.CheckPassword(input.Password) {
		c.JSON(401, gin.H{"error": "密码错误"})
		return
	}

	c.JSON(200, gin.H{"message": "登录成功，使用 Basic Auth 访问受保护接口"})
}
