package main

import (
	"./concurrency"
	"time"
)

func main() {
	go concurrency.CheckTimeout()
	concurrency.GenerateHandlerList()
	time.Sleep(time.Duration(30) * time.Second)
}
