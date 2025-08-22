package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"aiforum/models"
)

// 技术分享页面
func TechSharePage(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	category := c.Query("category")
	keyword := c.Query("q")
	sort := c.DefaultQuery("sort", "latest")
	topic := c.Query("topic")

	// 获取文章列表
	articles, err := models.GetTechArticles(page, 12, category, keyword, sort, topic)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取文章失败",
		})
		return
	}

	// 获取热门作者
	popularAuthors, _ := models.GetPopularAuthors(5)

	// 获取相关专题
	relatedTopics, _ := models.GetRelatedTopics(6)

	// 获取当前用户信息（如果已登录）
	var user *models.User
	if userID, exists := c.Get("user_id"); exists {
		user, _ = models.GetUserByID(userID.(int))
	}

	// 计算总页数
	totalCount, _ := models.GetTechArticleCount(category, keyword, topic)
	totalPages := (totalCount + 11) / 12

	c.HTML(http.StatusOK, "tech_share.html", gin.H{
		"title":           "技术分享",
		"articles":        articles,
		"popularAuthors":  popularAuthors,
		"relatedTopics":   relatedTopics,
		"user":            user,
		"category":        category,
		"keyword":         keyword,
		"sort":            sort,
		"currentTopic":    topic,
		"page":            page,
		"totalPages":      totalPages,
		"currentPage":     page,
	})
}

// TechShareDetailPage 技术分享详情页
func TechShareDetailPage(c *gin.Context) {
	articleID := c.Param("id")
	
	// 获取文章详情
	article, err := models.GetTechArticleByIDString(articleID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "文章不存在",
		})
		return
	}
	
	// 获取当前用户信息（如果已登录）
	var user *models.User
	if userID, exists := c.Get("user_id"); exists {
		user, _ = models.GetUserByID(userID.(int))
	}
	
	// 计算作者等级
	article.AuthorLevel = calculateUserLevel(article.UserID)
	
	// 检查用户是否已点赞和收藏
	if user != nil {
		article.IsLiked = models.IsArticleLiked(article.ID, user.ID)
		article.IsFavorited = models.IsArticleFavorited(article.ID, user.ID)
	}
	
	// 获取文章评论
	comments, err := models.GetArticleComments(article.ID)
	if err != nil {
		comments = []models.Comment{}
	}
	
	// 为评论添加用户等级和点赞状态
	for i := range comments {
		comments[i].UserLevel = calculateUserLevel(comments[i].UserID)
		if user != nil {
			comments[i].IsLiked = models.IsCommentLiked(comments[i].ID, user.ID)
		}
	}
	
	// 获取相关文章推荐
	relatedArticles, err := models.GetRelatedArticles(article.ID, article.Category, 3)
	if err != nil {
		relatedArticles = []models.TechArticle{}
	}
	
	// 获取作者其他文章
	authorArticles, err := models.GetAuthorArticles(article.UserID, article.ID, 3)
	if err != nil {
		authorArticles = []models.TechArticle{}
	}
	
	// 增加文章阅读量
	go models.IncrementArticleViews(article.ID)
	
	c.HTML(http.StatusOK, "tech_share_detail.html", gin.H{
		"title":           article.Title,
		"article":         article,
		"comments":        comments,
		"relatedArticles": relatedArticles,
		"authorArticles":  authorArticles,
		"user":            user,
	})
}

// 专题页面
func TopicPage(c *gin.Context) {
	topicSlug := c.Param("slug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// 获取专题信息
	topic, err := models.GetTopicBySlug(topicSlug)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "专题不存在",
		})
		return
	}

	// 获取专题下的文章
	articles, err := models.GetTechArticles(page, 12, "", "", "latest", topicSlug)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "获取文章失败",
		})
		return
	}

	// 计算总页数
	totalCount, _ := models.GetTechArticleCount("", "", topicSlug)
	totalPages := (totalCount + 11) / 12

	c.HTML(http.StatusOK, "topic.html", gin.H{
		"title":       topic.Name,
		"topic":       topic,
		"articles":    articles,
		"page":        page,
		"totalPages":  totalPages,
		"currentPage": page,
	})
}

// 发布技术分享页面
func PublishTechSharePage(c *gin.Context) {
	c.HTML(http.StatusOK, "publish_tech_share.html", gin.H{
		"title": "发布技术分享",
	})
}

// 发布技术分享
func PublishTechShare(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		Title    string `form:"title" binding:"required"`
		Content  string `form:"content" binding:"required"`
		Category string `form:"category" binding:"required"`
		Tags     string `form:"tags"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整的文章信息"})
		return
	}

	// 处理封面图片上传
	coverImage := ""
	if file, err := c.FormFile("cover"); err == nil {
		// 这里应该保存文件到服务器，这里简化处理
		coverImage = "/uploads/covers/" + file.Filename
	}

	// 创建文章
	articleID, err := models.CreateTechArticle(req.Title, req.Content, req.Category, userID, req.Tags, coverImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "发布成功",
		"article_id": articleID,
	})
}

// 点赞文章
func LikeTechArticle(c *gin.Context) {
	userID := c.GetInt("user_id")
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	err = models.LikeTechArticle(articleID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取更新后的点赞数
	article, err := models.GetTechArticleByID(articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "点赞成功",
		"likeCount": article.LikeCount,
	})
}

// 关注作者
func FollowAuthor(c *gin.Context) {
	userID := c.GetInt("user_id")
	authorID, err := strconv.Atoi(c.Param("author_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作者ID"})
		return
	}

	isFollowed, err := models.ToggleFollowAuthor(authorID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"isFollowed": isFollowed,
		"message": func() string {
			if isFollowed {
				return "关注成功"
			}
			return "取消关注成功"
		}(),
	})
} 