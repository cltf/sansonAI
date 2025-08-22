package models

import (
	"database/sql"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 学习资料模型
type LearningResource struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	TypeIcon       string    `json:"type_icon"`
	Level          string    `json:"level"`
	DifficultyText string    `json:"difficulty_text"`
	Category       string    `json:"category"`
	CategoryName   string    `json:"category_name"`
	UserID         int       `json:"user_id"`
	UploaderName   string    `json:"uploader_name"`
	UploaderAvatar string    `json:"uploader_avatar"`
	CoverImage     string    `json:"cover_image"`
	FilePaths      string    `json:"file_paths"`
	FileSize       string    `json:"file_size"`
	TotalSize      int64     `json:"total_size"`
	Tags           string    `json:"tags"`
	TagsArray      []string  `json:"tags_array"`
	Rating         float64   `json:"rating"`
	DownloadCount  int       `json:"download_count"`
	CommentCount   int       `json:"comment_count"`
	DownloadURL    string    `json:"download_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 分类模型
type ResourceCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Count       int    `json:"count"`
}

// 创建学习资料
func CreateLearningResource(title, description, resourceType, level, category, tags, coverImage string, filePaths []string, totalSize int64, userID int) (int, error) {
	// 生成文件路径字符串
	filePathsStr := strings.Join(filePaths, ",")
	
	// 生成下载URL
	downloadURL := "/downloads/" + filepath.Base(filePaths[0])
	
	result, err := DB.Exec(`
		INSERT INTO learning_resources (title, description, type, level, category, user_id, cover_image, file_paths, total_size, tags, download_url) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, title, description, resourceType, level, category, userID, coverImage, filePathsStr, totalSize, tags, downloadURL)
	
	if err != nil {
		return 0, err
	}

	resourceID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 给上传用户加积分
	err = UpdateUserPoints(userID, 20)
	if err != nil {
		return 0, err
	}

	return int(resourceID), nil
}

