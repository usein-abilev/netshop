package db

import (
	"context"
	"fmt"
	"netshop/main/tools/sqb"
	"time"

	"github.com/jackc/pgx/v5"
)

type OrderItemEntity struct {
	Id               int64                 `json:"id"`
	OrderId          int64                 `json:"order_id"`
	ProductVariantId int64                 `json:"product_variant_id"`
	ProductVariant   *ProductVariantEntity `json:"product_variant"`
	Price            float64               `json:"price"`
	Quantity         uint32                `json:"quantity"`
}

type OrderEntity struct {
	Id              int64              `json:"id"`
	CustomerId      int64              `json:"customer_id"`
	Customer        *CustomerEntity    `json:"customer"`
	Status          string             `json:"status"`
	DeliveryAddress string             `json:"delivery_address"`
	DeliveryZipcode string             `json:"delivery_zipcode"`
	DeliveryCity    string             `json:"delivery_city"`
	DeliveryCountry string             `json:"delivery_country"`
	StatusDate      time.Time          `json:"status_date"`
	OrderDate       time.Time          `json:"order_date"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Items           []*OrderItemEntity `json:"items"`
}

type orderItemCreateUpdate struct {
	Price            float64
	ProductVariantId int64
	Quantity         int
}

type OrderGetAllOptions struct {
	CustomerId *int64
	Status     *string
}

type OrderCreateUpdateOptions struct {
	CustomerId int64
	OrderDate  *time.Time
	Customer   *CustomerCreateUpdate
	Status     string
	Delivery   struct{ Address, Zipcode, City, Country string }
	Items      []*orderItemCreateUpdate
}

type OrderEntityStore struct {
	db *DatabaseConnection
}

func NewOrderEntity(database *DatabaseConnection) *OrderEntityStore {
	return &OrderEntityStore{
		db: database,
	}
}

func (c *OrderEntityStore) Exists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "orders" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *OrderEntityStore) GetAll(options *OrderGetAllOptions) ([]OrderEntity, error) {
	builder := sqb.NewSQLQueryBuilder().
		Select(
			"orders.id",
			"orders.order_date",
			"orders.customer_id",
			"orders.status",
			"orders.delivery_address",
			"orders.delivery_zipcode",
			"orders.delivery_city",
			"orders.delivery_country",
			"orders.status_date",
			"orders.created_at",
			"orders.updated_at",
			"order_items.id",
			"order_items.order_id",
			"order_items.product_variant_id",
			"order_items.price",
			"order_items.quantity",
		).
		From("orders").
		LeftJoin("order_items", "order_items.order_id = orders.id").
		LeftJoin("product_variants", "product_variants.id = order_items.product_variant_id").
		InnerJoin("customers", "customers.id = orders.customer_id").
		OrderBy("orders.id", "desc")

	if options.CustomerId != nil {
		builder.AndWhere("orders.customer_id = $customerId")
		builder.SetParameter("customerId", options.CustomerId)
	}

	query, args := builder.Build()

	rows, err := c.db.Connection.Query(c.db.Context, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]OrderEntity, 0)
	for rows.Next() {
		var order OrderEntity
		err := rows.Scan(
			&order.Id,
			&order.OrderDate,
			&order.CustomerId,
			&order.Status,
			&order.DeliveryAddress,
			&order.DeliveryZipcode,
			&order.DeliveryCity,
			&order.DeliveryCountry,
			&order.StatusDate,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, nil
}

// Creates a new order in the database
// This methods can create a new customer by given customer object if it does not exist
func (c *OrderEntityStore) Create(ctx context.Context, options *OrderCreateUpdateOptions) (result *OrderEntity, err error) {
	tx, err := c.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	customerStore := NewCustomerEntityStore(c.db)
	customer, err := customerStore.GetById(options.CustomerId)
	if err != nil {
		return result, fmt.Errorf("failed to get customer: %w", err)
	}

	result = &OrderEntity{
		CustomerId:      customer.Id,
		Customer:        customer,
		Status:          options.Status,
		DeliveryAddress: options.Delivery.Address,
		DeliveryZipcode: options.Delivery.Zipcode,
		DeliveryCity:    options.Delivery.City,
		DeliveryCountry: options.Delivery.Country,
	}

	if options.OrderDate == nil {
		result.OrderDate = time.Now()
	}

	err = tx.QueryRow(c.db.Context, `
		insert into "orders" (
			customer_id, 
			status, 
			delivery_address, 
			delivery_zipcode, 
			delivery_city, 
			delivery_country, 
			status_date, 
			order_date) 
		values ($1, $2, $3, $4, $5, $6, $7, $8) 
		returning id, created_at, updated_at`,
		options.CustomerId,
		options.Status,
		options.Delivery.Address,
		options.Delivery.Zipcode,
		options.Delivery.City,
		options.Delivery.Country,
		time.Now(),
		options.OrderDate,
	).Scan(&result.Id, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return result, fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range options.Items {
		itemResult, err := c.createOrderItem(tx, result.Id, item)
		if err != nil {
			return result, fmt.Errorf("failed to create order item: %w", err)
		}
		result.Items = append(result.Items, itemResult)
	}

	return result, tx.Commit(ctx)
}

func (c *OrderEntityStore) createOrderItem(tx pgx.Tx, orderId int64, item *orderItemCreateUpdate) (*OrderItemEntity, error) {
	var result OrderItemEntity
	err := tx.QueryRow(c.db.Context, `
		with updated_stock as (
		 	update "product_variants"
			set stock = stock - $1
			where id = $2
			returning *
		)
		insert into "order_items" (order_id, product_variant_id, price, quantity)
		values ($3, $2, (select price from updated_stock), $1)
		returning id`,
		item.Quantity,
		item.ProductVariantId,
		orderId,
	).Scan(&result.Id)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
