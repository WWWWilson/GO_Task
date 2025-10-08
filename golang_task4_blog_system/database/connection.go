package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// MySQLConfig 包含MySQL连接配置
type MySQLConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func InitDB(cfg *MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	log.Printf("Connecting to MySQL database: %s@%s:%s/%s", 
		cfg.User, cfg.Host, cfg.Port, cfg.Database)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("✅ MySQL connection established successfully")

	// Auto migrate tables
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get generic database object:", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database migrated successfully")
}

// 关闭数据库连接
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Println("Failed to get generic database object:", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}