package db

import (
	"context"
	"errors"
	"log"
)

type Product struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Id            int64  `json:"id"`
	Price         int32  `json:"price"`
}

type ProductEntityStore struct {
	db *DatabaseConnection
}

func NewProductEntityStore(database *DatabaseConnection) *ProductEntityStore {
	return &ProductEntityStore{
		db: database,
	}
}

func (p *ProductEntityStore) GetProductById(id int64) (Product, error) {
	row := p.db.Connection.QueryRow(p.db.Context, `select "id", "name", "description", "price" from "products" where id = $1`, id)
	var product Product
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (p *ProductEntityStore) GetProducts() ([]Product, error) {
	query := `select "id", "name", "description", "price" from "products"`
	rows, err := p.db.Connection.Query(p.db.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err != nil {
			log.Println("Error while scanning rows", err)
			return nil, errors.New("Error while scanning rows")
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductEntityStore) CreateProduct(ctx context.Context, product *Product) error {
	_, err := p.db.Connection.Exec(ctx, `insert into "products" (name, description, price) values ($1, $2, $3)`, product.Name, product.Description, "0")
	if err != nil {
		return err
	}
	return nil
}
