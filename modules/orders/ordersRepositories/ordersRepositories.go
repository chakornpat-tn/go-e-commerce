package ordersRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/chakornpat-tn/go-rest-api/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOrder(orderId string) (*orders.Order, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) FindOrder(orderId string) (*orders.Order, error) {
	query := `
	SELECT
  to_jsonb("t")
FROM
  (
    SELECT
      "o"."id",
      "o"."user_id",
      "o"."transfer_slip",
      (
        SELECT
          array_to_json(array_agg("pt"))
        FROM
          (
            SELECT
              "spo"."id",
              "spo"."qty",
              "spo"."product"
            FROM
              "products_orders" "spo"
            WHERE
              "spo"."order_id" = "o"."id"
          ) AS "pt"
      ) AS "products",
      "o"."address",
      "o"."contact",
      (
        SELECT
          SUM(
            COALESCE(
              ("po"."product" ->> 'price') :: FLOAT * ("po"."qty") :: FLOAT,
              0
            )
          )
        FROM
          "products_orders" "po"
        WHERE
          "po"."order_id" = "o"."id"
      ) AS "total_paid",
      "o"."status",
      "o"."created_at",
      "o"."updated_at"
    FROM
      "orders" "o"
    WHERE
      "o"."id" = $1
  ) AS "t";`

	orderData := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Products:     make([]*orders.ProductsOrder, 0),
	}

	raw := make([]byte, 0)

	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order data: %v", err)
	}

	return orderData, nil
}
