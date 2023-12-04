package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id           primitive.ObjectID `json:"id, omitempty"`
	Nama_product string             `json:"nama_product, omitempty" validate:"required"`
	Jumlah       int                `json:"jumlah, omitempty" validate:"required"`
	Dibuat       string             `json:"dibuat, omitempty" validate:"required"`
	Ketersediaan bool               `json:"ketersediaan, omitempty" validate:"required"`
}
