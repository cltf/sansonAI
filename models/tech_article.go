package models

import (
	"database/sql"
	"strings"
	"time"
)

// 技术文章模型
type TechArticle struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Summary      string    `json:"summary"`
	Category     string    `json:"category"`
	CategoryName string    `json:"category_name"`
	UserID       int       `json:"user_id"`
	AuthorName   string    `json:"author_name"`
	AuthorAvatar string    `json:"author_avatar"`
	AuthorLevel  int       `json:"author_level"`
	AuthorBio    string    `json:"author_bio"`
	CoverImage   string    `json:"cover_image"`
	Tags         string    `json:"tags"`
	TagsArray    []string  `json:"tags_array"`
	ViewCount    int       `json:"view_count"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	TopicSlug    string    `json:"topic_slug"`
	IsLiked      bool      `json:"is_liked"`
	IsFavorited  bool      `json:"is_favorited"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 评论模型
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	UserAvatar string   `json:"user_avatar"`
	UserLevel int       `json:"user_level"`
	LikeCount int       `json:"like_count"`
	IsLiked   bool      `json:"is_liked"`
	Replies   []Comment `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
}

// 专题模型
type Topic struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ArticleCount int   `json:"article_count"`
}

// 热门作者模型
type PopularAuthor struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	ArticleCount  int    `json:"article_count"`
	FollowerCount int    `json:"follower_count"`
}

// 创建技术文章
func CreateTechArticle(title, content, category string, userID int, tags, coverImage string) (int, error) {
	// 生成文章摘要
	summary := generateTechArticleSummary(content)
	
	result, err := DB.Exec(`
		INSERT INTO tech_articles (title, content, summary, category, user_id, cover_image, tags) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, title, content, summary, category, userID, coverImage, tags)
	
	if err != nil {
		return 0, err
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 给发布用户加积分
	err = UpdateUserPoints(userID, 10)
	if err != nil {
		return 0, err
	}

	return int(articleID), nil
}

