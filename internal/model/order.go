package model

import "time"

type Order struct {
	ID       int       `db:"id" json:"id"`
	PetID    int       `db:"pet_id" json:"petId"`
	Quantity int       `db:"quantity" json:"quantity"`
	ShipDate time.Time `db:"ship_date" json:"shipDate"`
	Status   string    `db:"status" json:"status"`
	Complete bool      `db:"complete" json:"complete"`
}
