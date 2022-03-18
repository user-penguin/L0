package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"sync"
)

func main() {
	sc, _ := stan.Connect("test-cluster", "reader")
	defer sc.Close()
	sc.Subscribe("wild", func(msg *stan.Msg) {
		fmt.Printf("Get: %s\n", string(msg.Data))
	})

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
