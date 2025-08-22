package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// 计算用户等级
func calculateUserLevel(userID int) int {
	user, err := models.GetUserByID(userID)
	if err != nil {
		return 1
	}
	return models.GetUserLevel(user.Points)
}

// 问答页面
func QAPage(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category", "0"))
	keyword := c.Query("q")
	tag := c.Query("tag")
	sort := c.DefaultQuery("sort", "latest")

	// 获取分类
	categories, err := models.GetCategories()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取分类失败",
		})
		return
	}

	// 获取标签
	tags, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取标签失败",
		})
		return
	}

	// 获取问答列表
	var questions []*models.Question
	var totalCount, solvedCount, unsolvedCount, todayCount int

	if keyword != "" {
		// 搜索问答
		questions, err = models.SearchQuestions(keyword, page, 10)
	} else if tag != "" {
		// 按标签筛选
		questions, err = models.GetQuestionsByTag(tag, page, 10)
	} else {
		// 获取问答列表
		questions, err = models.GetQuestions(page, 10, categoryID, sort)
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取问答失败",
		})
		return
	}

	// 获取统计信息
	totalCount, _ = models.GetQuestionCount(categoryID)
	solvedCount, _ = models.GetSolvedQuestionCount(categoryID)
	unsolvedCount, _ = models.GetUnsolvedQuestionCount(categoryID)
	todayCount, _ = models.GetTodayQuestionCount(categoryID)

	// 获取推荐问答
	pendingQuestions, _ := models.GetPendingQuestions(5)
	rewardQuestions, _ := models.GetHighRewardQuestions(5)

	// 获取分类名称
	var categoryName string
	if categoryID > 0 {
		for _, cat := range categories {
			if cat.ID == categoryID {
				categoryName = cat.Name
				break
			}
		}
	}

	// 计算总页数
	totalPages := (totalCount + 9) / 10

	c.HTML(http.StatusOK, "qa.html", gin.H{
		"title":             "知识问答",
		"questions":         questions,
		"categories":        categories,
		"tags":              tags,
		"categoryID":        categoryID,
		"categoryName":      categoryName,
		"keyword":           keyword,
		"tag":               tag,
		"sort":              sort,
		"page":              page,
		"totalPages":        totalPages,
		"totalCount":        totalCount,
		"solvedCount":       solvedCount,
		"unsolvedCount":     unsolvedCount,
		"todayCount":        todayCount,
		"pendingQuestions":  pendingQuestions,
		"rewardQuestions":   rewardQuestions,
	})
}

// 高级搜索
func AdvancedSearch(c *gin.Context) {
	// 获取搜索参数
	title := c.Query("title")
	content := c.Query("content")
	author := c.Query("author")
	status := c.Query("status")
	timeRange := c.Query("time")
	rewardRange := c.Query("reward")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// 构建搜索条件
	conditions := make(map[string]interface{})
	if title != "" {
		conditions["title"] = title
	}
	if content != "" {
		conditions["content"] = content
	}
	if author != "" {
		conditions["author"] = author
	}
	if status != "" {
		conditions["status"] = status
	}
	if timeRange != "" {
		conditions["time_range"] = timeRange
	}
	if rewardRange != "" {
		conditions["reward_range"] = rewardRange
	}

	// 执行高级搜索
	questions, err := models.AdvancedSearchQuestions(conditions, page, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"page":      page,
	})
}

// 提问页面
func AskQuestionPage(c *gin.Context) {
	categories, err := models.GetCategories()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取分类失败",
		})
		return
	}

	tags, err := models.GetTags()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取标签失败",
		})
		return
	}

	c.HTML(http.StatusOK, "ask.html", gin.H{
		"title":      "发布新问题",
		"categories": categories,
		"tags":       tags,
	})
}

// 发布问题
func AskQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Title      string `form:"title" binding:"required"`
		Content    string `form:"content" binding:"required"`
		CategoryID int    `form:"category_id" binding:"required"`
		Tags       string `form:"tags"`
		Reward     int    `form:"reward"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整的提问信息"})
		return
	}

	// 检查用户积分是否足够
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	if user.Points < req.Reward {
		c.JSON(http.StatusBadRequest, gin.H{"error": "积分不足，无法设置悬赏"})
		return
	}

	// 创建问题
	questionID, err := models.CreateQuestion(req.Title, req.Content, req.CategoryID, userID, req.Tags, req.Reward)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布问题失败"})
		return
	}

	// 扣除悬赏积分
	if req.Reward > 0 {
		err = models.UpdateUserPoints(userID, -req.Reward)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "扣除积分失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "问题发布成功",
		"question_id": questionID,
	})
}

// 查看问题详情
func ViewQuestion(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "无效的问题ID",
		})
		return
	}

	// 获取问题详情
	question, err := models.GetQuestionByID(questionID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "问题不存在",
		})
		return
	}

	// 增加浏览量
	models.IncrementQuestionViewCount(questionID)

	// 获取回答
	answers, err := models.GetAnswersByQuestionID(questionID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取回答失败",
		})
		return
	}

	// 获取相关推荐
	relatedQuestions, _ := models.GetRelatedQuestions(questionID, 5)

	// 获取当前用户信息（如果已登录）
	var user *models.User
	if userID, exists := c.Get("user_id"); exists {
		user, _ = models.GetUserByID(userID.(int))
	}

	// 计算用户等级
	question.UserLevel = calculateUserLevel(question.UserID)

	c.HTML(http.StatusOK, "question_detail.html", gin.H{
		"title":             question.Title,
		"question":          question,
		"answers":           answers,
		"relatedQuestions":  relatedQuestions,
		"user":              user,
	})
}

// 回答问题
func AnswerQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写回答内容"})
		return
	}

	// 创建回答
	answerID, err := models.CreateAnswer(questionID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回答失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回答成功",
		"answer_id": answerID,
	})
}

// 采纳回答
func AcceptAnswer(c *gin.Context) {
	userID := c.GetInt("user_id")
	answerID, err := strconv.Atoi(c.Param("answer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的回答ID"})
		return
	}

	// 检查权限并采纳回答
	err = models.AcceptAnswer(answerID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "回答已采纳",
	})
}

// 点赞回答
func LikeAnswer(c *gin.Context) {
	userID := c.GetInt("user_id")
	answerID, err := strconv.Atoi(c.Param("answer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的回答ID"})
		return
	}

	err = models.LikeAnswer(answerID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的点赞数
	answer, err := models.GetAnswerByID(answerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取回答信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "点赞成功",
		"likeCount": answer.LikeCount,
	})
}

// 收藏问题
func FavoriteQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	isFavorited, err := models.ToggleQuestionFavorite(questionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"isFavorited": isFavorited,
		"message": func() string {
			if isFavorited {
				return "收藏成功"
			}
			return "取消收藏成功"
		}(),
	})
}

// 举报问题
func ReportQuestion(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的问题ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写举报原因"})
		return
	}

	err = models.ReportQuestion(questionID, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "举报失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "举报成功",
	})
} 