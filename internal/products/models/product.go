package models

type Product struct {
	Id          string
	Title       string
	Description string
	Price       float64
	Available   bool
	Images      []string
}
