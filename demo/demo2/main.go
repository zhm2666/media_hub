package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var logger *log.Logger

// ======================
// 协程池
// ======================

type Task func(ctx context.Context) error

type Pool struct {
	tasks   chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	size    int
	running int32
}

func NewPool(size int, queueSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		tasks:  make(chan Task, queueSize),
		size:   size,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	atomic.AddInt32(&p.running, 1)
	defer atomic.AddInt32(&p.running, -1)
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[Worker %d] panic: %v", id, r)
					}
				}()
				_ = task(p.ctx)
			}()
		}
	}
}

func (p *Pool) Submit(task Task) {
	select {
	case p.tasks <- task:
	case <-p.ctx.Done():
	}
}

func (p *Pool) Stop() {
	close(p.tasks)
	p.cancel()
	p.wg.Wait()
}

// ======================
// 抽奖业务
// ======================

type Prize struct {
	Name  string
	Total int32
	Left  int32
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
	f, _ := os.Create("./lottery_pool.log")
	logger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
}

func draw(r *rand.Rand) error {
	code := r.Int31n(115)
	switch {
	case code < 1:
		return award(prizes[0], "一等奖")
	case code < 4:
		return award(prizes[1], "二等奖")
	case code < 14:
		return award(prizes[2], "三等奖")
	default:
		logger.Println("未中奖")
	}
	return nil
}

func award(p *Prize, name string) error {
	newLeft := atomic.AddInt32(&p.Left, -1)
	if newLeft >= 0 {
		logger.Printf("中奖：%s，剩余：%d", name, newLeft)
		return nil
	}
	atomic.AddInt32(&p.Left, 1)
	logger.Println("奖品已发完")
	return nil
}

func main() {
	InitLog()
	pool := NewPool(5, 100)
	pool.Start()
	for i := 0; i < 10050; i++ {
		i := i
		pool.Submit(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return nil
			default:
				return draw(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
			}
		})
	}
	pool.Stop()
	fmt.Println("抽奖完成")
}