// 获取学习资料列表
func GetLearningResources(page, limit int, keyword, resourceType, level, timeFilter, rating, category string) ([]*LearningResource, int, error) {
	offset := (page - 1) * limit
	
	var query string
	var countQuery string
	var args []interface{}
	
	baseQuery := `
		SELECT r.id, r.title, r.description, r.type, r.level, r.category, r.user_id, 
			   u.username, u.avatar, r.cover_image, r.file_paths, r.total_size, r.tags, 
			   r.rating, r.download_count, r.comment_count, r.download_url, r.created_at, r.updated_at
		FROM learning_resources r
		JOIN users u ON r.user_id = u.id
	`
	
	countQuery = `
		SELECT COUNT(*) FROM learning_resources r
		JOIN users u ON r.user_id = u.id
	`
	
	var whereConditions []string
	
	if keyword != "" {
		whereConditions = append(whereConditions, "(r.title LIKE ? OR r.description LIKE ? OR r.tags LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if resourceType != "" {
		whereConditions = append(whereConditions, "r.type = ?")
		args = append(args, resourceType)
	}
	
	if level != "" {
		whereConditions = append(whereConditions, "r.level = ?")
		args = append(args, level)
	}
	
	if category != "" {
		whereConditions = append(whereConditions, "r.category = ?")
		args = append(args, category)
	}
	
	if timeFilter != "" {
		switch timeFilter {
		case "today":
			whereConditions = append(whereConditions, "DATE(r.created_at) = CURDATE()")
		case "week":
			whereConditions = append(whereConditions, "r.created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)")
		case "month":
			whereConditions = append(whereConditions, "r.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)")
		case "year":
			whereConditions = append(whereConditions, "r.created_at >= DATE_SUB(NOW(), INTERVAL 1 YEAR)")
		}
	}
	
	if rating != "" {
		ratingInt, _ := strconv.Atoi(rating)
		whereConditions = append(whereConditions, "r.rating >= ?")
		args = append(args, float64(ratingInt))
	}
	
	if len(whereConditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(whereConditions, " AND ")
		countQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	} else {
		query = baseQuery
	}
	
	query += " ORDER BY r.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	// 获取总数
	var total int
	err := DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	// 获取资料列表
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var resources []*LearningResource
	for rows.Next() {
		resource := &LearningResource{}
		err := rows.Scan(&resource.ID, &resource.Title, &resource.Description, &resource.Type, 
			&resource.Level, &resource.Category, &resource.UserID, &resource.UploaderName, 
			&resource.UploaderAvatar, &resource.CoverImage, &resource.FilePaths, &resource.TotalSize, 
			&resource.Tags, &resource.Rating, &resource.DownloadCount, &resource.CommentCount, 
			&resource.DownloadURL, &resource.CreatedAt, &resource.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		
		// 处理标签数组
		if resource.Tags != "" {
			resource.TagsArray = strings.Split(resource.Tags, ",")
			for i, tag := range resource.TagsArray {
				resource.TagsArray[i] = strings.TrimSpace(tag)
			}
		}
		
		// 设置类型图标
		resource.TypeIcon = getResourceTypeIcon(resource.Type)
		
		// 设置难度文本
		resource.DifficultyText = getDifficultyText(resource.Level)
		
		// 设置分类名称
		resource.CategoryName = getCategoryName(resource.Category)
		
		// 格式化文件大小
		resource.FileSize = formatFileSize(resource.TotalSize)
		
		resources = append(resources, resource)
	}

	return resources, total, nil
}

// 获取最新上传资料
func GetLatestResources(limit int) ([]LearningResource, error) {
	query := `
		SELECT r.id, r.title, r.cover_image, r.type, r.category, r.rating, r.download_count, r.created_at
		FROM learning_resources r
		ORDER BY r.created_at DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []LearningResource
	for rows.Next() {
		resource := LearningResource{}
		err := rows.Scan(&resource.ID, &resource.Title, &resource.CoverImage, &resource.Type, 
			&resource.Category, &resource.Rating, &resource.DownloadCount, &resource.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		resource.TypeIcon = getResourceTypeIcon(resource.Type)
		resource.CategoryName = getCategoryName(resource.Category)
		
		resources = append(resources, resource)
	}

	return resources, nil
}

// 获取最高评分资料
func GetTopRatedResources(limit int) ([]LearningResource, error) {
	query := `
		SELECT r.id, r.title, r.cover_image, r.category, r.rating, r.download_count
		FROM learning_resources r
		WHERE r.rating > 0
		ORDER BY r.rating DESC, r.download_count DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []LearningResource
	for rows.Next() {
		resource := LearningResource{}
		err := rows.Scan(&resource.ID, &resource.Title, &resource.CoverImage, &resource.Category, 
			&resource.Rating, &resource.DownloadCount)
		if err != nil {
			return nil, err
		}
		
		resource.CategoryName = getCategoryName(resource.Category)
		
		resources = append(resources, resource)
	}

	return resources, nil
}

// 根据ID获取学习资料
func GetLearningResourceByID(resourceID string) (*LearningResource, error) {
	resource := &LearningResource{}
	err := DB.QueryRow(`
		SELECT r.id, r.title, r.description, r.type, r.level, r.category, r.user_id, 
			   u.username, u.avatar, r.cover_image, r.file_paths, r.total_size, r.tags, 
			   r.rating, r.download_count, r.comment_count, r.download_url, r.created_at, r.updated_at
		FROM learning_resources r
		JOIN users u ON r.user_id = u.id
		WHERE r.id = ?
	`, resourceID).Scan(&resource.ID, &resource.Title, &resource.Description, &resource.Type, 
		&resource.Level, &resource.Category, &resource.UserID, &resource.UploaderName, 
		&resource.UploaderAvatar, &resource.CoverImage, &resource.FilePaths, &resource.TotalSize, 
		&resource.Tags, &resource.Rating, &resource.DownloadCount, &resource.CommentCount, 
		&resource.DownloadURL, &resource.CreatedAt, &resource.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	// 处理标签数组
	if resource.Tags != "" {
		resource.TagsArray = strings.Split(resource.Tags, ",")
		for i, tag := range resource.TagsArray {
			resource.TagsArray[i] = strings.TrimSpace(tag)
		}
	}
	
	// 设置类型图标
	resource.TypeIcon = getResourceTypeIcon(resource.Type)
	
	// 设置难度文本
	resource.DifficultyText = getDifficultyText(resource.Level)
	
	// 设置分类名称
	resource.CategoryName = getCategoryName(resource.Category)
	
	// 格式化文件大小
	resource.FileSize = formatFileSize(resource.TotalSize)
	
	return resource, nil
}

// 根据slug获取分类
func GetCategoryBySlug(slug string) (*ResourceCategory, error) {
	category := &ResourceCategory{}
	err := DB.QueryRow(`
		SELECT id, name, slug, description, icon
		FROM resource_categories
		WHERE slug = ?
	`, slug).Scan(&category.ID, &category.Name, &category.Slug, &category.Description, &category.Icon)
	
	if err != nil {
		return nil, err
	}
	
	// 获取分类下的资料数量
	err = DB.QueryRow("SELECT COUNT(*) FROM learning_resources WHERE category = ?", slug).Scan(&category.Count)
	if err != nil {
		return nil, err
	}
	
	return category, nil
}

// 增加下载次数
func IncrementResourceDownloads(resourceID string, userID int) error {
	// 记录下载历史
	_, err := DB.Exec("INSERT INTO resource_downloads (resource_id, user_id) VALUES (?, ?)", resourceID, userID)
	if err != nil {
		return err
	}
	
	// 更新下载计数
	_, err = DB.Exec("UPDATE learning_resources SET download_count = download_count + 1 WHERE id = ?", resourceID)
	return err
}

// 评分学习资料
func RateLearningResource(resourceID string, userID int, rating int) error {
	// 检查是否已评分
	var exists int
	err := DB.QueryRow("SELECT 1 FROM resource_ratings WHERE resource_id = ? AND user_id = ?", resourceID, userID).Scan(&exists)
	if err == nil {
		// 已评分，更新评分
		_, err = DB.Exec("UPDATE resource_ratings SET rating = ?, updated_at = NOW() WHERE resource_id = ? AND user_id = ?", rating, resourceID, userID)
	} else if err == sql.ErrNoRows {
		// 未评分，添加评分
		_, err = DB.Exec("INSERT INTO resource_ratings (resource_id, user_id, rating) VALUES (?, ?, ?)", resourceID, userID, rating)
	}
	
	if err != nil {
		return err
	}
	
	// 更新资料的平均评分
	_, err = DB.Exec(`
		UPDATE learning_resources 
		SET rating = (SELECT AVG(rating) FROM resource_ratings WHERE resource_id = ?)
		WHERE id = ?
	`, resourceID, resourceID)
	
	return err
}

// 评论学习资料
func CommentLearningResource(resourceID string, userID int, content string) error {
	// 添加评论
	_, err := DB.Exec("INSERT INTO resource_comments (resource_id, user_id, content) VALUES (?, ?, ?)", resourceID, userID, content)
	if err != nil {
		return err
	}
	
	// 更新评论计数
	_, err = DB.Exec("UPDATE learning_resources SET comment_count = comment_count + 1 WHERE id = ?", resourceID)
	return err
}

// 上传文件
func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// 创建上传目录
	uploadDir := "uploads/" + folder
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	
	// 生成文件名
	filename := filepath.Base(file.Filename)
	filepath := filepath.Join(uploadDir, filename)
	
	// 保存文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	// 复制文件内容
	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", err
	}
	
	return filepath, nil
}

// 工具函数
func getResourceTypeIcon(resourceType string) string {
	iconMap := map[string]string{
		"ebook": "fas fa-book",
		"video": "fas fa-video",
		"slides": "fas fa-presentation",
		"dataset": "fas fa-database",
		"code": "fas fa-code",
		"paper": "fas fa-file-alt",
	}
	
	if icon, exists := iconMap[resourceType]; exists {
		return icon
	}
	return "fas fa-file"
}

func getDifficultyText(level string) string {
	textMap := map[string]string{
		"beginner": "入门",
		"intermediate": "进阶",
		"advanced": "高级",
	}
	
	if text, exists := textMap[level]; exists {
		return text
	}
	return level
}

func getCategoryName(category string) string {
	nameMap := map[string]string{
		"machine-learning": "机器学习",
		"deep-learning": "深度学习",
		"robotics": "机器人学",
		"computer-vision": "计算机视觉",
		"nlp": "自然语言处理",
		"reinforcement-learning": "强化学习",
		"data-science": "数据科学",
		"ai-ethics": "AI伦理",
	}
	
	if name, exists := nameMap[category]; exists {
		return name
	}
	return category
}

func formatFileSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string(rune(float64(bytes)/float64(div))) + " " + string("KMGTPE"[exp]) + "B"
} 