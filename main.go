package main

import (
	"L0/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
)

func main() {
	sc, err := stan.Connect("test-cluster", "reader")
	if err != nil {
		log.Printf("Connections error, reason:%s", err)
		panic("error connection") // если не законнектили, нечего продолжать
	}
	defer func(sc stan.Conn) {
		err := sc.Close()
		if err != nil {
			log.Println(err)
		}
	}(sc)

	sc.Subscribe("wild", func(msg *stan.Msg) {
		var order = &model.Order{}
		err := json.Unmarshal(msg.Data, order)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("order_uid:%s, delivery.name: %s\n", order.OrderUid, order.Delivery.Name)
	})

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
