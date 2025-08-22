package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// 获取分类列表
func GetCategories(c *gin.Context) {
	categories, err := models.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// 获取标签列表
func GetTags(c *gin.Context) {
	tags, err := models.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

// 获取用户资料页面
func ProfilePage(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取用户信息失败",
		})
		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"title": "个人资料",
		"user":  user,
	})
}

// 更新用户资料
func UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整的用户信息"})
		return
	}

	// 检查用户名是否已被其他用户使用
	user, err := models.GetUserByUsername(req.Username)
	if err == nil && user.ID != userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已被使用"})
		return
	}

	// 更新用户信息
	err = models.UpdateUser(userID, req.Username, req.Email, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
} 