// 获取技术文章列表
func GetTechArticles(page, limit int, category, keyword, sort, topic string) ([]*TechArticle, error) {
	offset := (page - 1) * limit
	
	var query string
	var args []interface{}
	
	baseQuery := `
		SELECT a.id, a.title, a.content, a.summary, a.category, a.user_id, 
			   u.username, u.avatar, a.cover_image, a.tags, a.view_count, 
			   a.like_count, a.comment_count, a.topic_slug, a.created_at, a.updated_at
		FROM tech_articles a
		JOIN users u ON a.user_id = u.id
	`
	
	var whereConditions []string
	
	if category != "" {
		whereConditions = append(whereConditions, "a.category = ?")
		args = append(args, category)
	}
	
	if keyword != "" {
		whereConditions = append(whereConditions, "(a.title LIKE ? OR a.content LIKE ? OR a.tags LIKE ?)")
		keyword = "%" + keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}
	
	if topic != "" {
		whereConditions = append(whereConditions, "a.topic_slug = ?")
		args = append(args, topic)
	}
	
	if len(whereConditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(whereConditions, " AND ")
	} else {
		query = baseQuery
	}
	
	// 排序
	switch sort {
	case "likes":
		query += " ORDER BY a.like_count DESC"
	case "comments":
		query += " ORDER BY a.comment_count DESC"
	case "views":
		query += " ORDER BY a.view_count DESC"
	default:
		query += " ORDER BY a.created_at DESC"
	}
	
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*TechArticle
	for rows.Next() {
		article := &TechArticle{}
		err := rows.Scan(
			&article.ID, &article.Title, &article.Content, &article.Summary, &article.Category, &article.UserID,
			&article.AuthorName, &article.AuthorAvatar, &article.CoverImage, &article.Tags, &article.ViewCount,
			&article.LikeCount, &article.CommentCount, &article.TopicSlug, &article.CreatedAt, &article.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// 设置分类名称
		article.CategoryName = getTechCategoryName(article.Category)
		
		articles = append(articles, article)
	}

	return articles, nil
}

// 根据ID获取技术文章
func GetTechArticleByID(id int) (*TechArticle, error) {
	article := &TechArticle{}
	err := DB.QueryRow(`
		SELECT a.id, a.title, a.content, a.summary, a.category, a.user_id, 
			   u.username, u.avatar, a.cover_image, a.tags, a.view_count, 
			   a.like_count, a.comment_count, a.topic_slug, a.created_at, a.updated_at
		FROM tech_articles a
		JOIN users u ON a.user_id = u.id
		WHERE a.id = ?
	`, id).Scan(
		&article.ID, &article.Title, &article.Content, &article.Summary, &article.Category, &article.UserID,
		&article.AuthorName, &article.AuthorAvatar, &article.CoverImage, &article.Tags, &article.ViewCount,
		&article.LikeCount, &article.CommentCount, &article.TopicSlug, &article.CreatedAt, &article.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	article.CategoryName = getTechCategoryName(article.Category)
	
	return article, nil
}

// 获取相关技术文章
func GetRelatedTechArticles(articleID, limit int) ([]*TechArticle, error) {
	query := `
		SELECT a.id, a.title, a.summary, a.cover_image, a.view_count, a.like_count
		FROM tech_articles a
		WHERE a.id != ? AND a.category = (SELECT category FROM tech_articles WHERE id = ?)
		ORDER BY a.view_count DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, articleID, articleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*TechArticle
	for rows.Next() {
		article := &TechArticle{}
		err := rows.Scan(&article.ID, &article.Title, &article.Summary, &article.CoverImage, &article.ViewCount, &article.LikeCount)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 获取技术文章数量
func GetTechArticleCount(category, keyword, topic string) (int, error) {
	var count int
	var query string
	var args []interface{}
	
	baseQuery := "SELECT COUNT(*) FROM tech_articles"
	var whereConditions []string
	
	if category != "" {
		whereConditions = append(whereConditions, "category = ?")
		args = append(args, category)
	}
	
	if keyword != "" {
		whereConditions = append(whereConditions, "(title LIKE ? OR content LIKE ? OR tags LIKE ?)")
		keyword = "%" + keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}
	
	if topic != "" {
		whereConditions = append(whereConditions, "topic_slug = ?")
		args = append(args, topic)
	}
	
	if len(whereConditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(whereConditions, " AND ")
	} else {
		query = baseQuery
	}
	
	err := DB.QueryRow(query, args...).Scan(&count)
	return count, err
}

// 增加文章浏览量
func IncrementTechArticleViewCount(id int) error {
	_, err := DB.Exec("UPDATE tech_articles SET view_count = view_count + 1 WHERE id = ?", id)
	return err
}

// 点赞文章
func LikeTechArticle(articleID, userID int) error {
	// 检查是否已点赞
	var exists int
	err := DB.QueryRow("SELECT 1 FROM tech_article_likes WHERE article_id = ? AND user_id = ?", articleID, userID).Scan(&exists)
	if err == nil {
		// 已点赞，取消点赞
		_, err = DB.Exec("DELETE FROM tech_article_likes WHERE article_id = ? AND user_id = ?", articleID, userID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("UPDATE tech_articles SET like_count = like_count - 1 WHERE id = ?", articleID)
		return err
	} else if err == sql.ErrNoRows {
		// 未点赞，添加点赞
		_, err = DB.Exec("INSERT INTO tech_article_likes (article_id, user_id) VALUES (?, ?)", articleID, userID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("UPDATE tech_articles SET like_count = like_count + 1 WHERE id = ?", articleID)
		return err
	}
	
	return err
}

// 获取热门作者
func GetPopularAuthors(limit int) ([]*PopularAuthor, error) {
	query := `
		SELECT u.id, u.username, u.avatar,
		       COUNT(DISTINCT a.id) as article_count,
		       COUNT(DISTINCT f.follower_id) as follower_count
		FROM users u
		LEFT JOIN tech_articles a ON u.id = a.user_id
		LEFT JOIN user_follows f ON u.id = f.following_id
		GROUP BY u.id
		HAVING article_count > 0
		ORDER BY article_count DESC, follower_count DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*PopularAuthor
	for rows.Next() {
		author := &PopularAuthor{}
		err := rows.Scan(&author.ID, &author.Username, &author.Avatar, &author.ArticleCount, &author.FollowerCount)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

// 获取相关专题
func GetRelatedTopics(limit int) ([]*Topic, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.icon,
		       COUNT(a.id) as article_count
		FROM topics t
		LEFT JOIN tech_articles a ON t.slug = a.topic_slug
		GROUP BY t.id
		ORDER BY article_count DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*Topic
	for rows.Next() {
		topic := &Topic{}
		err := rows.Scan(&topic.ID, &topic.Name, &topic.Slug, &topic.Description, &topic.Icon, &topic.ArticleCount)
		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

// 根据slug获取专题
func GetTopicBySlug(slug string) (*Topic, error) {
	topic := &Topic{}
	err := DB.QueryRow(`
		SELECT id, name, slug, description, icon
		FROM topics
		WHERE slug = ?
	`, slug).Scan(&topic.ID, &topic.Name, &topic.Slug, &topic.Description, &topic.Icon)
	
	if err != nil {
		return nil, err
	}
	
	// 获取文章数量
	err = DB.QueryRow("SELECT COUNT(*) FROM tech_articles WHERE topic_slug = ?", slug).Scan(&topic.ArticleCount)
	if err != nil {
		return nil, err
	}
	
	return topic, nil
}

// 关注/取消关注作者
func ToggleFollowAuthor(followingID, followerID int) (bool, error) {
	// 检查是否已关注
	var exists int
	err := DB.QueryRow("SELECT 1 FROM user_follows WHERE following_id = ? AND follower_id = ?", followingID, followerID).Scan(&exists)
	if err == nil {
		// 已关注，取消关注
		_, err = DB.Exec("DELETE FROM user_follows WHERE following_id = ? AND follower_id = ?", followingID, followerID)
		return false, err
	} else if err == sql.ErrNoRows {
		// 未关注，添加关注
		_, err = DB.Exec("INSERT INTO user_follows (following_id, follower_id) VALUES (?, ?)", followingID, followerID)
		return true, err
	}
	
	return false, err
}

// 生成文章摘要
func generateTechArticleSummary(content string) string {
	if len(content) <= 200 {
		return content
	}
	return content[:200] + "..."
}

// 获取分类名称
func getTechCategoryName(category string) string {
	categoryMap := map[string]string{
		"algorithm": "算法研究",
		"development": "应用开发",
		"industry": "行业动态",
		"tutorial": "教程指南",
		"research": "研究论文",
	}
	
	if name, exists := categoryMap[category]; exists {
		return name
	}
	return category
}

// 根据ID获取技术文章详情（字符串ID版本）
func GetTechArticleByIDString(articleID string) (*TechArticle, error) {
	article := &TechArticle{}
	err := DB.QueryRow(`
		SELECT a.id, a.title, a.content, a.summary, a.category, a.user_id, 
			   u.username, u.avatar, u.bio, a.cover_image, a.tags, a.view_count, 
			   a.like_count, a.comment_count, a.topic_slug, a.created_at, a.updated_at
		FROM tech_articles a
		JOIN users u ON a.user_id = u.id
		WHERE a.id = ?
	`, articleID).Scan(&article.ID, &article.Title, &article.Content, &article.Summary, 
		&article.Category, &article.UserID, &article.AuthorName, &article.AuthorAvatar, 
		&article.AuthorBio, &article.CoverImage, &article.Tags, &article.ViewCount, 
		&article.LikeCount, &article.CommentCount, &article.TopicSlug, 
		&article.CreatedAt, &article.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	// 处理标签数组
	if article.Tags != "" {
		article.TagsArray = strings.Split(article.Tags, ",")
		for i, tag := range article.TagsArray {
			article.TagsArray[i] = strings.TrimSpace(tag)
		}
	}
	
	// 获取分类名称
	article.CategoryName = getTechCategoryName(article.Category)
	
	return article, nil
}

// 检查用户是否已点赞文章
func IsArticleLiked(articleID, userID int) bool {
	var exists int
	err := DB.QueryRow("SELECT 1 FROM tech_article_likes WHERE article_id = ? AND user_id = ?", articleID, userID).Scan(&exists)
	return err == nil
}

// 检查用户是否已收藏文章
func IsArticleFavorited(articleID, userID int) bool {
	var exists int
	err := DB.QueryRow("SELECT 1 FROM article_favorites WHERE article_id = ? AND user_id = ?", articleID, userID).Scan(&exists)
	return err == nil
}

// 获取文章评论
func GetArticleComments(articleID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.content, c.user_id, u.username, u.avatar, c.like_count, c.created_at
		FROM article_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = ? AND c.parent_id IS NULL
		ORDER BY c.created_at DESC
	`
	
	rows, err := DB.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		comment := Comment{}
		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, 
			&comment.Username, &comment.UserAvatar, &comment.LikeCount, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		// 获取回复
		replies, err := GetCommentReplies(comment.ID)
		if err == nil {
			comment.Replies = replies
		}
		
		comments = append(comments, comment)
	}

	return comments, nil
}

// 获取评论回复
func GetCommentReplies(commentID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.content, c.user_id, u.username, u.avatar, c.like_count, c.created_at
		FROM article_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.parent_id = ?
		ORDER BY c.created_at ASC
	`
	
	rows, err := DB.Query(query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Comment
	for rows.Next() {
		reply := Comment{}
		err := rows.Scan(&reply.ID, &reply.Content, &reply.UserID, 
			&reply.Username, &reply.UserAvatar, &reply.LikeCount, &reply.CreatedAt)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, nil
}

// 检查用户是否已点赞评论
func IsCommentLiked(commentID, userID int) bool {
	var exists int
	err := DB.QueryRow("SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&exists)
	return err == nil
}

// 获取相关文章
func GetRelatedArticles(articleID int, category string, limit int) ([]TechArticle, error) {
	query := `
		SELECT a.id, a.title, a.cover_image, a.user_id, u.username, a.view_count, a.created_at
		FROM tech_articles a
		JOIN users u ON a.user_id = u.id
		WHERE a.id != ? AND a.category = ?
		ORDER BY a.view_count DESC, a.created_at DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, articleID, category, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []TechArticle
	for rows.Next() {
		article := TechArticle{}
		err := rows.Scan(&article.ID, &article.Title, &article.CoverImage, 
			&article.UserID, &article.AuthorName, &article.ViewCount, &article.CreatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 获取作者其他文章
func GetAuthorArticles(userID, excludeArticleID int, limit int) ([]TechArticle, error) {
	query := `
		SELECT a.id, a.title, a.cover_image, a.view_count, a.created_at
		FROM tech_articles a
		WHERE a.user_id = ? AND a.id != ?
		ORDER BY a.created_at DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, userID, excludeArticleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []TechArticle
	for rows.Next() {
		article := TechArticle{}
		err := rows.Scan(&article.ID, &article.Title, &article.CoverImage, 
			&article.ViewCount, &article.CreatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 增加文章阅读量
func IncrementArticleViews(articleID int) error {
	_, err := DB.Exec("UPDATE tech_articles SET view_count = view_count + 1 WHERE id = ?", articleID)
	return err
} 