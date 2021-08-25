package models

type Product struct {
	Name  string `json:"name" bson:"name"`
	Price uint32 `json:"price" bson:"price"`
}
