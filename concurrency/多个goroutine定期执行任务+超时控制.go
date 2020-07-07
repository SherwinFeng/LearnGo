package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var(
	taskList = sync.Map{}
	taskNotify = sync.Map{}
)

const TimeOut = 3

func GenerateHandlerList(){
	for i := 0; i < 10; i++{
		quit := make(chan struct{})
		taskList.LoadOrStore(i,time.Now().Unix())
		taskNotify.LoadOrStore(i,quit)
		go handleTask(quit)
	}
}

func handleTask(quit chan struct{}){
	t := time.NewTicker(time.Duration(1) * time.Second)
	select {
	case <-quit:
		t.Stop()
		return
	case <-t.C:
		fmt.Println("I am running")
		}
}

func checkTimeout(){
	t := time.NewTicker(time.Duration(1) * time.Second)
	for{
		select {
		case <-t.C:
			taskList.Range(func(key, value interface{}) bool {
				taskID := key.(int)
				createTime := value.(int64)
				if time.Now().Unix() - createTime > TimeOut{
					quit, _ := taskNotify.Load(taskID)
					//注意：非nil的结构体需要加两对花括号
					quit.(chan struct{}) <- struct{}{}
				}
				return true
			})

		}
	}

}

