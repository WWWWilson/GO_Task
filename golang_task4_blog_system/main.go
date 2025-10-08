package main

import (
	"golang_task4_blog_system/controllers"
	"golang_task4_blog_system/database"
	"golang_task4_blog_system/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

// JWTSecret 用于签名JWT令牌
const JWTSecret = "your-super-secret-jwt-key-change-in-production"

func main() {
  // 初始化数据库连接
	dbConfig := &database.MySQLConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "wilson1234",
		Database: "blog_system",
		MaxIdleConns:    100,
	    MaxOpenConns:100,
	}

	database.InitDB(dbConfig)
	defer database.CloseDB()

    router := gin.Default()

	// 公开路由
	public := router.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
		public.GET("/posts", controllers.GetPosts)
		public.GET("/posts/:id", controllers.GetPost)
	}

	// 需要认证的路由
	auth := router.Group("/api")
	auth.Use(middleware.BasicAuth())
	{
		// 文章管理
		auth.POST("/posts", controllers.CreatePost)
		auth.PUT("/posts/:id", controllers.UpdatePost)
		auth.DELETE("/posts/:id", controllers.DeletePost)
		
		// 评论管理
		auth.POST("/comments", controllers.CreateComment)    // 创建评论（需要认证）
		auth.GET("/comments/:id", controllers.GetComment)    // 获取评论详情（需要认证）
		auth.PUT("/comments/:id", controllers.UpdateComment) // 更新评论（需要认证+作者权限）
		auth.DELETE("/comments/:id", controllers.DeleteComment) // 删除评论（需要认证+作者权限）
		auth.GET("/comments/my", controllers.GetMyComments)  // 获取我的评论（需要认证）
	}

  // 启动服务器
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}