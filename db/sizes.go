package db

// SizeEntity represents a category of products in the database
type SizeEntity struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SizeEntityStore struct {
	db *DatabaseConnection
}

func NewSizeEntityStore(database *DatabaseConnection) *SizeEntityStore {
	return &SizeEntityStore{
		db: database,
	}
}

func (c *SizeEntityStore) Exists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "sizes" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *SizeEntityStore) GetById(id int64) (SizeEntity, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select "id", "name" from "sizes" where id = $1`, id)
	var category SizeEntity
	err := row.Scan(&category.Id, &category.Name)
	if err != nil {
		return SizeEntity{}, err
	}
	return category, nil
}

func (c *SizeEntityStore) GetEntities() ([]SizeEntity, error) {
	query := `select "id", "name" from "sizes"`
	rows, err := c.db.Connection.Query(c.db.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sizes := make([]SizeEntity, 0)
	for rows.Next() {
		var category SizeEntity
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}
		sizes = append(sizes, category)
	}

	return sizes, nil
}
