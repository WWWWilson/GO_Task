package controllers

import (
	"golang_task4_blog_system/database"
	"golang_task4_blog_system/middleware"
	"golang_task4_blog_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建文章
func CreatePost(c *gin.Context) {
	var req models.Post

	// 验证输入数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "输入验证失败",
			"message": err.Error(),
		})
		return
	}

	// 获取当前用户（通过BasicAuth中间件认证）
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "需要认证才能创建文章",
		})
		return
	}

	// 创建文章
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  currentUser.ID,
	}

	// 保存到数据库
	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建文章失败",
		})
		return
	}

	// 返回创建的文章信息（包含用户信息）
	var createdPost models.Post
	if err := database.DB.Preload("User").First(&createdPost, post.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文章详情失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    createdPost,
	})
}

// 获取所有文章
func GetPosts(c *gin.Context) {
	var posts []models.Post
	var total int64

	// 获取文章总数
	database.DB.Model(&models.Post{}).Count(&total)

	// 获取文章列表（包含用户信息）
	if err := database.DB.Preload("User").
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文章列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"total": total,
		},
	})
}

// 获取单个文章
func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	// 获取文章详情，包含用户信息和评论
	if err := database.DB.Preload("User").Preload("Comments.User").
		First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// 更新文章
func UpdatePost(c *gin.Context) {
	var req models.Post
	PostID := c.Param("id")

	// 验证输入
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "输入验证失败",
			"message": err.Error(),
		})
		return
	}

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查找文章失败",
		})
		return
	}

	// 检查权限：只有文章作者可以更新
	postIDUint, _ := strconv.ParseUint(PostID, 10, 32)
	if !middleware.IsPostAuthor(c, uint(postIDUint)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "无权更新此文章",
		})
		return
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	// 更新文章
	if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新文章失败",
		})
		return
	}

	// 重新获取更新后的文章
	database.DB.Preload("User").First(&post, PostID)

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// 删除文章
func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查找文章失败",
		})
		return
	}

	// 检查权限：只有文章作者可以删除
	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	if !middleware.IsPostAuthor(c, uint(postIDUint)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "无权删除此文章",
		})
		return
	}

	// 删除文章（GORM 会自动处理关联的评论，因为有外键约束）
	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
	})
}

// GetMyPosts 获取当前用户的文章
func GetMyPosts(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "需要认证",
		})
		return
	}

	var posts []models.Post
	if err := database.DB.Preload("User").
		Where("user_id = ?", currentUser.ID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}


