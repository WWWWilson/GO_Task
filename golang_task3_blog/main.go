package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	FirstName    string    `gorm:"size:50" json:"first_name"`
	LastName     string    `gorm:"size:50" json:"last_name"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	PostCount    int       `gorm:"default:0" json:"post_count"`      // 用户文章数量统计
	CommentCount int       `gorm:"default:0" json:"comment_count"`   // 用户评论数量统计
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// Post 文章模型
type Post struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Title        string    `gorm:"size:200;not null" json:"title"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	Summary      string    `gorm:"size:500" json:"summary"`
	Slug         string    `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Status       string    `gorm:"size:20;default:draft" json:"status"` // draft, published, archived
	CommentStatus string   `gorm:"size:20;default:no_comments" json:"comment_status"` // no_comments, has_comments
	ViewCount    uint      `gorm:"default:0" json:"view_count"`
	CommentCount int       `gorm:"default:0" json:"comment_count"` // 文章评论数量统计
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	UserID   uint      `gorm:"not null;index" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}


// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    string    `gorm:"size:20;default:pending" json:"status"` // pending, approved, rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 外键：关联文章
	PostID uint `gorm:"not null;index" json:"post_id"`
	// 多对一关系：评论属于文章
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	
	// 外键：关联用户（评论作者）
	UserID uint `gorm:"not null;index" json:"user_id"`
	// 多对一关系：评论属于用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	
	// 自引用关系：支持回复评论
	ParentID *uint    `gorm:"index" json:"parent_id"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

var db *gorm.DB

func initDB() {
	dsn := "username:password@tcp(localhost:3306)/blog_system?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
}

// 1. 查询某个用户发布的所有文章及其对应的评论信息
func getUserPostsWithComments(userID uint) ([]Post, error) {
	var posts []Post
	
	// 使用 Preload 预加载关联数据
	err := db.Where("user_id = ?", userID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, email, first_name, last_name")
		}).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved"). // 只加载已批准的评论
					Preload("User", func(db *gorm.DB) *gorm.DB {
						return db.Select("id, username, first_name, last_name")
					})
		}).
		Find(&posts).Error
	
	if err != nil {
		return nil, fmt.Errorf("查询用户文章失败: %v", err)
	}
	
	return posts, nil
}

// 2. 查询评论数量最多的文章信息
func getMostCommentedPost() (*Post, int, error) {
	var post Post
	var commentCount int
	
	// 方法1：使用子查询
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id AND comments.status = 'approved'").
		Group("posts.id").
		Order("comment_count DESC").
		First(&post).
		Scan(&commentCount).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("查询评论最多文章失败: %v", err)
	}
	
	return &post, commentCount, nil
}

// 2.1 另一种方式：查询评论数量最多的文章（带完整信息）
func getMostCommentedPostWithDetails() (map[string]interface{}, error) {
	var result map[string]interface{}
	
	err := db.Model(&Post{}).
		Select("posts.*, users.username, users.first_name, users.last_name, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN users ON users.id = posts.user_id").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id AND comments.status = 'approved'").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&result).Error
	
	if err != nil {
		return nil, fmt.Errorf("查询评论最多文章详情失败: %v", err)
	}
	
	return result, nil
}


// ==================== Post 模型钩子函数 ====================

// BeforeCreate 在创建文章之前自动更新用户的文章数量统计
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeCreate 钩子: 正在为文章 '%s' 更新用户文章统计\n", p.Title)
	
	// 更新用户的文章数量
	result := tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1))
	
	if result.Error != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在: %d", p.UserID)
	}
	
	fmt.Printf("用户 %d 的文章数量已更新\n", p.UserID)
	return nil
}

// AfterCreate 在创建文章之后设置默认的评论状态
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("AfterCreate 钩子: 文章 '%s' 创建完成，ID: %d\n", p.Title, p.ID)
	return nil
}

// BeforeDelete 在删除文章之前更新用户的文章数量统计
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeDelete 钩子: 正在删除文章 '%s'，更新用户统计\n", p.Title)
	
	// 减少用户的文章数量
	result := tx.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("GREATEST(0, post_count - 1)"))
	
	if result.Error != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", result.Error)
	}
	
	fmt.Printf("用户 %d 的文章数量已减少\n", p.UserID)
	return nil
}

// ==================== Comment 模型钩子函数 ====================

// BeforeCreate 在创建评论之前更新文章和用户的评论数量统计
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeCreate 钩子: 正在为文章 %d 创建评论\n", c.PostID)
	
	// 更新文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).
		Updates(map[string]interface{}{
			"comment_count": gorm.Expr("comment_count + 1"),
			"comment_status": "has_comments",
		})
	
	if result.Error != nil {
		return fmt.Errorf("更新文章评论数量失败: %v", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("文章不存在: %d", c.PostID)
	}
	
	// 更新用户的评论数量
	result = tx.Model(&User{}).Where("id = ?", c.UserID).
		Update("comment_count", gorm.Expr("comment_count + 1"))
	
	if result.Error != nil {
		return fmt.Errorf("更新用户评论数量失败: %v", result.Error)
	}
	
	fmt.Printf("文章 %d 和用户 %d 的评论数量已更新\n", c.PostID, c.UserID)
	return nil
}

// BeforeDelete 在删除评论之前保存相关信息用于后续处理
func (c *Comment) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("BeforeDelete 钩子: 正在删除评论 %d\n", c.ID)
	
	// 这里可以保存评论的相关信息，用于AfterDelete钩子
	// 但由于Gorm的限制，我们主要在AfterDelete中处理
	return nil
}

// AfterDelete 在删除评论之后检查文章的评论数量并更新状态
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("AfterDelete 钩子: 评论 %d 已删除，检查文章评论状态\n", c.ID)
	
	// 减少文章的评论数量
	result := tx.Model(&Post{}).Where("id = ?", c.PostID).
		Update("comment_count", gorm.Expr("GREATEST(0, comment_count - 1)"))
	
	if result.Error != nil {
		return fmt.Errorf("更新文章评论数量失败: %v", result.Error)
	}
	
	// 减少用户的评论数量
	result = tx.Model(&User{}).Where("id = ?", c.UserID).
		Update("comment_count", gorm.Expr("GREATEST(0, comment_count - 1)"))
	
	if result.Error != nil {
		return fmt.Errorf("更新用户评论数量失败: %v", result.Error)
	}
	
	// 检查文章的评论数量，如果为0则更新评论状态
	var post Post
	if err := tx.First(&post, "id = ?", c.PostID).Error; err != nil {
		return fmt.Errorf("查询文章失败: %v", err)
	}
	
	if post.CommentCount == 0 {
		result := tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", "no_comments")
		
		if result.Error != nil {
			return fmt.Errorf("更新文章评论状态失败: %v", result.Error)
		}
		
		fmt.Printf("文章 %d 的评论状态已更新为 'no_comments'\n", c.PostID)
	} else {
		fmt.Printf("文章 %d 还有 %d 条评论，评论状态保持不变\n", c.PostID, post.CommentCount)
	}
	
	return nil
}

// ==================== 业务函数 ====================

// CreatePost 创建文章（封装了业务逻辑）
func CreatePost(post *Post) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 这里会触发 BeforeCreate 钩子
		if err := tx.Create(post).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteComment 删除评论（封装了业务逻辑）
func DeleteComment(commentID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var comment Comment
		if err := tx.First(&comment, commentID).Error; err != nil {
			return fmt.Errorf("评论不存在: %v", err)
		}
		
		// 这里会触发 BeforeDelete 和 AfterDelete 钩子
		if err := tx.Delete(&comment).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// GetUserStats 获取用户统计信息
func GetUserStats(userID uint) (*User, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetPostStats 获取文章统计信息
func GetPostStats(postID uint) (*Post, error) {
	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}



func main(){
	// 数据库连接配置
	initDB()
	// 自动迁移创建表
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal("数据库迁移失败: ", err)
	}
	
	fmt.Println("=== 博客系统钩子函数演示 ===")

	// 示例用户ID
	userID := uint(1)
	
	// 创建测试用户
	user := User{
		Username:  "test_user",
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "Test",
		LastName:  "User",
		Status:    "active",
	}
	
	fmt.Println("=== 查询用户所有文章及评论 ===")
	
	// 查询用户的所有文章及其评论
	posts, err := getUserPostsWithComments(userID)
	if err != nil {
		log.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("用户 %d 共有 %d 篇文章\n", userID, len(posts))
		for i, post := range posts {
			fmt.Printf("文章 %d: 《%s》 (评论数: %d)\n", i+1, post.Title, len(post.Comments))
			for j, comment := range post.Comments {
				fmt.Printf("  - 评论 %d: %s (作者: %s)\n", j+1, comment.Content, comment.User.Username)
			}
		}
	}
	
	fmt.Println("\n=== 查询评论数量最多的文章 ===")
	mostCommentedPost, count, err := getMostCommentedPost()
	if err != nil {
		log.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("评论最多的文章: 《%s》 (评论数: %d)\n", mostCommentedPost.Title, count)
	}
	
	fmt.Println("\n=== 查询评论最多文章详情 ===")
	postDetail, err := getMostCommentedPostWithDetails()
	if err != nil {
		log.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("文章详情: 《%s》, 作者: %s, 评论数: %d\n", 
			postDetail["title"], postDetail["username"], postDetail["comment_count"])
	}

	//钩子
	// 检查用户是否已存在
	var existingUser User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&user).Error; err != nil {
				log.Fatal("创建用户失败: ", err)
			}
			fmt.Printf("创建用户成功: %s (ID: %d)\n", user.Username, user.ID)
		} else {
			log.Fatal("查询用户失败: ", err)
		}
	} else {
		user = existingUser
		fmt.Printf("使用现有用户: %s (ID: %d)\n", user.Username, user.ID)
	}
	
	// 演示1: 创建文章（触发Post的BeforeCreate钩子）
	fmt.Println("\n--- 演示1: 创建文章 ---")
	post := Post{
		Title:   "测试钩子函数的文章",
		Content: "这篇文章用于测试Gorm钩子函数",
		Summary: "测试文章摘要",
		Slug:    "test-hook-article",
		Status:  "published",
		UserID:  user.ID,
	}
	
	if err := CreatePost(&post); err != nil {
		log.Printf("创建文章失败: %v\n", err)
	} else {
		fmt.Printf("文章创建成功: 《%s》 (ID: %d)\n", post.Title, post.ID)
	}
	
	// 检查用户文章数量
	userStats, _ := GetUserStats(user.ID)
	fmt.Printf("用户当前文章数量: %d\n", userStats.PostCount)
	
	// 演示2: 创建评论（触发Comment的BeforeCreate钩子）
	fmt.Println("\n--- 演示2: 创建评论 ---")
	comment := Comment{
		Content: "这是一条测试评论",
		Status:  "approved",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	
	if err := db.Create(&comment).Error; err != nil {
		log.Printf("创建评论失败: %v\n", err)
	} else {
		fmt.Printf("评论创建成功: ID %d\n", comment.ID)
	}
	
	// 检查文章评论状态
	postStats, _ := GetPostStats(post.ID)
	fmt.Printf("文章当前评论数量: %d, 评论状态: %s\n", postStats.CommentCount, postStats.CommentStatus)
	
	// 演示3: 删除评论（触发Comment的AfterDelete钩子）
	fmt.Println("\n--- 演示3: 删除评论 ---")
	if err := DeleteComment(comment.ID); err != nil {
		log.Printf("删除评论失败: %v\n", err)
	} else {
		fmt.Println("评论删除成功")
	}
	
	// 再次检查文章评论状态
	postStats, _ = GetPostStats(post.ID)
	fmt.Printf("删除后文章评论数量: %d, 评论状态: %s\n", postStats.CommentCount, postStats.CommentStatus)
	
	// 演示4: 删除文章（触发Post的BeforeDelete钩子）
	fmt.Println("\n--- 演示4: 删除文章 ---")
	if err := db.Delete(&post).Error; err != nil {
		log.Printf("删除文章失败: %v\n", err)
	} else {
		fmt.Println("文章删除成功")
	}
	
	// 最终检查用户统计
	userStats, _ = GetUserStats(user.ID)
	fmt.Printf("最终用户文章数量: %d, 评论数量: %d\n", userStats.PostCount, userStats.CommentCount)
	
	fmt.Println("\n=== 演示完成 ===")


}

// 辅助函数：打印查询结果（用于调试）
func printQueryResult(result interface{}) {
	switch v := result.(type) {
	case []Post:
		for i, post := range v {
			fmt.Printf("文章 %d: ID=%d, Title=%s, Comments=%d\n", 
				i+1, post.ID, post.Title, len(post.Comments))
		}
	case map[string]interface{}:
		for key, value := range v {
			fmt.Printf("%s: %v\n", key, value)
		}
	case []map[string]interface{}:
		for i, item := range v {
			fmt.Printf("结果 %d:\n", i+1)
			for key, value := range item {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}
}