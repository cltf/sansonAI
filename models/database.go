package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"aiforum/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Avatar    string    `json:"avatar"`
	Level     int       `json:"level"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 帖子模型
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CategoryID  int       `json:"category_id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	UserAvatar  string    `json:"user_avatar"`
	ViewCount   int       `json:"view_count"`
	ReplyCount  int       `json:"reply_count"`
	LikeCount   int       `json:"like_count"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 回复模型
type Reply struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	UserAvatar string   `json:"user_avatar"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// 分类模型
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PostCount   int    `json:"post_count"`
}

// 标签模型
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 初始化数据库
func InitDB() error {
	config.Init()
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 测试连接
	err = DB.Ping()
	if err != nil {
		return err
	}

	// 创建表
	err = createTables()
	if err != nil {
		return err
	}

	log.Println("数据库连接成功")
	return nil
}

// 创建数据表
func createTables() error {
	// 用户表
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		avatar VARCHAR(255) DEFAULT '/images/user.jpg',
		level INT DEFAULT 1,
		points INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 分类表
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		description TEXT,
		post_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 帖子表
	postTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		category_id INT NOT NULL,
		user_id INT NOT NULL,
		view_count INT DEFAULT 0,
		reply_count INT DEFAULT 0,
		like_count INT DEFAULT 0,
		tags VARCHAR(500),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 回复表
	replyTable := `
	CREATE TABLE IF NOT EXISTS replies (
		id INT AUTO_INCREMENT PRIMARY KEY,
		post_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 问题表
	questionTable := `
	CREATE TABLE IF NOT EXISTS questions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		category_id INT NOT NULL,
		user_id INT NOT NULL,
		view_count INT DEFAULT 0,
		answer_count INT DEFAULT 0,
		like_count INT DEFAULT 0,
		tags VARCHAR(500),
		reward INT DEFAULT 0,
		is_solved BOOLEAN DEFAULT FALSE,
		summary TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 回答表
	answerTable := `
	CREATE TABLE IF NOT EXISTS answers (
		id INT AUTO_INCREMENT PRIMARY KEY,
		question_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		like_count INT DEFAULT 0,
		is_accepted BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 回答点赞表
	answerLikeTable := `
	CREATE TABLE IF NOT EXISTS answer_likes (
		id INT AUTO_INCREMENT PRIMARY KEY,
		answer_id INT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY unique_like (answer_id, user_id),
		FOREIGN KEY (answer_id) REFERENCES answers(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	// 标签表
	tagTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	tables := []string{userTable, categoryTable, postTable, replyTable, questionTable, answerTable, answerLikeTable, tagTable}
	
	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	// 插入默认分类
	insertDefaultCategories()
	
	return nil
}

// 插入默认分类
func insertDefaultCategories() {
	categories := []string{
		"知识问答",
		"技术分享", 
		"学习资料",
		"社群交流",
	}

	for _, name := range categories {
		_, err := DB.Exec("INSERT IGNORE INTO categories (name) VALUES (?)", name)
		if err != nil {
			log.Printf("插入分类失败: %v", err)
		}
	}
} 