package models

import (
	"database/sql"
	"strings"
	"time"
)

// 问题模型
type Question struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	CategoryID      int       `json:"category_id"`
	UserID          int       `json:"user_id"`
	Username        string    `json:"username"`
	UserAvatar      string    `json:"user_avatar"`
	UserLevel       int       `json:"user_level"`
	ViewCount       int       `json:"view_count"`
	AnswerCount     int       `json:"answer_count"`
	LikeCount       int       `json:"like_count"`
	Tags            string    `json:"tags"`
	Reward          int       `json:"reward"`
	IsSolved        bool      `json:"is_solved"`
	AcceptedAnswer  string    `json:"accepted_answer"`
	Summary         string    `json:"summary"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// 回答模型
type Answer struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	UserAvatar string    `json:"user_avatar"`
	Content    string    `json:"content"`
	LikeCount  int       `json:"like_count"`
	IsAccepted bool      `json:"is_accepted"`
	CreatedAt  time.Time `json:"created_at"`
}

// 创建问题
func CreateQuestion(title, content string, categoryID, userID int, tags string, reward int) (int, error) {
	// 生成问题摘要
	summary := generateSummary(content)
	
	result, err := DB.Exec(`
		INSERT INTO questions (title, content, category_id, user_id, tags, reward, summary) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, title, content, categoryID, userID, tags, reward, summary)
	
	if err != nil {
		return 0, err
	}

	questionID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 更新分类问题数量
	_, err = DB.Exec("UPDATE categories SET post_count = post_count + 1 WHERE id = ?", categoryID)
	if err != nil {
		return 0, err
	}

	// 给提问用户加积分
	err = UpdateUserPoints(userID, 5)
	if err != nil {
		return 0, err
	}

	return int(questionID), nil
}

// 获取问题列表
func GetQuestions(page, limit, categoryID int, sort string) ([]*Question, error) {
	offset := (page - 1) * limit
	
	var query string
	var args []interface{}
	
	baseQuery := `
		SELECT q.id, q.title, q.content, q.category_id, q.user_id, 
			   u.username, u.avatar, q.view_count, q.answer_count, q.like_count, 
			   q.tags, q.reward, q.is_solved, q.summary, q.created_at, q.updated_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
	`
	
	whereClause := ""
	if categoryID > 0 {
		whereClause = "WHERE q.category_id = ?"
		args = append(args, categoryID)
	}
	
	orderClause := "ORDER BY q.created_at DESC"
	switch sort {
	case "hot":
		orderClause = "ORDER BY (q.view_count + q.answer_count * 2 + q.like_count * 3) DESC"
	case "reward":
		orderClause = "ORDER BY q.reward DESC"
	case "unsolved":
		orderClause = "ORDER BY q.is_solved ASC, q.created_at DESC"
	}
	
	query = baseQuery + " " + whereClause + " " + orderClause + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(
			&question.ID, &question.Title, &question.Content, &question.CategoryID, &question.UserID,
			&question.Username, &question.UserAvatar, &question.ViewCount, &question.AnswerCount, &question.LikeCount,
			&question.Tags, &question.Reward, &question.IsSolved, &question.Summary, &question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// 获取采纳的回答
		if question.IsSolved {
			question.AcceptedAnswer, _ = getAcceptedAnswer(question.ID)
		}
		
		questions = append(questions, question)
	}

	return questions, nil
}

