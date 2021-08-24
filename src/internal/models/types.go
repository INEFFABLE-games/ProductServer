package models

type Product struct {
	Name string `json:"name" bson:"name"`
	Price uint64 `json:"price" bson:"price"`
}