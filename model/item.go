package model

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (i *Item) Insert(orderPK uint64, conn pgx.Conn) (uint64, error) {
	row := conn.QueryRow(context.Background(),
		"INSERT INTO item ("+
			"document_id,"+
			"chrt_id,"+
			"track_number,"+
			"price,"+
			"rid,"+
			"name,"+
			"sale,"+
			"size,"+
			"total_price,"+
			"nm_id,"+
			"brand,"+
			"status) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id",
		orderPK,
		i.ChrtId,
		i.TrackNumber,
		i.Price,
		i.Rid,
		i.Name,
		i.Sale,
		i.Size,
		i.TotalPrice,
		i.NmId,
		i.Brand,
		i.Status,
	)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id, nil
}
