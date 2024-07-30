package db

import "github.com/jackc/pgx/v5"

type Person struct {
	Id            int64   `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Metadata      *string `json:"metadata"`
	EmailVerified bool    `json:"email_verified"`
}

type PersonCreateUpdate struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Phone     string  `json:"phone"`
	Email     string  `json:"email"`
	Metadata  *string `json:"metadata"`
}

type PersonEntityStore struct {
	db *DatabaseConnection
}

func NewPersonEntityStore(database *DatabaseConnection) *PersonEntityStore {
	return &PersonEntityStore{
		db: database,
	}
}

func (p *PersonEntityStore) TxCreate(tx *pgx.Tx, options *PersonCreateUpdate) (result *Person, err error) {
	conn := *tx
	result = &Person{
		FirstName:     options.FirstName,
		LastName:      options.LastName,
		Phone:         options.Phone,
		Email:         options.Email,
		Metadata:      options.Metadata,
		EmailVerified: false,
	}

	err = conn.QueryRow(p.db.Context, `
		insert into "person" (first_name, last_name, phone, email, metadata, email_verified) values ($1, $2, $3, $4, $5, $6) 
		returning id`, options.FirstName, options.LastName, options.Phone, options.Email, options.Metadata, false).Scan(&result.Id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
