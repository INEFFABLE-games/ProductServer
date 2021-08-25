package repository

import (
	"github.com/jackc/pgx"
	"main/src/internal/models"
)

type PostgreProductRepository struct {
	conn *pgx.Conn
}

func (p *PostgreProductRepository) AddProduct(prod models.Product) error {

	_, err := p.conn.Exec("insert into products(name,price) values($1,$2)", string(prod.Name), prod.Price)

	return err
}

func NewPostgreProductRepository(conn *pgx.Conn) PostgreProductRepository {
	return PostgreProductRepository{conn: conn}
}
