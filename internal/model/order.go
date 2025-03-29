package model

import "time"

type Order struct {
	ID       int       `db:"id" json:"id" example:"10"`
	PetID    int       `db:"pet_id" json:"petId" example:"3"`
	Quantity int       `db:"quantity" json:"quantity" example:"2"`
	ShipDate time.Time `db:"ship_date" json:"shipDate" example:"2025-03-29T15:04:05Z"`
	Status   string    `db:"status" json:"status" example:"placed"`
	Complete bool      `db:"complete" json:"complete" example:"false"`
}
