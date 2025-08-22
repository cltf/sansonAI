package models

import (
	"database/sql"
	"time"
	"aiforum/utils"
)

// 创建用户
func CreateUser(username, email, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, hashedPassword)
	return err
}

// 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT id, username, email, password, avatar, level, points, created_at, updated_at FROM users WHERE username = ?",
		username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Level, &user.Points, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 根据ID获取用户
func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := DB.QueryRow("SELECT id, username, email, password, avatar, level, points, created_at, updated_at FROM users WHERE id = ?",
		id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Level, &user.Points, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 检查用户名是否存在
func UsernameExists(username string) (bool, error) {
	var exists int
	err := DB.QueryRow("SELECT 1 FROM users WHERE username = ?", username).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// 检查邮箱是否存在
func EmailExists(email string) (bool, error) {
	var exists int
	err := DB.QueryRow("SELECT 1 FROM users WHERE email = ?", email).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// 更新用户信息
func UpdateUser(id int, username, email, avatar string) error {
	_, err := DB.Exec("UPDATE users SET username = ?, email = ?, avatar = ?, updated_at = ? WHERE id = ?",
		username, email, avatar, time.Now(), id)
	return err
}

// 更新用户积分
func UpdateUserPoints(id, points int) error {
	_, err := DB.Exec("UPDATE users SET points = points + ?, updated_at = ? WHERE id = ?",
		points, time.Now(), id)
	return err
}

// 获取用户等级
func GetUserLevel(points int) int {
	if points >= 1000 {
		return 5
	} else if points >= 500 {
		return 4
	} else if points >= 200 {
		return 3
	} else if points >= 50 {
		return 2
	}
	return 1
} 