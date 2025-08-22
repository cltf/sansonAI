package models

// 创建回复
func CreateReply(postID, userID int, content string) (int, error) {
	result, err := DB.Exec("INSERT INTO replies (post_id, user_id, content) VALUES (?, ?, ?)",
		postID, userID, content)
	if err != nil {
		return 0, err
	}

	replyID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 更新帖子回复数量
	_, err = DB.Exec("UPDATE posts SET reply_count = reply_count + 1 WHERE id = ?", postID)
	if err != nil {
		return 0, err
	}

	// 给回复用户加积分
	err = UpdateUserPoints(userID, 2)
	if err != nil {
		return 0, err
	}

	return int(replyID), nil
}

// 根据帖子ID获取回复
func GetRepliesByPostID(postID int) ([]*Reply, error) {
	query := `
		SELECT r.id, r.post_id, r.user_id, u.username, u.avatar, r.content, r.created_at
		FROM replies r
		JOIN users u ON r.user_id = u.id
		WHERE r.post_id = ?
		ORDER BY r.created_at ASC
	`
	
	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*Reply
	for rows.Next() {
		reply := &Reply{}
		err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.UserID, &reply.Username, &reply.UserAvatar,
			&reply.Content, &reply.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, nil
} 