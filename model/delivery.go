package model

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) Insert(orderPK uint64, conn pgx.Conn) (uint64, error) {
	row := conn.QueryRow(context.Background(),
		"INSERT INTO delivery ("+
			"document_id,"+
			"name,"+
			"phone,"+
			"zip,"+
			"city,"+
			"address,"+
			"region,"+
			"email) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id",
		orderPK,
		d.Name,
		d.Phone,
		d.Zip,
		d.City,
		d.Address,
		d.Region,
		d.Email,
	)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id, nil
}
