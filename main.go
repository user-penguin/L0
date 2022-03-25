package main

import (
	"L0/model"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
)

const psqlUrl = "postgresql://localhost/L0?user=intern&password=.hbqufufhby"

func main() {
	// мапа для временного хранения данных
	cache := make(map[string]model.Order)
	// устанавливаем соединение с nats-streaming-server
	sc, err := stan.Connect("test-cluster", "reader")
	if err != nil {
		log.Printf("Connections error, reason:%s", err)
		panic("error connection") // если не законнектили, нечего продолжать
	}
	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {
			log.Printf("Error closing stan connection. %s", err)
		}
	}(sc)

	// устанавливаем соединение с постгресом
	conn, err := pgx.Connect(context.Background(), psqlUrl)
	if err != nil {
		log.Printf("Error with psql connection. %s", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Printf("Error closing psql connection. %s", err)
		}
	}(conn, context.Background())

	// основное тело, где происходит принятие сообщений из nats-streaming
	_, err = sc.Subscribe("wild", func(msg *stan.Msg) {
		var order = &model.Order{}
		err := json.Unmarshal(msg.Data, order)
		if err != nil {
			log.Printf("Error in NATS-Order unmarhaling. %s", err)
		}
		_, _ = insertOrder(*order, *conn)
		cache[order.OrderUid] = *order
	})
	if err != nil {
		log.Println(err)
	}

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

func insertOrder(ord model.Order, conn pgx.Conn) (string, error) {
	orderId, _ := ord.Insert(conn)
	_, _ = ord.Delivery.Insert(orderId, conn)
	_, _ = ord.Payment.Insert(orderId, conn)
	for _, item := range ord.Items {
		_, _ = item.Insert(orderId, conn)
	}
	log.Printf("Last inserted order id: %s", orderId)
	return orderId, nil
}
