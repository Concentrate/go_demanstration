package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	Subsriber chan interface{}
	FliterFun func(msg interface{}) bool
)

type Publisher struct {
	lock        sync.RWMutex
	subscirbers map[Subsriber]FliterFun
	timeOut     time.Duration
	bufferCount int
}

func NewPublisher(timeOut time.Duration, bufferCount int) *Publisher {
	return &Publisher{
		lock:        sync.RWMutex{},
		subscirbers: make(map[Subsriber]FliterFun),
		timeOut:     timeOut,
		bufferCount: bufferCount,
	}
}

func (pub *Publisher) Subscriber(fliter FliterFun) Subsriber {
	sub := make(Subsriber,pub.bufferCount)
	pub.lock.Lock()
	pub.subscirbers[sub] = fliter
	pub.lock.Unlock()
	return sub
}

func (pub *Publisher) RemoveSubscribe(tmpSub Subsriber) {
	pub.lock.Lock()
	defer pub.lock.Unlock()
	delete(pub.subscirbers, tmpSub)
	close(tmpSub)
}

func (pub *Publisher) ClosePublisher() {
	pub.lock.Lock()
	defer pub.lock.Unlock()
	tmpSaveSubArray := make([]Subsriber, 0)
	for sub, _ := range pub.subscirbers {
		tmpSaveSubArray = append(tmpSaveSubArray, sub)
		close(sub)
	}
	for _, va := range tmpSaveSubArray {
		delete(pub.subscirbers, va)
	}
}

func (pub *Publisher) SubscribeAllMessage() Subsriber {
	return pub.Subscriber(nil)
}

func (pub *Publisher) PublishMessage(message interface{}) {
	pub.lock.RLock()
	defer pub.lock.RUnlock()
	waitGroup := sync.WaitGroup{}
	for subscriber, tmpFilter := range pub.subscirbers {
		waitGroup.Add(1)
		go pub.publishToSubscribers(tmpFilter, message, subscriber, &waitGroup)
	}
	waitGroup.Wait()
}

func (pub *Publisher) publishToSubscribers(tmpFilter FliterFun, message interface{}, subscriber Subsriber, group *sync.WaitGroup) {
	defer group.Done()
	if tmpFilter != nil && !tmpFilter(message) {
		return
	}
	select {
	case subscriber <- message:
	case <-time.After(pub.timeOut):
	}
}

func main() {
	var Publisher = NewPublisher(10, 1)
	golangSubsci := Publisher.Subscriber(func(msg interface{}) bool {
		if va, ok := msg.(string); ok {
			return strings.Contains(va, "golang")
		}
		return false
	})
	normal := Publisher.SubscribeAllMessage()
	go func() {
		for {
			tmpStr := <-golangSubsci
			fmt.Println(tmpStr)
		}
	}()
	go func() {
		for {
			tmp := <-normal
			fmt.Println(tmp)
		}
	}()
	Publisher.PublishMessage("hello world")
	Publisher.PublishMessage("golang world")
	time.Sleep(5 * time.Second)
}
