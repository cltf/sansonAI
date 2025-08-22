package models

// 创建帖子
func CreatePost(title, content string, categoryID, userID int, tags string) (int, error) {
	result, err := DB.Exec("INSERT INTO posts (title, content, category_id, user_id, tags) VALUES (?, ?, ?, ?, ?)",
		title, content, categoryID, userID, tags)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 更新分类帖子数量
	_, err = DB.Exec("UPDATE categories SET post_count = post_count + 1 WHERE id = ?", categoryID)
	if err != nil {
		return 0, err
	}

	// 给发帖用户加积分
	err = UpdateUserPoints(userID, 10)
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

// 获取帖子列表
func GetPosts(page, limit int, categoryID int) ([]*Post, error) {
	offset := (page - 1) * limit
	
	var query string
	var args []interface{}
	
	if categoryID > 0 {
		query = `
			SELECT p.id, p.title, p.content, p.category_id, p.user_id, 
				   u.username, u.avatar, p.view_count, p.reply_count, p.like_count, 
				   p.tags, p.created_at, p.updated_at
			FROM posts p
			JOIN users u ON p.user_id = u.id
			WHERE p.category_id = ?
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{categoryID, limit, offset}
	} else {
		query = `
			SELECT p.id, p.title, p.content, p.category_id, p.user_id, 
				   u.username, u.avatar, p.view_count, p.reply_count, p.like_count, 
				   p.tags, p.created_at, p.updated_at
			FROM posts p
			JOIN users u ON p.user_id = u.id
			ORDER BY p.created_at DESC
			LIMIT ? OFFSET ?
		`
		args = []interface{}{limit, offset}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.UserID,
			&post.Username, &post.UserAvatar, &post.ViewCount, &post.ReplyCount, &post.LikeCount,
			&post.Tags, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// 根据ID获取帖子
func GetPostByID(id int) (*Post, error) {
	post := &Post{}
	err := DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.category_id, p.user_id, 
			   u.username, u.avatar, p.view_count, p.reply_count, p.like_count, 
			   p.tags, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.UserID,
		&post.Username, &post.UserAvatar, &post.ViewCount, &post.ReplyCount, &post.LikeCount,
		&post.Tags, &post.CreatedAt, &post.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	return post, nil
}

// 增加帖子浏览量
func IncrementPostViewCount(id int) error {
	_, err := DB.Exec("UPDATE posts SET view_count = view_count + 1 WHERE id = ?", id)
	return err
}

// 搜索帖子
func SearchPosts(keyword string, page, limit int) ([]*Post, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT p.id, p.title, p.content, p.category_id, p.user_id, 
			   u.username, u.avatar, p.view_count, p.reply_count, p.like_count, 
			   p.tags, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.title LIKE ? OR p.content LIKE ? OR p.tags LIKE ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	keyword = "%" + keyword + "%"
	rows, err := DB.Query(query, keyword, keyword, keyword, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.UserID,
			&post.Username, &post.UserAvatar, &post.ViewCount, &post.ReplyCount, &post.LikeCount,
			&post.Tags, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// 获取热门帖子
func GetHotPosts(limit int) ([]*Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.category_id, p.user_id, 
			   u.username, u.avatar, p.view_count, p.reply_count, p.like_count, 
			   p.tags, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY (p.view_count + p.reply_count * 2 + p.like_count * 3) DESC
		LIMIT ?
	`
	
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.CategoryID, &post.UserID,
			&post.Username, &post.UserAvatar, &post.ViewCount, &post.ReplyCount, &post.LikeCount,
			&post.Tags, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
} 