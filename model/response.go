package model

import (
	"golang.org/x/exp/constraints"
)

type CustomData interface {
	constraints.Ordered | map[string]string | []ApiError | []ApiBookData | ApiBookData | []Author | Author | []Book | Book | []Publisher | Publisher
}

type Response[T CustomData] struct {
	Message	string	`json:"message"`
	Count	int64	`json:"count"`
	Page	int64	`json:"page"`
	Data	T		`json:"data"`
	Errors	[]ApiError	`json:"errors"`	
}
