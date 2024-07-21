package db

// CategoryEntity represents a category of products in the database
type CategoryEntity struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type CategoryEntityStore struct {
	db *DatabaseConnection
}

func NewCategoryEntityStore(database *DatabaseConnection) *CategoryEntityStore {
	return &CategoryEntityStore{
		db: database,
	}
}

func (c *CategoryEntityStore) CheckCategoryExists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "categories" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *CategoryEntityStore) GetCategoryById(id int64) (CategoryEntity, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select "id", "name" from "categories" where id = $1`, id)
	var category CategoryEntity
	err := row.Scan(&category.Id, &category.Name)
	if err != nil {
		return CategoryEntity{}, err
	}
	return category, nil
}

func (c *CategoryEntityStore) GetCategories() ([]CategoryEntity, error) {
	query := `select "id", "name" from "categories"`
	rows, err := c.db.Connection.Query(c.db.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]CategoryEntity, 0)
	for rows.Next() {
		var category CategoryEntity
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
