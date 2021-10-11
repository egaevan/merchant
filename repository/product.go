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

func (p *Product) FindOne(ctx context.Context, productID int) (*model.Product, error) {
	query := `
			SELECT 
				ID,
				name,
				sku,
				COALESCE(path,'')
			FROM 
				product
			WHERE
				ID = ? AND flag_aktif = 1`

	prod := model.Product{}

	err := p.DB.QueryRowContext(ctx, query, productID).Scan(&prod.ID, &prod.Name, &prod.Sku, &prod.Path)
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

func (p *Product) Fetch(ctx context.Context) (result []model.Product, err error) {
	query := `
			SELECT 
				ID,
				name,
				sku,
				COALESCE(path,'')
			FROM 
				product
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

	result = make([]model.Product, 0)

	for rows.Next() {
		t := model.Product{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Sku,
			&t.Path,
		)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (p *Product) Store(ctx context.Context, product model.Product) error {
	query := `
				INSERT INTO product
					(ID, name, sku, path)
				VALUES
					(?, ?, ?, ?)
			`

	_, err := p.DB.ExecContext(ctx, query,
		product.ID, product.Name, product.Sku, product.Path)

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Update(ctx context.Context, product model.ProductUpdate, productID int) error {
	query := `
				UPDATE 
					product
				SET
					name = ?, 
					sku = ?, 
					path = ?
				WHERE
					ID = ?
			`

	_, err := p.DB.ExecContext(ctx, query,
		product.Name, product.Sku, product.Path, productID)

	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete(ctx context.Context, productID int) error {
	query := `
				UPDATE 
					product
				SET
					flag_aktif = 0
				WHERE
					ID = ?
			`

	_, err := p.DB.ExecContext(ctx, query, productID)

	if err != nil {
		return err
	}

	return nil
}
