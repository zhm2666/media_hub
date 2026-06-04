package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
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

type PrizeMeta struct {
	Prize *Prize
	End   int32 // 权重前缀和
}

var (
	//prizes      []*Prize
	prizeMeta   []PrizeMeta
	totalWeight int32
)

func init() {
	raw := []*Prize{
		{Name: "一等奖", Total: 5},
		{Name: "二等奖", Total: 10},
		{Name: "三等奖", Total: 100},
	}
	var sum int32
	for _, p := range raw {
		p.Left = p.Total
		sum += p.Total
		prizeMeta = append(prizeMeta, PrizeMeta{
			Prize: p,
			End:   sum,
		})
	}
	totalWeight = sum
	//prizes = raw
}

func InitLog() {
	f, _ := os.Create("./lottery.log")
	logger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
}

func GetPrize(r *rand.Rand) {
	code := r.Int31n(totalWeight)
	// ✅ 二分查找
	idx := sort.Search(len(prizeMeta), func(i int) bool {
		return prizeMeta[i].End > code
	})
	selected := prizeMeta[idx].Prize
	newLeft := atomic.AddInt32(&selected.Left, -1)
	if newLeft >= 0 {
		logger.Printf("中奖：%s，剩余：%d", selected.Name, newLeft)
	} else {
		atomic.AddInt32(&selected.Left, 1)
		logger.Println("奖品已发完")
	}
}

func main() {
	InitLog()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 1005; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for j := 0; j < 10; j++ {
				GetPrize(r)
			}
		}()
	}
	wg.Wait()
	fmt.Println("耗时：", time.Since(start))
}
