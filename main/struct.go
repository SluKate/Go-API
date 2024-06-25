package main

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
	Measure  int     `json:"measure"`
}

type Measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
