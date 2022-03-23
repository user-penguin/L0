package model

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func (p *Payment) Insert(orderPK uint64, conn pgx.Conn) (uint64, error) {
	row := conn.QueryRow(context.Background(),
		"INSERT INTO payment ("+
			"document_id,"+
			"transaction,"+
			"request_id,"+
			"currency,"+
			"provider,"+
			"amount,"+
			"payment_dt,"+
			"bank,"+
			"delivery_cost,"+
			"goods_total,"+
			"custom_fee) "+
			"VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id",
		orderPK,
		p.Transaction,
		p.RequestId,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDt,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee,
	)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id, nil
}
