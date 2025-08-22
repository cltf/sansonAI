package models

import (
	"time"
	"aiforum/utils"
)

// 用户动态结构
type UserActivity struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Type      string    `json:"type"` // question, answer, share, comment, resource
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TargetID  int       `json:"target_id"`
	Views     int       `json:"views"`
	Comments  int       `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

// 用户提问结构
type UserQuestion struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Status      string    `json:"status"` // open, answered, closed
	AnswerCount int       `json:"answer_count"`
	Views       int       `json:"views"`
	CreatedAt   time.Time `json:"created_at"`
}

// 用户回答结构
type UserAnswer struct {
	ID            int       `json:"id"`
	QuestionID    int       `json:"question_id"`
	QuestionTitle string    `json:"question_title"`
	Content       string    `json:"content"`
	IsAccepted    bool      `json:"is_accepted"`
	Likes         int       `json:"likes"`
	Comments      int       `json:"comments"`
	CreatedAt     time.Time `json:"created_at"`
}

// 用户分享结构
type UserShare struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Views     int       `json:"views"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

// 用户资料结构
type UserResource struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	FileSize    int64     `json:"file_size"`
	Downloads   int       `json:"downloads"`
	Views       int       `json:"views"`
	CreatedAt   time.Time `json:"created_at"`
}

