package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/egaevan/merchant/model"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &Product{
		DB: db,
	}
}

func (p *Product) FindOne(ctx context.Context, productID int) (*model.ProductDetail, error) {
	query := `
			SELECT 
				p.id,
				p.name,
				p.sku,
				p.path,
				pd.price,
				pd.stock 
			FROM 
				product p
			JOIN
				product_detail pd ON p.id = pd.product_id
			WHERE
				p.id = ? AND p.flag_aktif = 1`

	prod := model.ProductDetail{}

	err := p.DB.QueryRowContext(ctx, query, productID).Scan(&prod.Product.ID, &prod.Product.Name, &prod.Product.Sku, &prod.Product.Path, &prod.Price, &prod.Stock)
	if err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error querying find product uuid %s", err.Error())
	}

	return &prod, nil
}

func (p *Product) Fetch(ctx context.Context) (result []model.ProductDetail, err error) {
	query := `
			SELECT 
				p.id,
				p.name,
				p.sku,
				p.path,
				pd.price,
				pd.stock 
			FROM 
				product p
			JOIN
				product_detail pd ON p.id = pd.product_id
			WHERE
				flag_aktif = 1`

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Error(errRow)
		}
	}()

	result = make([]model.ProductDetail, 0)

	for rows.Next() {
		t := model.ProductDetail{}
		err = rows.Scan(
			&t.Product.ID,
			&t.Product.Name,
			&t.Product.Sku,
			&t.Product.Path,
			&t.Price,
			&t.Stock,
		)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (p *Product) Store(ctx context.Context, product model.ProductDetail) error {

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query1 := `
				INSERT INTO product
					(name, sku, path)
				VALUES
					(?, ?, ?)
			`

	res, err := tx.ExecContext(ctx, query1,
		product.Product.Name, product.Product.Sku, product.Product.Path)

	if err != nil {
		return err
	}

	producID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	query2 := `
				INSERT INTO product_detail
					(product_id, price, stock)
				VALUES
					(?, ?, ?)
			`

	_, err = tx.ExecContext(ctx, query2,
		producID, product.Price, product.Stock)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *Product) Update(ctx context.Context, product model.ProductDetail, productID int) error {

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query1 := `
				UPDATE 
					product
				SET
					name = ?, 
					sku = ?, 
					path = ?
				WHERE
					id = ?
			`

	_, err = tx.ExecContext(ctx, query1,
		product.Product.Name, product.Product.Sku, product.Product.Path, productID)

	if err != nil {
		return err
	}

	query2 := `
				UPDATE 
					product_detail
				SET
					price = ?, 
					stock = ?
				WHERE
					product_id = ?
			`

	_, err = tx.ExecContext(ctx, query2,
		product.Price, product.Stock, productID)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// func (p *Product) Update(ctx context.Context, product model.ProductUpdate, productID int) error {
// 	query := `
// 				UPDATE
// 					product
// 				SET
// 					name = ?,
// 					sku = ?,
// 					path = ?
// 				WHERE
// 					id = ?
// 			`

// 	_, err := p.DB.ExecContext(ctx, query,
// 		product.Name, product.Sku, product.Path, productID)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (p *Product) Delete(ctx context.Context, productID int) error {
	query := `
				UPDATE 
					product
				SET
					flag_aktif = 0
				WHERE
					id = ?
			`

	_, err := p.DB.ExecContext(ctx, query, productID)

	if err != nil {
		return err
	}

	return nil
}
