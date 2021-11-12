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

type ProductCount struct {
	Count int
}

type ProductType struct {
	TypeName string
}

type Master struct {
	ID             int
	FIO            string
	Specialization string
}
type NewMaster struct {
	FIO            string
	Specialization string
}

type ErrorMessage struct {
	Message string
}
