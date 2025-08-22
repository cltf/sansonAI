package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// 获取用户动态
func GetUserActivity(c *gin.Context) {
	userID := c.GetInt("user_id")
	filter := c.Query("filter")

	activities, err := models.GetUserActivity(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户动态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"activities": activities,
	})
}

// 获取用户提问
func GetUserQuestions(c *gin.Context) {
	userID := c.GetInt("user_id")

	questions, err := models.GetUserQuestions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户提问失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"questions": questions,
	})
}

// 获取用户回答
func GetUserAnswers(c *gin.Context) {
	userID := c.GetInt("user_id")

	answers, err := models.GetUserAnswers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户回答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"answers": answers,
	})
}

// 获取用户分享
func GetUserShares(c *gin.Context) {
	userID := c.GetInt("user_id")

	shares, err := models.GetUserShares(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户分享失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"shares":  shares,
	})
}

// 获取用户资料
func GetUserResources(c *gin.Context) {
	userID := c.GetInt("user_id")

	resources, err := models.GetUserResources(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户资料失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"resources": resources,
	})
}

// 获取用户收藏
func GetUserFavorites(c *gin.Context) {
	userID := c.GetInt("user_id")

	favorites, err := models.GetUserFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"favorites": favorites,
	})
}

// 获取用户关注
func GetUserFollowing(c *gin.Context) {
	userID := c.GetInt("user_id")

	following, err := models.GetUserFollowing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"following": following,
	})
}

// 获取用户粉丝
func GetUserFollowers(c *gin.Context) {
	userID := c.GetInt("user_id")

	followers, err := models.GetUserFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户粉丝失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"followers": followers,
	})
}

// 获取用户消息
func GetUserMessages(c *gin.Context) {
	userID := c.GetInt("user_id")

	messages, err := models.GetUserMessages(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"messages": messages,
	})
}

// 更新用户头像
func UpdateUserAvatar(c *gin.Context) {
	userID := c.GetInt("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请选择头像文件",
		})
		return
	}

	// 检查文件类型
	if !isValidImageFile(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请选择有效的图片文件（JPG、PNG、GIF）",
		})
		return
	}

	// 检查文件大小（限制为5MB）
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "头像文件大小不能超过5MB",
		})
		return
	}

	// 生成文件名
	filename := generateAvatarFilename(userID, file.Filename)
	filepath := "images/avatars/" + filename

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "头像上传失败",
		})
		return
	}

	// 更新数据库中的头像路径
	avatarURL := "/" + filepath
	err = models.UpdateUserAvatar(userID, avatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "更新头像信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"avatar_url": avatarURL,
		"message":    "头像上传成功",
	})
}

// 更新用户个人资料
func UpdateUserProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Username      string `json:"username" binding:"required"`
		Email         string `json:"email" binding:"required,email"`
		Bio           string `json:"bio"`
		Phone         string `json:"phone"`
		Website       string `json:"website"`
		ProfilePublic bool   `json:"profile_public"`
		ShowEmail     bool   `json:"show_email"`
		ShowPhone     bool   `json:"show_phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请填写完整的用户信息",
		})
		return
	}

	// 检查用户名是否已被其他用户使用
	user, err := models.GetUserByUsername(req.Username)
	if err == nil && user.ID != userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "用户名已被使用",
		})
		return
	}

	// 更新用户信息
	err = models.UpdateUserProfile(userID, req.Username, req.Email, req.Bio, req.Phone, req.Website, req.ProfilePublic, req.ShowEmail, req.ShowPhone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "更新个人资料失败",
		})
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    updatedUser,
		"message": "个人资料更新成功",
	})
}

// 修改用户密码
func ChangeUserPassword(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请填写完整的密码信息",
		})
		return
	}

	// 验证当前密码
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取用户信息失败",
		})
		return
	}

	if !models.CheckPassword(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "当前密码错误",
		})
		return
	}

	// 更新密码
	err = models.UpdateUserPassword(userID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "密码修改失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码修改成功",
	})
}

// 保存用户通知设置
func SaveNotificationSettings(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		EmailNotifications    bool `json:"email_notifications"`
		BrowserNotifications  bool `json:"browser_notifications"`
		QuestionNotifications bool `json:"question_notifications"`
		FollowNotifications   bool `json:"follow_notifications"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "请填写完整的通知设置",
		})
		return
	}

	err := models.UpdateNotificationSettings(userID, req.EmailNotifications, req.BrowserNotifications, req.QuestionNotifications, req.FollowNotifications)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "保存通知设置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "通知设置保存成功",
	})
}

// 标记所有消息为已读
func MarkAllMessagesRead(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := models.MarkAllMessagesRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "标记已读失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已标记所有消息为已读",
	})
}

// 标记单条消息为已读
func MarkMessageRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的消息ID",
		})
		return
	}

	err = models.MarkMessageRead(userID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "标记已读失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已标记为已读",
	})
}

// 删除消息
func DeleteMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的消息ID",
		})
		return
	}

	err = models.DeleteMessage(userID, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除消息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息删除成功",
	})
}

// 关注用户
func FollowUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的用户ID",
		})
		return
	}

	if userID == targetUserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "不能关注自己",
		})
		return
	}

	err = models.FollowUser(userID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "关注成功",
	})
}

// 取消关注用户
func UnfollowUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的用户ID",
		})
		return
	}

	err = models.UnfollowUser(userID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "取消关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已取消关注",
	})
}

// 取消收藏
func RemoveFavorite(c *gin.Context) {
	userID := c.GetInt("user_id")
	favoriteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的收藏ID",
		})
		return
	}

	err = models.RemoveFavorite(userID, favoriteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "取消收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "已取消收藏",
	})
}

// 删除用户提问
func DeleteUserQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的问题ID",
		})
		return
	}

	err = models.DeleteUserQuestion(userID, questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除问题失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "问题删除成功",
	})
}

// 删除用户回答
func DeleteUserAnswer(c *gin.Context) {
	userID := c.GetInt("user_id")
	answerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的回答ID",
		})
		return
	}

	err = models.DeleteUserAnswer(userID, answerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除回答失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回答删除成功",
	})
}

// 删除用户分享
func DeleteUserShare(c *gin.Context) {
	userID := c.GetInt("user_id")
	shareID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的分享ID",
		})
		return
	}

	err = models.DeleteUserShare(userID, shareID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除分享失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分享删除成功",
	})
}

// 删除用户资料
func DeleteUserResource(c *gin.Context) {
	userID := c.GetInt("user_id")
	resourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "无效的资料ID",
		})
		return
	}

	err = models.DeleteUserResource(userID, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除资料失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "资料删除成功",
	})
}

// 工具函数
func isValidImageFile(filename string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range validExtensions {
		if len(filename) >= len(ext) && filename[len(filename)-len(ext):] == ext {
			return true
		}
	}
	return false
}

func generateAvatarFilename(userID int, originalFilename string) string {
	timestamp := time.Now().Unix()
	ext := ""
	if len(originalFilename) >= 4 {
		ext = originalFilename[len(originalFilename)-4:]
	}
	return strconv.Itoa(userID) + "_" + strconv.FormatInt(timestamp, 10) + ext
} 