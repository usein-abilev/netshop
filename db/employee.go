package db

type Employee struct {
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

func (p *EmployeeEntityStore) GetEmployeeById(id int64) (Employee, error) {
	row := p.db.Connection.QueryRow(p.db.Context, `select "id", "username", "password" from "employees" where id = $1`, id)
	var employee Employee
	err := row.Scan(&employee.Id, &employee.Username, &employee.Password)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func (e *EmployeeEntityStore) GetEmployeeByUsername(username string) (Employee, error) {
	row := e.db.Connection.QueryRow(e.db.Context, `select "id", "username", "password" from "employees" where username = $1`, username)
	var employee Employee
	err := row.Scan(&employee.Id, &employee.Username, &employee.Password)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}
