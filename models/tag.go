package models

// 获取所有标签
func GetTags() ([]*Tag, error) {
	query := `
		SELECT id, name
		FROM tags
		ORDER BY name ASC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := &Tag{}
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// 创建标签
func CreateTag(name string) error {
	_, err := DB.Exec("INSERT IGNORE INTO tags (name) VALUES (?)", name)
	return err
} 