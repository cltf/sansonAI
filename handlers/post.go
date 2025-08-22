package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// 首页
func HomePage(c *gin.Context) {
	// 获取最新帖子
	posts, err := models.GetPosts(1, 10, 0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取帖子失败",
		})
		return
	}

	// 获取热门帖子
	hotPosts, err := models.GetHotPosts(5)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取热门帖子失败",
		})
		return
	}

	// 获取分类
	categories, err := models.GetCategories()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取分类失败",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "AI论坛 - 知识分享与交流平台",
		"posts":      posts,
		"hotPosts":   hotPosts,
		"categories": categories,
	})
}

// 发帖页面
func NewPostPage(c *gin.Context) {
	categories, err := models.GetCategories()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取分类失败",
		})
		return
	}

	c.HTML(http.StatusOK, "new_post.html", gin.H{
		"title":      "发布新帖",
		"categories": categories,
	})
}

// 创建帖子
func CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID int    `json:"category_id" binding:"required"`
		Tags       string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整的帖子信息"})
		return
	}

	// 创建帖子
	postID, err := models.CreatePost(req.Title, req.Content, req.CategoryID, userID, req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发帖失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "发帖成功",
		"post_id": postID,
	})
}

// 查看帖子
func ViewPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "无效的帖子ID",
		})
		return
	}

	// 获取帖子信息
	post, err := models.GetPostByID(postID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "帖子不存在",
		})
		return
	}

	// 增加浏览量
	models.IncrementPostViewCount(postID)

	// 获取回复
	replies, err := models.GetRepliesByPostID(postID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取回复失败",
		})
		return
	}

	c.HTML(http.StatusOK, "post.html", gin.H{
		"title":  post.Title,
		"post":   post,
		"replies": replies,
	})
}

// 创建回复
func CreateReply(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写回复内容"})
		return
	}

	// 创建回复
	replyID, err := models.CreateReply(postID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "回复成功",
		"reply_id": replyID,
	})
}

// 获取帖子列表（API）
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))

	posts, err := models.GetPosts(page, limit, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取帖子失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"page":  page,
		"limit": limit,
	})
}

// 获取单个帖子（API）
func GetPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的帖子ID"})
		return
	}

	post, err := models.GetPostByID(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// 搜索帖子
func Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := models.SearchPosts(keyword, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":   posts,
		"keyword": keyword,
		"page":    page,
		"limit":   limit,
	})
} 