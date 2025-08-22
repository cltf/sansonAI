package models

import (
	"database/sql"
)

// 创建回答
func CreateAnswer(questionID, userID int, content string) (int, error) {
	result, err := DB.Exec("INSERT INTO answers (question_id, user_id, content) VALUES (?, ?, ?)",
		questionID, userID, content)
	if err != nil {
		return 0, err
	}

	answerID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 更新问题回答数量
	_, err = DB.Exec("UPDATE questions SET answer_count = answer_count + 1 WHERE id = ?", questionID)
	if err != nil {
		return 0, err
	}

	// 给回答用户加积分
	err = UpdateUserPoints(userID, 3)
	if err != nil {
		return 0, err
	}

	return int(answerID), nil
}

// 根据问题ID获取回答
func GetAnswersByQuestionID(questionID int) ([]*Answer, error) {
	query := `
		SELECT a.id, a.question_id, a.user_id, u.username, u.avatar, 
			   a.content, a.like_count, a.is_accepted, a.created_at
		FROM answers a
		JOIN users u ON a.user_id = u.id
		WHERE a.question_id = ?
		ORDER BY a.is_accepted DESC, a.like_count DESC, a.created_at ASC
	`
	
	rows, err := DB.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*Answer
	for rows.Next() {
		answer := &Answer{}
		err := rows.Scan(
			&answer.ID, &answer.QuestionID, &answer.UserID, &answer.Username, &answer.UserAvatar,
			&answer.Content, &answer.LikeCount, &answer.IsAccepted, &answer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

// 采纳回答
func AcceptAnswer(answerID, userID int) error {
	// 获取回答信息
	var questionID int
	var questionUserID int
	var answerUserID int
	var reward int
	var isSolved bool
	
	err := DB.QueryRow(`
		SELECT q.id, q.user_id, a.user_id, q.reward, q.is_solved 
		FROM questions q 
		JOIN answers a ON q.id = a.question_id 
		WHERE a.id = ?
	`, answerID).Scan(&questionID, &questionUserID, &answerUserID, &reward, &isSolved)
	
	if err != nil {
		return err
	}
	
	// 检查权限（只有提问者可以采纳回答）
	if userID != questionUserID {
		return sql.ErrNoRows
	}
	
	// 检查问题是否已解决
	if isSolved {
		return sql.ErrNoRows
	}
	
	// 开始事务
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 标记回答为采纳
	_, err = tx.Exec("UPDATE answers SET is_accepted = 1 WHERE id = ?", answerID)
	if err != nil {
		return err
	}
	
	// 标记问题为已解决
	_, err = tx.Exec("UPDATE questions SET is_solved = 1 WHERE id = ?", questionID)
	if err != nil {
		return err
	}
	
	// 分配悬赏积分
	if reward > 0 {
		// 给回答者加积分
		_, err = tx.Exec("UPDATE users SET points = points + ? WHERE id = ?", reward, answerUserID)
		if err != nil {
			return err
		}
	}
	
	// 提交事务
	return tx.Commit()
}

// 点赞回答
func LikeAnswer(answerID, userID int) error {
	// 检查是否已点赞
	var exists int
	err := DB.QueryRow("SELECT 1 FROM answer_likes WHERE answer_id = ? AND user_id = ?", answerID, userID).Scan(&exists)
	if err == nil {
		// 已点赞，取消点赞
		_, err = DB.Exec("DELETE FROM answer_likes WHERE answer_id = ? AND user_id = ?", answerID, userID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("UPDATE answers SET like_count = like_count - 1 WHERE id = ?", answerID)
		return err
	} else if err == sql.ErrNoRows {
		// 未点赞，添加点赞
		_, err = DB.Exec("INSERT INTO answer_likes (answer_id, user_id) VALUES (?, ?)", answerID, userID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("UPDATE answers SET like_count = like_count + 1 WHERE id = ?", answerID)
		return err
	}
	
	return err
}

// 根据ID获取回答
func GetAnswerByID(answerID int) (*Answer, error) {
	answer := &Answer{}
	err := DB.QueryRow(`
		SELECT a.id, a.question_id, a.user_id, u.username, u.avatar, 
			   a.content, a.like_count, a.is_accepted, a.created_at
		FROM answers a
		JOIN users u ON a.user_id = u.id
		WHERE a.id = ?
	`, answerID).Scan(
		&answer.ID, &answer.QuestionID, &answer.UserID, &answer.Username, &answer.UserAvatar,
		&answer.Content, &answer.LikeCount, &answer.IsAccepted, &answer.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return answer, nil
} 