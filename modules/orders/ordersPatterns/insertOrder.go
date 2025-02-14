package ordersPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/chakornpat-tn/go-rest-api/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrder() error
	commit() error
	getOrderId() string
}

type insertOrderBuilder struct {
	db  *sqlx.DB
	req *orders.Order
	tx  *sqlx.Tx
}

func (b *insertOrderBuilder) getOrderId() string {
	return b.req.Id
}

func (b *insertOrderBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}
func (b *insertOrderBuilder) insertOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	query := `
	INSERT INTO "orders" (
		"user_id",
		"contact",
		"address",
		"transfer_slip",
		"status"
	)
	VALUES
	($1, $2, $3, $4, $5)
		RETURNING "id";`

	if err := b.tx.QueryRowxContext(ctx, query, b.req.UserId, b.req.Contact, b.req.Address, b.req.TransferSlip, b.req.Status).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return err
	}
	return nil
}
func (b *insertOrderBuilder) insertProductsOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	query := `
	INSERT INTO "products_orders" (
		"order_id",
		"qty",
		"product"
	)
	VALUES`

	values := make([]any, 0)
	lastIdx := 0
	for i := range b.req.Products {
		values = append(values, b.req.Id, b.req.Products[i].Qty, b.req.Products[i].Product)

		if i != len(b.req.Products)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, lastIdx+1, lastIdx+2, lastIdx+3,
			)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, lastIdx+1, lastIdx+2, lastIdx+3,
			)
		}

		lastIdx += 3
	}
	if _, err := b.tx.ExecContext(ctx, query, values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product attachment failed: %v", err)
	}

	return nil
}
func (b *insertOrderBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertOrderBuilder(db *sqlx.DB, req *orders.Order) IInsertOrderBuilder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

type insertOrderEngineer struct {
	builder IInsertOrderBuilder
}

func InsertOrderEngineer(b IInsertOrderBuilder) *insertOrderEngineer {
	return &insertOrderEngineer{
		builder: b,
	}
}

func (en *insertOrderEngineer) InsertOrder() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertOrder(); err != nil {
		return "", err
	}
	if err := en.builder.insertProductsOrder(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getOrderId(), nil
}
