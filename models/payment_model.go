package models

type Payment struct {
	Id              string
	Status          string
	CardNumber      string
	ExpirationMonth int
	ExpirationYear  int
	CurrencyCode    string
	Amount          int
	cvv             string
}
