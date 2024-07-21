package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type ProductEntity struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Id          int64  `json:"id"`
	Price       int32  `json:"price"`
}

type ProductVariantCreateUpdate struct {
	SizeId  int64 `json:"size_id"`
	ColorId int64 `json:"color_id"`
	Price   int32 `json:"price"`
	Stock   int32 `json:"stock"`
}

type ProductCreateUpdate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	CategoryId  int64                        `json:"category_id"`
	EmployeeId  int64                        `json:"employee_id"`
	BasePrice   int32                        `json:"base_price"`
	Variants    []ProductVariantCreateUpdate `json:"variants"`
}

type ProductEntityStore struct {
	db *DatabaseConnection
}

func NewProductEntityStore(database *DatabaseConnection) *ProductEntityStore {
	return &ProductEntityStore{
		db: database,
	}
}

func (p *ProductEntityStore) GetProductById(id int64) (ProductEntity, error) {
	row := p.db.Connection.QueryRow(p.db.Context, `select "id", "name", "description", "base_price" from "products" where id = $1`, id)
	var product ProductEntity
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err != nil {
		return ProductEntity{}, err
	}
	return product, nil
}

func (p *ProductEntityStore) GetProducts() ([]ProductEntity, error) {
	query := `select "id", "name", "description", "base_price" from "products"`
	rows, err := p.db.Connection.Query(p.db.Context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]ProductEntity, 0)
	for rows.Next() {
		var product ProductEntity
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
		if err != nil {
			log.Println("Error while scanning rows", err)
			return nil, errors.New("error while scanning rows")
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductEntityStore) CreateProduct(ctx context.Context, opts *ProductCreateUpdate) error {
	tx, err := p.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if category exists
	if err := p.checkCategoryExists(ctx, tx, opts.CategoryId); err != nil {
		return err
	}

	// Create a base product
	productId, err := p.createBaseProduct(ctx, tx, opts)
	if err != nil {
		return err
	}

	// Create product variants
	if err := p.createProductVariants(ctx, tx, productId, opts.Variants); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *ProductEntityStore) AddProductVariant(ctx context.Context, productId int64, opts *ProductVariantCreateUpdate) error {
	tx, err := p.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := p.checkProductExists(ctx, tx, productId); err != nil {
		return err
	}

	if err := p.addProductVariant(ctx, tx, productId, opts); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *ProductEntityStore) createBaseProduct(ctx context.Context, tx pgx.Tx, opts *ProductCreateUpdate) (int64, error) {
	var productId int64
	err := tx.QueryRow(ctx,
		`INSERT INTO "products" ("name", "description", "base_price", "category_id", "employee_id") 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING "id"`,
		opts.Name, opts.Description, opts.BasePrice, opts.CategoryId, opts.EmployeeId,
	).Scan(&productId)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}
	return productId, nil
}

func (p *ProductEntityStore) createProductVariants(ctx context.Context, tx pgx.Tx, productId int64, variants []ProductVariantCreateUpdate) error {
	for _, variant := range variants {
		if err := p.addProductVariant(ctx, tx, productId, &variant); err != nil {
			return fmt.Errorf("failed to create product variant: %w", err)
		}
	}
	return nil
}

func (p *ProductEntityStore) addProductVariant(ctx context.Context, tx pgx.Tx, productId int64, opts *ProductVariantCreateUpdate) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO "product_variants" ("product_id", "size_id", "color_id", "price", "stock")
			VALUES ($1, $2, $3, $4, $5)`,
		productId, opts.SizeId, opts.ColorId, opts.Price, opts.Stock,
	)
	if err != nil {
		return fmt.Errorf("failed to add product variant: %w", err)
	}
	return nil
}

func (p *ProductEntityStore) checkCategoryExists(ctx context.Context, tx pgx.Tx, categoryId int64) error {
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "categories" WHERE "id" = $1)`, categoryId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check category: %w", err)
	}
	if !exists {
		return fmt.Errorf("category with id '%d' not found", categoryId)
	}
	return nil
}

func (p *ProductEntityStore) checkProductExists(ctx context.Context, tx pgx.Tx, productId int64) error {
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "products" WHERE "id" = $1)`, productId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check product: %w", err)
	}
	if !exists {
		return fmt.Errorf("product with id '%d' not found", productId)
	}
	return nil
}
