package controllers

import (
	"golang_task4_blog_system/database"
	"golang_task4_blog_system/middleware"
	"golang_task4_blog_system/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 创建评论
func CreateComment(c *gin.Context) {
	var req models.Comment

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
			"error": "需要认证才能发表评论",
		})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, req.PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查文章失败",
		})
		return
	}

	// 创建评论
	comment := models.Comment{
		Content: req.Content,
		UserID:  currentUser.ID,
		PostID:  req.PostID,
	}

	// 保存到数据库
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建评论失败",
		})
		return
	}

	// 返回创建的评论信息（包含用户信息）
	var createdComment models.Comment
	if err := database.DB.Preload("User").First(&createdComment, comment.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取评论详情失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": createdComment,
	})
}

// GetPostComments 获取某篇文章的所有评论列表
func GetPostComments(c *gin.Context) {
	postID := c.Param("postId")

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查文章失败",
		})
		return
	}

	var comments []models.Comment
	var total int64

	// 获取评论总数
	database.DB.Model(&models.Comment{}).Where("post_id = ?", postID).Count(&total)

	// 获取评论列表（包含用户信息）
	if err := database.DB.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at ASC"). // 按创建时间正序排列
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取评论列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"post_id":  postID,
		"pagination": gin.H{
			"total": total,
		},
	})
}

// GetComment 获取单个评论详情
func GetComment(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if err := database.DB.Preload("User").Preload("Post").
		First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "评论不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

func UpdateComment(c *gin.Context) {
	var req models.Comment
	commentID := c.Param("id")

	// 验证输入
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "输入验证失败",
			"message": err.Error(),
		})
		return
	}

	// 查找评论
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "评论不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查找评论失败",
		})
		return
	}

	// 检查权限：只有评论作者可以更新
	currentUser := middleware.GetCurrentUser(c)
	if comment.UserID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "无权更新此评论",
		})
		return
	}

	// 检查文章是否存在（如果修改了文章ID）
	if req.PostID != comment.PostID {
		var post models.Post
		if err := database.DB.First(&post, req.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "目标文章不存在",
			})
			return
		}
	}

	// 更新评论
	comment.Content = req.Content
	comment.PostID = req.PostID

	if err := database.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新评论失败",
		})
		return
	}

	// 重新获取更新后的评论
	database.DB.Preload("User").First(&comment, commentID)

	c.JSON(http.StatusOK, gin.H{
		"message": "评论更新成功",
		"comment": comment,
	})
}


func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	// 查找评论
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "评论不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查找评论失败",
		})
		return
	}

	// 检查权限：只有评论作者可以删除
	currentUser := middleware.GetCurrentUser(c)
	if comment.UserID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "无权删除此评论",
		})
		return
	}

	// 删除评论
	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论删除成功",
	})
}

// GetMyComments 获取当前用户的所有评论
func GetMyComments(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "需要认证",
		})
		return
	}

	var comments []models.Comment
	var total int64

	// 获取评论总数
	database.DB.Model(&models.Comment{}).Where("user_id = ?", currentUser.ID).Count(&total)

	// 获取评论列表（包含用户和文章信息）
	if err := database.DB.Preload("User").Preload("Post").
		Where("user_id = ?", currentUser.ID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"total": total,
		},
	})
}