// 用户收藏结构
type UserFavorite struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` // question, share, resource
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// 用户关注结构
type UserFollowing struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Followers  int    `json:"followers"`
	Questions  int    `json:"questions"`
	Answers    int    `json:"answers"`
	CreatedAt  time.Time `json:"created_at"`
}

// 用户粉丝结构
type UserFollower struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Bio        string `json:"bio"`
	Followers  int    `json:"followers"`
	Questions  int    `json:"questions"`
	Answers    int    `json:"answers"`
	CreatedAt  time.Time `json:"created_at"`
}

// 用户消息结构
type UserMessage struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` // system, question, follow, comment
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Sender    string    `json:"sender"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

// 获取用户动态
func GetUserActivity(userID int, filter string) ([]UserActivity, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			'question' as type,
			q.id,
			q.user_id,
			q.title,
			LEFT(q.content, 200) as content,
			q.id as target_id,
			q.views,
			(SELECT COUNT(*) FROM answers WHERE question_id = q.id) as comments,
			q.created_at
		FROM questions q
		WHERE q.user_id = ?
		
		UNION ALL
		
		SELECT 
			'answer' as type,
			a.id,
			a.user_id,
			CONCAT('回答了：', q.title) as title,
			LEFT(a.content, 200) as content,
			a.question_id as target_id,
			q.views,
			(SELECT COUNT(*) FROM comments WHERE answer_id = a.id) as comments,
			a.created_at
		FROM answers a
		JOIN questions q ON a.question_id = q.id
		WHERE a.user_id = ?
		
		ORDER BY created_at DESC
		LIMIT 50
	`

	if filter != "all" && filter != "" {
		query = `
			SELECT * FROM (
				` + baseQuery + `
			) as activities
			WHERE type = ?
			ORDER BY created_at DESC
		`
		args = []interface{}{userID, userID, filter}
	} else {
		query = baseQuery
		args = []interface{}{userID, userID}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []UserActivity
	for rows.Next() {
		var activity UserActivity
		err := rows.Scan(
			&activity.Type,
			&activity.ID,
			&activity.UserID,
			&activity.Title,
			&activity.Content,
			&activity.TargetID,
			&activity.Views,
			&activity.Comments,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// 获取用户提问
func GetUserQuestions(userID int) ([]UserQuestion, error) {
	query := `
		SELECT 
			id,
			title,
			LEFT(content, 200) as content,
			status,
			(SELECT COUNT(*) FROM answers WHERE question_id = questions.id) as answer_count,
			views,
			created_at
		FROM questions 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []UserQuestion
	for rows.Next() {
		var question UserQuestion
		err := rows.Scan(
			&question.ID,
			&question.Title,
			&question.Content,
			&question.Status,
			&question.AnswerCount,
			&question.Views,
			&question.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// 获取用户回答
func GetUserAnswers(userID int) ([]UserAnswer, error) {
	query := `
		SELECT 
			a.id,
			a.question_id,
			q.title as question_title,
			LEFT(a.content, 200) as content,
			a.is_accepted,
			a.likes,
			(SELECT COUNT(*) FROM comments WHERE answer_id = a.id) as comments,
			a.created_at
		FROM answers a
		JOIN questions q ON a.question_id = q.id
		WHERE a.user_id = ?
		ORDER BY a.created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []UserAnswer
	for rows.Next() {
		var answer UserAnswer
		err := rows.Scan(
			&answer.ID,
			&answer.QuestionID,
			&answer.QuestionTitle,
			&answer.Content,
			&answer.IsAccepted,
			&answer.Likes,
			&answer.Comments,
			&answer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

// 获取用户分享
func GetUserShares(userID int) ([]UserShare, error) {
	query := `
		SELECT 
			id,
			title,
			LEFT(content, 200) as content,
			category,
			views,
			likes,
			(SELECT COUNT(*) FROM comments WHERE tech_article_id = tech_articles.id) as comments,
			created_at
		FROM tech_articles 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []UserShare
	for rows.Next() {
		var share UserShare
		err := rows.Scan(
			&share.ID,
			&share.Title,
			&share.Content,
			&share.Category,
			&share.Views,
			&share.Likes,
			&share.Comments,
			&share.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, share)
	}

	return shares, nil
}

// 获取用户资料
func GetUserResources(userID int) ([]UserResource, error) {
	query := `
		SELECT 
			id,
			title,
			description,
			type,
			file_size,
			downloads,
			views,
			created_at
		FROM learning_resources 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []UserResource
	for rows.Next() {
		var resource UserResource
		err := rows.Scan(
			&resource.ID,
			&resource.Title,
			&resource.Description,
			&resource.Type,
			&resource.FileSize,
			&resource.Downloads,
			&resource.Views,
			&resource.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// 获取用户收藏
func GetUserFavorites(userID int) ([]UserFavorite, error) {
	query := `
		SELECT 
			f.id,
			f.type,
			CASE 
				WHEN f.type = 'question' THEN q.title
				WHEN f.type = 'share' THEN t.title
				WHEN f.type = 'resource' THEN lr.title
			END as title,
			CASE 
				WHEN f.type = 'question' THEN LEFT(q.content, 100)
				WHEN f.type = 'share' THEN LEFT(t.content, 100)
				WHEN f.type = 'resource' THEN lr.description
			END as content,
			u.username as author,
			f.created_at
		FROM favorites f
		LEFT JOIN questions q ON f.type = 'question' AND f.target_id = q.id
		LEFT JOIN tech_articles t ON f.type = 'share' AND f.target_id = t.id
		LEFT JOIN learning_resources lr ON f.type = 'resource' AND f.target_id = lr.id
		LEFT JOIN users u ON 
			(f.type = 'question' AND q.user_id = u.id) OR
			(f.type = 'share' AND t.user_id = u.id) OR
			(f.type = 'resource' AND lr.user_id = u.id)
		WHERE f.user_id = ?
		ORDER BY f.created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []UserFavorite
	for rows.Next() {
		var favorite UserFavorite
		err := rows.Scan(
			&favorite.ID,
			&favorite.Type,
			&favorite.Title,
			&favorite.Content,
			&favorite.Author,
			&favorite.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, favorite)
	}

	return favorites, nil
}

// 获取用户关注
func GetUserFollowing(userID int) ([]UserFollowing, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.avatar,
			u.bio,
			(SELECT COUNT(*) FROM follows WHERE followed_id = u.id) as followers,
			(SELECT COUNT(*) FROM questions WHERE user_id = u.id) as questions,
			(SELECT COUNT(*) FROM answers WHERE user_id = u.id) as answers,
			f.created_at
		FROM follows f
		JOIN users u ON f.followed_id = u.id
		WHERE f.follower_id = ?
		ORDER BY f.created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []UserFollowing
	for rows.Next() {
		var user UserFollowing
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Avatar,
			&user.Bio,
			&user.Followers,
			&user.Questions,
			&user.Answers,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		following = append(following, user)
	}

	return following, nil
}

// 获取用户粉丝
func GetUserFollowers(userID int) ([]UserFollower, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.avatar,
			u.bio,
			(SELECT COUNT(*) FROM follows WHERE followed_id = u.id) as followers,
			(SELECT COUNT(*) FROM questions WHERE user_id = u.id) as questions,
			(SELECT COUNT(*) FROM answers WHERE user_id = u.id) as answers,
			f.created_at
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.followed_id = ?
		ORDER BY f.created_at DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []UserFollower
	for rows.Next() {
		var user UserFollower
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Avatar,
			&user.Bio,
			&user.Followers,
			&user.Questions,
			&user.Answers,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}

	return followers, nil
}

// 获取用户消息
func GetUserMessages(userID int) ([]UserMessage, error) {
	query := `
		SELECT 
			id,
			type,
			title,
			content,
			sender,
			is_read,
			created_at
		FROM messages 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []UserMessage
	for rows.Next() {
		var message UserMessage
		err := rows.Scan(
			&message.ID,
			&message.Type,
			&message.Title,
			&message.Content,
			&message.Sender,
			&message.IsRead,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// 更新用户头像
func UpdateUserAvatar(userID int, avatarURL string) error {
	_, err := DB.Exec("UPDATE users SET avatar = ?, updated_at = ? WHERE id = ?",
		avatarURL, time.Now(), userID)
	return err
}

// 更新用户个人资料
func UpdateUserProfile(userID int, username, email, bio, phone, website string, profilePublic, showEmail, showPhone bool) error {
	_, err := DB.Exec(`
		UPDATE users SET 
			username = ?, 
			email = ?, 
			bio = ?, 
			phone = ?, 
			website = ?, 
			profile_public = ?, 
			show_email = ?, 
			show_phone = ?, 
			updated_at = ? 
		WHERE id = ?
	`, username, email, bio, phone, website, profilePublic, showEmail, showPhone, time.Now(), userID)
	return err
}

// 更新用户密码
func UpdateUserPassword(userID int, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = DB.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?",
		hashedPassword, time.Now(), userID)
	return err
}

// 检查密码
func CheckPassword(password, hashedPassword string) bool {
	return utils.CheckPassword(password, hashedPassword)
}

// 更新通知设置
func UpdateNotificationSettings(userID int, emailNotifications, browserNotifications, questionNotifications, followNotifications bool) error {
	_, err := DB.Exec(`
		UPDATE users SET 
			email_notifications = ?, 
			browser_notifications = ?, 
			question_notifications = ?, 
			follow_notifications = ?, 
			updated_at = ? 
		WHERE id = ?
	`, emailNotifications, browserNotifications, questionNotifications, followNotifications, time.Now(), userID)
	return err
}

// 标记所有消息为已读
func MarkAllMessagesRead(userID int) error {
	_, err := DB.Exec("UPDATE messages SET is_read = 1 WHERE user_id = ?", userID)
	return err
}

// 标记单条消息为已读
func MarkMessageRead(userID, messageID int) error {
	_, err := DB.Exec("UPDATE messages SET is_read = 1 WHERE id = ? AND user_id = ?", messageID, userID)
	return err
}

// 删除消息
func DeleteMessage(userID, messageID int) error {
	_, err := DB.Exec("DELETE FROM messages WHERE id = ? AND user_id = ?", messageID, userID)
	return err
}

// 关注用户
func FollowUser(followerID, followedID int) error {
	// 检查是否已经关注
	var exists int
	err := DB.QueryRow("SELECT 1 FROM follows WHERE follower_id = ? AND followed_id = ?", followerID, followedID).Scan(&exists)
	if err == nil {
		return nil // 已经关注了
	}

	_, err = DB.Exec("INSERT INTO follows (follower_id, followed_id, created_at) VALUES (?, ?, ?)",
		followerID, followedID, time.Now())
	return err
}

// 取消关注用户
func UnfollowUser(followerID, followedID int) error {
	_, err := DB.Exec("DELETE FROM follows WHERE follower_id = ? AND followed_id = ?", followerID, followedID)
	return err
}

// 取消收藏
func RemoveFavorite(userID, favoriteID int) error {
	_, err := DB.Exec("DELETE FROM favorites WHERE id = ? AND user_id = ?", favoriteID, userID)
	return err
}

// 删除用户提问
func DeleteUserQuestion(userID, questionID int) error {
	_, err := DB.Exec("DELETE FROM questions WHERE id = ? AND user_id = ?", questionID, userID)
	return err
}

// 删除用户回答
func DeleteUserAnswer(userID, answerID int) error {
	_, err := DB.Exec("DELETE FROM answers WHERE id = ? AND user_id = ?", answerID, userID)
	return err
}

// 删除用户分享
func DeleteUserShare(userID, shareID int) error {
	_, err := DB.Exec("DELETE FROM tech_articles WHERE id = ? AND user_id = ?", shareID, userID)
	return err
}

// 删除用户资料
func DeleteUserResource(userID, resourceID int) error {
	_, err := DB.Exec("DELETE FROM learning_resources WHERE id = ? AND user_id = ?", resourceID, userID)
	return err
} 