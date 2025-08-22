package models

// 获取所有分类
func GetCategories() ([]*Category, error) {
	query := `
		SELECT id, name, description, post_count
		FROM categories
		ORDER BY post_count DESC, name ASC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.PostCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// 根据ID获取分类
func GetCategoryByID(id int) (*Category, error) {
	category := &Category{}
	err := DB.QueryRow("SELECT id, name, description, post_count FROM categories WHERE id = ?", id).
		Scan(&category.ID, &category.Name, &category.Description, &category.PostCount)
	
	if err != nil {
		return nil, err
	}
	return category, nil
} 