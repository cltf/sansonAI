package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// LearningResourcesPage 学习资料页面
func LearningResourcesPage(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	keyword := c.Query("keyword")
	resourceType := c.Query("type")
	level := c.Query("level")
	timeFilter := c.Query("time")
	rating := c.Query("rating")
	category := c.Query("category")
	
	// 设置默认值
	if page < 1 {
		page = 1
	}
	limit := 12
	
	// 获取学习资料列表
	resources, total, err := models.GetLearningResources(page, limit, keyword, resourceType, level, timeFilter, rating, category)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取资料列表失败",
		})
		return
	}
	
	// 获取最新上传资料
	latestResources, err := models.GetLatestResources(5)
	if err != nil {
		latestResources = []models.LearningResource{}
	}
	
	// 获取最高评分资料
	topRatedResources, err := models.GetTopRatedResources(5)
	if err != nil {
		topRatedResources = []models.LearningResource{}
	}
	
	// 获取当前用户信息（如果已登录）
	var user *models.User
	if userID, exists := c.Get("user_id"); exists {
		user, _ = models.GetUserByID(userID.(int))
	}
	
	// 计算总页数
	totalPages := (total + limit - 1) / limit
	
	c.HTML(http.StatusOK, "learning_resources.html", gin.H{
		"title":              "学习资料",
		"resources":          resources,
		"latestResources":    latestResources,
		"topRatedResources":  topRatedResources,
		"currentPage":        page,
		"totalPages":         totalPages,
		"keyword":            keyword,
		"type":               resourceType,
		"level":              level,
		"time":               timeFilter,
		"rating":             rating,
		"currentCategory":    category,
		"user":               user,
	})
}

// CategoryPage 分类页面
func CategoryPage(c *gin.Context) {
	category := c.Param("category")
	
	// 获取分类信息
	categoryInfo, err := models.GetCategoryBySlug(category)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "分类不存在",
		})
		return
	}
	
	// 获取该分类下的资料
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 12
	
	resources, total, err := models.GetLearningResources(page, limit, "", "", "", "", "", category)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取资料列表失败",
		})
		return
	}
	
	// 获取最新上传资料
	latestResources, err := models.GetLatestResources(5)
	if err != nil {
		latestResources = []models.LearningResource{}
	}
	
	// 获取最高评分资料
	topRatedResources, err := models.GetTopRatedResources(5)
	if err != nil {
		topRatedResources = []models.LearningResource{}
	}
	
	// 计算总页数
	totalPages := (total + limit - 1) / limit
	
	c.HTML(http.StatusOK, "learning_resources.html", gin.H{
		"title":              categoryInfo.Name,
		"resources":          resources,
		"latestResources":    latestResources,
		"topRatedResources":  topRatedResources,
		"currentPage":        page,
		"totalPages":         totalPages,
		"currentCategory":    category,
		"categoryInfo":       categoryInfo,
	})
}

// UploadLearningResource 上传学习资料
func UploadLearningResource(c *gin.Context) {
	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	
	// 获取表单数据
	title := c.PostForm("title")
	description := c.PostForm("description")
	resourceType := c.PostForm("type")
	level := c.PostForm("level")
	category := c.PostForm("category")
	tags := c.PostForm("tags")
	
	// 验证必填字段
	if title == "" || description == "" || resourceType == "" || level == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写所有必填字段"})
		return
	}
	
	// 处理文件上传
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	
	// 处理封面图片
	var coverImage string
	coverFiles := form.File["cover"]
	if len(coverFiles) > 0 {
		coverFile := coverFiles[0]
		coverImage, err = models.UploadFile(coverFile, "covers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "封面图片上传失败"})
			return
		}
	}
	
	// 上传文件
	var filePaths []string
	var totalSize int64
	
	for _, file := range files {
		filePath, err := models.UploadFile(file, "resources")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
			return
		}
		filePaths = append(filePaths, filePath)
		totalSize += file.Size
	}
	
	// 创建学习资料记录
	resourceID, err := models.CreateLearningResource(title, description, resourceType, level, category, tags, coverImage, filePaths, totalSize, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建资料记录失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "资料上传成功",
		"resourceID": resourceID,
	})
}

// DownloadLearningResource 下载学习资料
func DownloadLearningResource(c *gin.Context) {
	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	
	resourceID := c.Param("id")
	
	// 获取资料信息
	resource, err := models.GetLearningResourceByID(resourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "资料不存在"})
		return
	}
	
	// 检查用户是否有权限下载（这里可以添加权限检查逻辑）
	
	// 记录下载次数
	err = models.IncrementResourceDownloads(resourceID, userID.(int))
	if err != nil {
		// 记录失败不影响下载
	}
	
	// 返回下载链接
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"downloadUrl": resource.DownloadURL,
		"filename": resource.Title,
	})
}

// RateLearningResource 评分学习资料
func RateLearningResource(c *gin.Context) {
	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	
	resourceID := c.Param("id")
	
	var req struct {
		Rating int `json:"rating"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	
	if req.Rating < 1 || req.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在1-5之间"})
		return
	}
	
	// 提交评分
	err := models.RateLearningResource(resourceID, userID.(int), req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评分失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "评分成功",
	})
}

// CommentLearningResource 评论学习资料
func CommentLearningResource(c *gin.Context) {
	// 检查用户是否已登录
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}
	
	resourceID := c.Param("id")
	
	var req struct {
		Content string `json:"content"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}
	
	// 提交评论
	err := models.CommentLearningResource(resourceID, userID.(int), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "评论成功",
	})
} 