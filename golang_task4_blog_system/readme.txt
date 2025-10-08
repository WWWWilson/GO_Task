个人博客系统后端
一个基于 Go 语言、Gin 框架和 GORM 开发的个人博客系统后端，提供完整的博客文章管理和评论功能，支持用户认证和权限控制。


🛠️ 技术栈
后端框架: Gin
ORM: GORM
数据库: MySQL
认证: HTTP BasicAuth

环境要求:
Go 1.21 或更高版本
MySQL 5.7 或更高版本


安装步骤
1. 克隆项目
bash
git clone <项目地址>
cd blog-system
2. 安装依赖
bash
go mod tidy
3. 数据库配置
创建 MySQL 数据库：

sql
CREATE DATABASE blog_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
4. 配置数据库连接
修改 database/db.go 文件中的数据库连接信息：

go
dsn := "用户名:密码@tcp(主机:端口)/blog_system?charset=utf8mb4&parseTime=True&loc=Local"
示例：

go
dsn := "root:password@tcp(localhost:3306)/blog_system?charset=utf8mb4&parseTime=True&loc=Local" 

启动方式
开发环境启动
bash
go run main.go
生产环境编译



📁 项目结构
text
blog-system/
├── main.go                 # 程序入口
├── go.mod                 # Go 模块文件
├── go.sum                 # 依赖校验文件
├── controllers/           # 控制器层
│   ├── auth.go           # 认证控制器
│   ├── post.go           # 文章控制器
│   └── comment.go        # 评论控制器
├── models/               # 数据模型
│   ├── user.go          # 用户模型
│   ├── post.go          # 文章模型
│   └── comment.go       # 评论模型
├── middleware/           # 中间件
│   └── auth.go          # 认证中间件
└── database/            # 数据库层
    └── db.go            # 数据库连接

日志查看
启动时会在控制台输出服务状态和 API 文档：

🚀 博客系统服务启动在 :8080 端口
📝 API 文档:
   公开接口:
     POST /api/register              - 用户注册
     ...
