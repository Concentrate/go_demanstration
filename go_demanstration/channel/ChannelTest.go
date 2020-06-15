package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan interface{})

	var group = sync.WaitGroup{}
	group.Add(3)
	go  runGenerNum(ch, group)
	go  runGenerNum(ch, group)

	go func() {
		for {
			fmt.Printf("%v   ",<-ch)
			time.Sleep(time.Second)
		}
		group.Done()
	}()
	group.Wait()
}

func runGenerNum(ch chan interface{}, group sync.WaitGroup) {
		index := 0
		for {
			ch <- index
			index = index + 1
		}
		group.Done()

}
