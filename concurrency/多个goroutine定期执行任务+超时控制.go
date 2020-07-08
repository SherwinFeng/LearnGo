/*
 * FileName:
 * Author:      SherwinFeng
 * @contact:    gotzqf@163.com
 * @time:       2020/07/08 21:03
 * @version:    1.0
 * Description:	1.适用场景为：
				有多个goroutine需要周期性地执行任务,每次执行任务时会更新自己的时间戳
				如果通过超时检查发现某个goroutine没有执行任务了, 就退出这个goroutine
				不适用于：
				检测某个goroutine执行任务的时间是否超时，如果要求检测，需要在handleTask中的
				执行业务代码的case中加入超时控制
				2.例子：
				通过Websocket建立了一个room, 前端定期向这个room发renewal消息，更新room对应的
				时间戳，后端会定期向这个room里广播消息，也会检查对应的时间戳是否过期
				如果过期说明前端已经退出这个room了，那么就删除这个room, 结束往这个room发广播的goroutine的运行
				3.注：
				demo中并没有模拟room超时的情况，如果要模拟将handleTask函数中更新taskList注释即可
 * Changelog:
*/
package concurrency

import (
	"fmt"
	"sync"
	"time"
)

var (
	//<taskID, createTime> 用于判断是否超时
	taskList = sync.Map{}
	//<taskID, channel> 用于超时通知
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
	//需要加入for循环 因为select是自上往下执行不会回头 所以如果不加for
	//虽然每隔1s都会执行handleTask 但只是往t.C中写数据 不会执行select了
	for {
		select {
		case <-quit:
			t.Stop()
			fmt.Printf("%d%s", taskID, "-task timeout===========\n")
			return
		case <-t.C:
			fmt.Printf("%d%s", taskID, " is excuting\n")
			taskList.Store(taskID, time.Now().Unix())
		}
	}
}

func CheckTimeout() {
	t := time.NewTicker(time.Duration(3) * time.Second)
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
