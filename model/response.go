package model

import (
	"golang.org/x/exp/constraints"
)

type CustomData interface {
	constraints.Ordered | map[string]string | Book | []Book
}

type Response[T CustomData] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
