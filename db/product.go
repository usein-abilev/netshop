package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ProductVariantEntity struct {
	Id    int64       `json:"id"`
	Size  SizeEntity  `json:"size"`
	Color ColorEntity `json:"color"`
	Price float64     `json:"price"`
	Stock int32       `json:"stock"`
}

type ProductEntity struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Id          int64                  `json:"id"`
	BasePrice   float64                `json:"base_price"`
	Category    CategoryEntity         `json:"category"`
	Variants    []ProductVariantEntity `json:"variants"`
}

type ProductGetEntitiesQueryOpts struct {
	CategoryIds []int64  `json:"category_ids,omitempty"`
	SizeIds     []int64  `json:"size_ids,omitempty"`
	ColorIds    []int64  `json:"color_ids,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
}

type ProductGetEntitiesOptions struct {
	// Query options for filtering products
	// If nil, no filtering is applied
	Query *ProductGetEntitiesQueryOpts

	// Limit is the maximum number of products to return. If 0, no limit is applied
	Limit int64
	// Offset is the number of products to skip. If 0, no offset is applied
	Offset int64

	// If OrderColumn is empty, default "id" is used
	OrderColumn string
	// If OrderDesc is false, default descending order is used
	OrderAsc bool
}

type ProductVariantCreateUpdate struct {
	SizeId  int64   `json:"size_id"`
	ColorId int64   `json:"color_id"`
	Price   float64 `json:"price"`
	Stock   int32   `json:"stock"`
}

type ProductCreateUpdate struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	CategoryId  int64                        `json:"category_id"`
	EmployeeId  int64                        `json:"employee_id"`
	BasePrice   float64                      `json:"base_price"`
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

func (p *ProductEntityStore) GetById(id int64) (ProductEntity, error) {
	row := p.db.Connection.QueryRow(p.db.Context, `select "id", "name", "description", "base_price" from "products" where id = $1`, id)
	var product ProductEntity
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.BasePrice)
	if err != nil {
		return ProductEntity{}, err
	}
	return product, nil
}

func (p *ProductEntityStore) GetEntities(opts *ProductGetEntitiesOptions) ([]ProductEntity, error) {
	query := strings.Builder{}
	query.WriteString(`select
		"products"."id",
		"products"."name",
		"products"."description",
		"products"."base_price",
		"products"."category_id",
		"categories"."name",
		"product_variants"."id",
		"product_variants"."size_id",
		"product_variants"."color_id",
		"product_variants"."price",
		"product_variants"."stock",
		"sizes"."name",
		"colors"."name"
	from "products"
	join "categories" on "products"."category_id" = "categories"."id"
	join "product_variants" on "products"."id" = "product_variants"."product_id"
	join "sizes" on "product_variants"."size_id" = "sizes"."id"
	join "colors" on "product_variants"."color_id" = "colors"."id"`)

	whereConditions := []string{}
	convertToSqlSeq := func(ids []int64) string {
		return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")
	}
	addWhere := func(condition string) {
		whereConditions = append(whereConditions, condition)
	}

	if opts != nil {
		if opts.Query != nil {
			if len(opts.Query.CategoryIds) > 0 {
				addWhere(fmt.Sprintf(`"products"."category_id" in (%s)`, convertToSqlSeq(opts.Query.CategoryIds)))
			}
			if len(opts.Query.SizeIds) > 0 {
				addWhere(fmt.Sprintf(`"product_variants"."size_id" in (%s)`, convertToSqlSeq(opts.Query.SizeIds)))
			}
			if len(opts.Query.ColorIds) > 0 {
				addWhere(fmt.Sprintf(`"product_variants"."color_id" in (%s)`, convertToSqlSeq(opts.Query.ColorIds)))
			}
			if opts.Query.MinPrice != nil {
				addWhere(fmt.Sprintf(`"product_variants"."price" >= %f`, *opts.Query.MinPrice))
			}
			if opts.Query.MaxPrice != nil {
				addWhere(fmt.Sprintf(`"product_variants"."price" <= %f`, *opts.Query.MaxPrice))
			}
		}

		if opts.Limit > 0 {
			query.WriteString(fmt.Sprintf(" limit %d", opts.Limit))
		}
		if opts.Offset > 0 {
			query.WriteString(fmt.Sprintf(" offset %d", opts.Offset))
		}

		// TODO: Add order by
	}

	if len(whereConditions) > 0 {
		query.WriteString(" where ")
		query.WriteString(strings.Join(whereConditions, " and "))
	}

	rows, err := p.db.Connection.Query(p.db.Context, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productsMap := make(map[int64]*ProductEntity)
	for rows.Next() {
		var productId, variantId, sizeId, colorId int64
		var productName, productDescription, categoryName, sizeName, colorName string
		var basePrice, variantPrice float64
		var categoryId int64
		var stock int32

		err := rows.Scan(
			&productId, &productName, &productDescription, &basePrice, &categoryId, &categoryName,
			&variantId, &sizeId, &colorId, &variantPrice, &stock,
			&sizeName, &colorName,
		)
		if err != nil {
			return nil, err
		}

		product, exists := productsMap[productId]
		if !exists {
			product = &ProductEntity{
				Id:          productId,
				Name:        productName,
				Description: productDescription,
				BasePrice:   basePrice,
				Category:    CategoryEntity{Id: categoryId, Name: categoryName},
				Variants:    []ProductVariantEntity{},
			}
			productsMap[productId] = product
		}

		variant := ProductVariantEntity{
			Id:    variantId,
			Size:  SizeEntity{Id: sizeId, Name: sizeName},
			Color: ColorEntity{Id: colorId, Name: colorName},
			Price: variantPrice,
			Stock: stock,
		}
		product.Variants = append(product.Variants, variant)
	}

	products := make([]ProductEntity, 0, len(productsMap))
	for _, product := range productsMap {
		products = append(products, *product)
	}

	return products, nil
}

func (p *ProductEntityStore) Create(ctx context.Context, opts *ProductCreateUpdate) error {
	tx, err := p.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := p.checkCategoryExists(ctx, tx, opts.CategoryId); err != nil {
		return err
	}

	productId, err := p.createBaseProduct(ctx, tx, opts)
	if err != nil {
		return err
	}

	if err := p.createProductVariants(ctx, tx, productId, opts.Variants); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *ProductEntityStore) Update(ctx context.Context, opts *ProductCreateUpdate) error {
	// TODO: Implement
	return errors.New("not implemented")
}

func (p *ProductEntityStore) GetVariants(productId int64) ([]ProductVariantEntity, error) {
	query := `select 
			"product_variants"."id" as "variant_id", 
			"size_id", 
			"color_id", 
			"price", 
			"stock",
			"sizes"."name",
			"colors"."name"
		from "product_variants"
			left join "sizes" on "product_variants"."size_id" = "sizes"."id"
			left join "colors" on "product_variants"."color_id" = "colors"."id"
 		where "product_id" = $1
	`
	rows, err := p.db.Connection.Query(p.db.Context, query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []ProductVariantEntity = make([]ProductVariantEntity, 0)
	for rows.Next() {
		var variant ProductVariantEntity
		err := rows.Scan(&variant.Id, &variant.Size.Id, &variant.Color.Id, &variant.Price, &variant.Stock, &variant.Size.Name, &variant.Color.Name)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	return variants, nil
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
