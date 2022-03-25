package model

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (o *Order) Insert(conn pgx.Conn) (string, error) {
	row := conn.QueryRow(context.Background(),
		"INSERT INTO document ("+
			"order_uid,"+
			"track_number,"+
			"entry,"+
			"locale,"+
			"internal_signature,"+
			"customer_id,"+
			"delivery_service,"+
			"shardkey,"+
			"sm_id,"+
			"date_created,"+
			"oof_shard) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING order_uid",
		o.OrderUid,
		o.TrackNumber,
		o.Entry,
		o.Locale,
		o.InternalSignature,
		o.CustomerId,
		o.DeliveryService,
		o.Shardkey,
		o.SmId,
		o.DateCreated,
		o.OofShard,
	)
	var order_uid string
	err := row.Scan(&order_uid)
	if err != nil {
		log.Printf("Failed in id search: %s", err)
	}
	return order_uid, nil
}
