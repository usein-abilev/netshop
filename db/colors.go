package db

type ColorEntity struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ColorEntityStore struct {
	db *DatabaseConnection
}

func NewColorEntityStore(database *DatabaseConnection) *ColorEntityStore {
	return &ColorEntityStore{
		db: database,
	}
}

func (c *ColorEntityStore) Exists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "colors" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *ColorEntityStore) GetById(id int64) (ColorEntity, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select "id", "name" from "colors" where id = $1`, id)
	var category ColorEntity
	err := row.Scan(&category.Id, &category.Name)
	if err != nil {
		return ColorEntity{}, err
	}
	return category, nil
}

func (c *ColorEntityStore) GetEntities() ([]ColorEntity, error) {
	query := `select "id", "name" from "colors"`
	rows, err := c.db.Connection.Query(c.db.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	colors := make([]ColorEntity, 0)
	for rows.Next() {
		var category ColorEntity
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, err
		}
		colors = append(colors, category)
	}

	return colors, nil
}
