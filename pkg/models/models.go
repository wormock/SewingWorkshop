package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Product struct {
	ID          int
	Type        string
	Cost        int
	Size        string
	Material    string
	MasterId    int
	CustomerId  int
	MasterFIO   string
	CustomerFIO string
}

type ProductType struct {
	TypeName string
}