// 根据ID获取问题
func GetQuestionByID(id int) (*Question, error) {
	question := &Question{}
	err := DB.QueryRow(`
		SELECT q.id, q.title, q.content, q.category_id, q.user_id, 
			   u.username, u.avatar, q.view_count, q.answer_count, q.like_count, 
			   q.tags, q.reward, q.is_solved, q.summary, q.created_at, q.updated_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.id = ?
	`, id).Scan(
		&question.ID, &question.Title, &question.Content, &question.CategoryID, &question.UserID,
		&question.Username, &question.UserAvatar, &question.ViewCount, &question.AnswerCount, &question.LikeCount,
		&question.Tags, &question.Reward, &question.IsSolved, &question.Summary, &question.CreatedAt, &question.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	// 获取采纳的回答
	if question.IsSolved {
		question.AcceptedAnswer, _ = getAcceptedAnswer(question.ID)
	}
	
	return question, nil
}

// 搜索问题
func SearchQuestions(keyword string, page, limit int) ([]*Question, error) {
	offset := (page - 1) * limit
	keyword = "%" + keyword + "%"
	
	query := `
		SELECT q.id, q.title, q.content, q.category_id, q.user_id, 
			   u.username, u.avatar, q.view_count, q.answer_count, q.like_count, 
			   q.tags, q.reward, q.is_solved, q.summary, q.created_at, q.updated_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.title LIKE ? OR q.content LIKE ? OR q.tags LIKE ?
		ORDER BY q.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := DB.Query(query, keyword, keyword, keyword, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(
			&question.ID, &question.Title, &question.Content, &question.CategoryID, &question.UserID,
			&question.Username, &question.UserAvatar, &question.ViewCount, &question.AnswerCount, &question.LikeCount,
			&question.Tags, &question.Reward, &question.IsSolved, &question.Summary, &question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if question.IsSolved {
			question.AcceptedAnswer, _ = getAcceptedAnswer(question.ID)
		}
		
		questions = append(questions, question)
	}

	return questions, nil
}

// 按标签获取问题
func GetQuestionsByTag(tag string, page, limit int) ([]*Question, error) {
	offset := (page - 1) * limit
	tag = "%" + tag + "%"
	
	query := `
		SELECT q.id, q.title, q.content, q.category_id, q.user_id, 
			   u.username, u.avatar, q.view_count, q.answer_count, q.like_count, 
			   q.tags, q.reward, q.is_solved, q.summary, q.created_at, q.updated_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.tags LIKE ?
		ORDER BY q.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := DB.Query(query, tag, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(
			&question.ID, &question.Title, &question.Content, &question.CategoryID, &question.UserID,
			&question.Username, &question.UserAvatar, &question.ViewCount, &question.AnswerCount, &question.LikeCount,
			&question.Tags, &question.Reward, &question.IsSolved, &question.Summary, &question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if question.IsSolved {
			question.AcceptedAnswer, _ = getAcceptedAnswer(question.ID)
		}
		
		questions = append(questions, question)
	}

	return questions, nil
}

// 高级搜索
func AdvancedSearchQuestions(conditions map[string]interface{}, page, limit int) ([]*Question, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT q.id, q.title, q.content, q.category_id, q.user_id, 
			   u.username, u.avatar, q.view_count, q.answer_count, q.like_count, 
			   q.tags, q.reward, q.is_solved, q.summary, q.created_at, q.updated_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
	`
	
	var whereConditions []string
	var args []interface{}
	
	if title, ok := conditions["title"].(string); ok && title != "" {
		whereConditions = append(whereConditions, "q.title LIKE ?")
		args = append(args, "%"+title+"%")
	}
	
	if content, ok := conditions["content"].(string); ok && content != "" {
		whereConditions = append(whereConditions, "q.content LIKE ?")
		args = append(args, "%"+content+"%")
	}
	
	if author, ok := conditions["author"].(string); ok && author != "" {
		whereConditions = append(whereConditions, "u.username LIKE ?")
		args = append(args, "%"+author+"%")
	}
	
	if status, ok := conditions["status"].(string); ok && status != "" {
		if status == "solved" {
			whereConditions = append(whereConditions, "q.is_solved = 1")
		} else if status == "unsolved" {
			whereConditions = append(whereConditions, "q.is_solved = 0")
		}
	}
	
	if timeRange, ok := conditions["time_range"].(string); ok && timeRange != "" {
		var days int
		switch timeRange {
		case "1d":
			days = 1
		case "7d":
			days = 7
		case "30d":
			days = 30
		case "90d":
			days = 90
		}
		if days > 0 {
			whereConditions = append(whereConditions, "q.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)")
			args = append(args, days)
		}
	}
	
	if rewardRange, ok := conditions["reward_range"].(string); ok && rewardRange != "" {
		var minReward int
		switch rewardRange {
		case "10+":
			minReward = 10
		case "50+":
			minReward = 50
		case "100+":
			minReward = 100
		}
		if minReward > 0 {
			whereConditions = append(whereConditions, "q.reward >= ?")
			args = append(args, minReward)
		}
	}
	
	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	
	query += " ORDER BY q.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(
			&question.ID, &question.Title, &question.Content, &question.CategoryID, &question.UserID,
			&question.Username, &question.UserAvatar, &question.ViewCount, &question.AnswerCount, &question.LikeCount,
			&question.Tags, &question.Reward, &question.IsSolved, &question.Summary, &question.CreatedAt, &question.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if question.IsSolved {
			question.AcceptedAnswer, _ = getAcceptedAnswer(question.ID)
		}
		
		questions = append(questions, question)
	}

	return questions, nil
}

// 获取待解决问题
func GetPendingQuestions(limit int) ([]*Question, error) {
	query := `
		SELECT q.id, q.title, q.user_id, u.username, q.reward, q.created_at
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.is_solved = 0
		ORDER BY q.created_at DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(&question.ID, &question.Title, &question.UserID, &question.Username, &question.Reward, &question.CreatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// 获取高悬赏问题
func GetHighRewardQuestions(limit int) ([]*Question, error) {
	query := `
		SELECT q.id, q.title, q.user_id, u.username, q.reward
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.reward > 0
		ORDER BY q.reward DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(&question.ID, &question.Title, &question.UserID, &question.Username, &question.Reward)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// 获取相关问题
func GetRelatedQuestions(questionID, limit int) ([]*Question, error) {
	query := `
		SELECT q.id, q.title, q.user_id, u.username, q.view_count, q.answer_count
		FROM questions q
		JOIN users u ON q.user_id = u.id
		WHERE q.id != ? AND q.category_id = (SELECT category_id FROM questions WHERE id = ?)
		ORDER BY q.view_count DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, questionID, questionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		question := &Question{}
		err := rows.Scan(&question.ID, &question.Title, &question.UserID, &question.Username, &question.ViewCount, &question.AnswerCount)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// 增加问题浏览量
func IncrementQuestionViewCount(id int) error {
	_, err := DB.Exec("UPDATE questions SET view_count = view_count + 1 WHERE id = ?", id)
	return err
}

// 获取问题统计
func GetQuestionCount(categoryID int) (int, error) {
	var count int
	var err error
	
	if categoryID > 0 {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE category_id = ?", categoryID).Scan(&count)
	} else {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	}
	
	return count, err
}

func GetSolvedQuestionCount(categoryID int) (int, error) {
	var count int
	var err error
	
	if categoryID > 0 {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE category_id = ? AND is_solved = 1", categoryID).Scan(&count)
	} else {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE is_solved = 1").Scan(&count)
	}
	
	return count, err
}

func GetUnsolvedQuestionCount(categoryID int) (int, error) {
	var count int
	var err error
	
	if categoryID > 0 {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE category_id = ? AND is_solved = 0", categoryID).Scan(&count)
	} else {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE is_solved = 0").Scan(&count)
	}
	
	return count, err
}

func GetTodayQuestionCount(categoryID int) (int, error) {
	var count int
	var err error
	
	if categoryID > 0 {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE category_id = ? AND DATE(created_at) = CURDATE()", categoryID).Scan(&count)
	} else {
		err = DB.QueryRow("SELECT COUNT(*) FROM questions WHERE DATE(created_at) = CURDATE()").Scan(&count)
	}
	
	return count, err
}

// 生成问题摘要
func generateSummary(content string) string {
	if len(content) <= 200 {
		return content
	}
	return content[:200] + "..."
}

// 获取采纳的回答
func getAcceptedAnswer(questionID int) (string, error) {
	var content string
	err := DB.QueryRow(`
		SELECT content FROM answers 
		WHERE question_id = ? AND is_accepted = 1 
		LIMIT 1
	`, questionID).Scan(&content)
	
	if err == sql.ErrNoRows {
		return "", nil
	}
	
	if len(content) > 150 {
		content = content[:150] + "..."
	}
	
	return content, err
}

// 收藏/取消收藏问题
func ToggleQuestionFavorite(questionID, userID int) (bool, error) {
	// 检查是否已收藏
	var exists int
	err := DB.QueryRow("SELECT 1 FROM question_favorites WHERE question_id = ? AND user_id = ?", questionID, userID).Scan(&exists)
	if err == nil {
		// 已收藏，取消收藏
		_, err = DB.Exec("DELETE FROM question_favorites WHERE question_id = ? AND user_id = ?", questionID, userID)
		return false, err
	} else if err == sql.ErrNoRows {
		// 未收藏，添加收藏
		_, err = DB.Exec("INSERT INTO question_favorites (question_id, user_id) VALUES (?, ?)", questionID, userID)
		return true, err
	}
	
	return false, err
}

// 举报问题
func ReportQuestion(questionID, userID int, reason string) error {
	_, err := DB.Exec("INSERT INTO question_reports (question_id, user_id, reason) VALUES (?, ?, ?)", questionID, userID, reason)
	return err
} 