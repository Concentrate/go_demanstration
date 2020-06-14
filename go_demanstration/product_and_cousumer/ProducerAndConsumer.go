package main

import (
	"fmt"
	"time"
)

type blockingQueue struct {
	queue chan int
}

func (queue *blockingQueue) push(num int) {
	queue.queue <- num
}

func createBlockingQueue(size int) *blockingQueue {
	queue := blockingQueue{queue: make(chan int, size)}
	return &queue
}

func (bq *blockingQueue) pull() int {
	return <-bq.queue
}

func main() {
	tmpBq := createBlockingQueue(5)
	producer := func(bq *blockingQueue) {
		for i := 0; i < 1000; i++ {
			bq.push(i)
		}
	}

	for i := 0; i < 500; i++ {
		go producer(tmpBq)
	}
	go func(bq *blockingQueue) {
		for {
			fmt.Printf("%v ", bq.pull())
		}
	}(tmpBq)
	time.Sleep(time.Hour)
}
