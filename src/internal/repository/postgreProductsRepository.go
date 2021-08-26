package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"main/src/internal/models"
)

type PostgreProductRepository struct {
	conn *pgx.Conn
}

func (p *PostgreProductRepository) AddProduct(prod []models.Product) error {

	b := &pgx.Batch{}
	for _, v := range prod {
		b.Queue("INSERT INTO products(name,price) values($1,$2)", (v.Name), v.Price)
	}

	res := p.conn.SendBatch(context.TODO(), b)
	err := res.Close()
	if err != nil {
		log.Println(err)
	}
	return nil
}

func NewPostgreProductRepository(conn *pgx.Conn) PostgreProductRepository {
	return PostgreProductRepository{conn: conn}
}
