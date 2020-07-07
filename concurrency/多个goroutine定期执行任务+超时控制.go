package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var (
	taskList   = sync.Map{}
	taskNotify = sync.Map{}
)

const TimeOut = 3

func GenerateHandlerList() {
	for i := 0; i < 5; i++ {
		quit := make(chan struct{})
		taskList.LoadOrStore(i, time.Now().Unix())
		taskNotify.LoadOrStore(i, quit)
		go handleTask(i, quit)
	}
}

func handleTask(taskID int, quit chan struct{}) {
	t := time.NewTicker(time.Duration(1) * time.Second)
	select {
	case <-quit:
		t.Stop()
		fmt.Printf("%d%s", taskID, "-task timeout===========\n")
		return
	case <-t.C:
		fmt.Printf("%d%s", taskID, " is excuting\n")
	}
}

func CheckTimeout() {
	t := time.NewTicker(time.Duration(2) * time.Second)
	for {
		select {
		case <-t.C:
			taskList.Range(func(key, value interface{}) bool {
				taskID := key.(int)
				createTime := value.(int64)
				fmt.Printf("%s%d\n", "check task ", taskID)
				if time.Now().Unix()-createTime > TimeOut {
					quit, _ := taskNotify.Load(taskID)
					//注意：非nil的结构体需要加两对花括号
					quit.(chan struct{}) <- struct{}{}
				}
				return true
			})

		}
	}
}
