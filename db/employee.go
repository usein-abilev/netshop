package db

type EmployeeEntity struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EmployeeEntityStore struct {
	db *DatabaseConnection
}

func NewEmployeeEntityStore(database *DatabaseConnection) *EmployeeEntityStore {
	return &EmployeeEntityStore{
		db: database,
	}
}

func (p *EmployeeEntityStore) GetById(id int64) (EmployeeEntity, error) {
	row := p.db.Connection.QueryRow(p.db.Context, `select "id", "username", "password" from "employees" where id = $1`, id)
	var employee EmployeeEntity
	err := row.Scan(&employee.Id, &employee.Username, &employee.Password)
	if err != nil {
		return EmployeeEntity{}, err
	}
	return employee, nil
}

func (e *EmployeeEntityStore) GetByUsername(username string) (EmployeeEntity, error) {
	row := e.db.Connection.QueryRow(e.db.Context, `select "id", "username", "password" from "employees" where username = $1`, username)
	var employee EmployeeEntity
	err := row.Scan(&employee.Id, &employee.Username, &employee.Password)
	if err != nil {
		return EmployeeEntity{}, err
	}
	return employee, nil
}
