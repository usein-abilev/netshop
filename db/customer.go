package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type CustomerEntity struct {
	Id         int64     `json:"id"`
	PersonId   int64     `json:"person_id"`
	Person     *Person   `json:"person"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsVerified bool      `json:"is_verified"`
}

type CustomerCreateUpdate struct {
	Person   PersonCreateUpdate `json:"person"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}

type CustomerEntityStore struct {
	db *DatabaseConnection
}

func NewCustomerEntityStore(database *DatabaseConnection) *CustomerEntityStore {
	return &CustomerEntityStore{
		db: database,
	}
}

func (c *CustomerEntityStore) Exists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "customers" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *CustomerEntityStore) GetById(id int64) (result *CustomerEntity, err error) {
	result = &CustomerEntity{
		Person: &Person{},
	}
	err = c.db.Connection.QueryRow(c.db.Context, `
		select 
			id, 
			person_id,
			person.first_name,
			person.last_name,
			person.email,
			person.phone,
			person.email_verified,
			person.metadata,
			username, 
			password,
			created_at,
			updated_at,
			is_verified
		from "customers"
		left join "person" on "person".id = "customers".person_id
		where id = $1
		`, id).Scan(
		&result.Id,
		&result.PersonId,
		&result.Person.FirstName,
		&result.Person.LastName,
		&result.Person.Email,
		&result.Person.Phone,
		&result.Person.EmailVerified,
		&result.Person.Metadata,
		&result.Username,
		&result.Password,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.IsVerified,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *CustomerEntityStore) GetByUsername(username string) (*CustomerEntity, error) {
	row := e.db.Connection.QueryRow(e.db.Context, `select "id", "username", "password" from "customers" where username = $1`, username)
	customer := &CustomerEntity{}
	err := row.Scan(&customer.Id, &customer.Username, &customer.Password)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *CustomerEntityStore) Create(ctx context.Context, options *CustomerCreateUpdate) (result *CustomerEntity, err error) {
	tx, err := c.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var alreadyExists bool
	err = tx.QueryRow(c.db.Context, `select exists(
			select 1 from "person"
			inner join "customers" on "customers".person_id = "person".id 
			where "person"."phone" = $1 or "person"."email" = $2
		)`, options.Person.Phone, options.Person.Email).Scan(&alreadyExists)
	if err != nil {
		return nil, err
	}
	if alreadyExists {
		return nil, fmt.Errorf("customer with the given phone number or email already exists")
	}

	personStore := NewPersonEntityStore(c.db)
	person, err := personStore.TxCreate(&tx, &options.Person)
	if err != nil {
		return result, fmt.Errorf("failed to create person: %w", err)
	}

	result = &CustomerEntity{
		PersonId: person.Id,
		Person:   person,
		Username: options.Username,
		Password: "",
	}

	err = tx.QueryRow(c.db.Context, `
		insert into "customers" (username, password, person_id, is_verified) values ($1, $2, $3, $4)
		returning id, created_at, updated_at`, options.Username, options.Password, person.Id, false).Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return result, tx.Commit(ctx)
}
