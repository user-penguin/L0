package main

import (
	"L0/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"sync"
)

const psqlUrl = "postgresql://localhost/L0?user=intern&password=.hbqufufhby"

func main() {
	//cache := make(map[string]model.Order)
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

	// при старте сразу получаем мапу, где есть все ордера из БД
	// далее просто будем туда складывать новые, чтобы  можно было быстро их отдавать при запросе
	cache := getAllOrders(*conn)

	go serve()

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

func getAllOrders(conn pgx.Conn) map[string]model.Order {
	allOrders := make(map[string]model.Order)
	items := getItems(conn)
	documentSQL := "SELECT d.order_uid,d.track_number,d.entry,d.locale,d.internal_signature,d.customer_id,d.delivery_service," +
		"d.shardkey,d.sm_id,d.date_created,d.oof_shard,p.transaction,p.request_id,p.currency,p.provider,p.amount," +
		"cast(extract(epoch from p.payment_dt) as bigint),p.bank,p.delivery_cost,p.goods_total,p.custom_fee,del.name," +
		"del.phone,del.zip,del.city,del.address,del.region,del.email " +
		"FROM document d " +
		"JOIN payment p ON d.order_uid = p.document_id " +
		"JOIN delivery del ON d.order_uid = del.document_id"
	rows, err := conn.Query(context.Background(), documentSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var o model.Order
		err := rows.Scan(
			&o.OrderUid,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerId,
			&o.DeliveryService,
			&o.Shardkey,
			&o.SmId,
			&o.DateCreated,
			&o.OofShard,
			&o.Payment.Transaction,
			&o.Payment.RequestId,
			&o.Payment.Currency,
			&o.Payment.Provider,
			&o.Payment.Amount,
			&o.Payment.PaymentDt,
			&o.Payment.Bank,
			&o.Payment.DeliveryCost,
			&o.Payment.GoodsTotal,
			&o.Payment.CustomFee,
			&o.Delivery.Name,
			&o.Delivery.Phone,
			&o.Delivery.Zip,
			&o.Delivery.City,
			&o.Delivery.Address,
			&o.Delivery.Region,
			&o.Delivery.Email)
		if err != nil {
			log.Fatal(err)
		}
		itemsByUid := items[o.OrderUid]
		for _, item := range itemsByUid {
			o.Items = append(o.Items, item)
		}
		allOrders[o.OrderUid] = o
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return allOrders
}

func getItems(conn pgx.Conn) map[string][]model.Item {
	itemsSQL := "SELECT document_id,chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status " +
		"FROM item"
	rows, err := conn.Query(context.Background(), itemsSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	items := make(map[string][]model.Item)
	var item model.Item
	var uid string
	for rows.Next() {
		err := rows.Scan(
			&uid,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status)
		if err != nil {
			log.Fatal(err)
		}
		items[uid] = append(items[uid], item)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return items
}

func serve() {
	port := ":8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Добро пожаловать в мир Golang\n")
		log.Printf("Method:%s, Request URI:%s", r.Method, r.RequestURI)
	})
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("Server dosn't start: %s", err)
	}
	log.Printf("Server is running now on localhost%s", port)
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
