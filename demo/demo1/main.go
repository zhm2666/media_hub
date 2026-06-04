package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var logger *log.Logger

type Prize struct {
	Name  string
	Total int32
	Left  int32
}

type Task struct {
	ID int
	R  *rand.Rand
}

type Result struct {
	Msg string
}

var prizes []*Prize

func init() {
	prizes = []*Prize{
		{Name: "一等奖", Total: 5},
		{Name: "二等奖", Total: 10},
		{Name: "三等奖", Total: 100},
	}
	for _, p := range prizes {
		p.Left = p.Total
	}
}

func InitLog() {
	f, _ := os.Create("./lottery_channel.log")
	logger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
}

// Worker
func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		result := draw(task.R)
		results <- Result{Msg: result}
	}
}

// 抽奖逻辑（无 channel，纯计算）
func draw(r *rand.Rand) string {
	code := r.Int31n(115) // 1+3+10+101（含未中奖）
	switch {
	case code < 1:
		return tryAward(prizes[0], "一等奖")
	case code < 4:
		return tryAward(prizes[1], "二等奖")
	case code < 14:
		return tryAward(prizes[2], "三等奖")
	default:
		return "未中奖"
	}
}

func tryAward(p *Prize, name string) string {
	newLeft := atomic.AddInt32(&p.Left, -1)
	if newLeft >= 0 {
		return "中奖：" + name
	}
	atomic.AddInt32(&p.Left, 1)
	return "奖品已发完"
}

func main() {
	InitLog()
	taskChan := make(chan Task, 100)
	resultChan := make(chan Result, 100)
	var wg sync.WaitGroup
	// 启动 Worker Pool
	workerCount := 5
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(i, taskChan, resultChan, &wg)
	}
	// 发送任务
	go func() {
		for i := 0; i < 10050; i++ {
			taskChan <- Task{
				ID: i,
				R:  rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))),
			}
		}
		close(taskChan)
	}()
	// 结果收集
	done := make(chan struct{})
	go func() {
		for res := range resultChan {
			logger.Println(res.Msg)
		}
		close(done)
	}()
	wg.Wait()
	close(resultChan)
	<-done
	fmt.Println("抽奖完成")
}
