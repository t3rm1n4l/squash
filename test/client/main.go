package main

import (
	"fmt"
	"github.com/t3rm1n4l/squash"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	reqSize  = 200
	respSize = 200
)

var (
	pool sync.Pool
	req  []byte
)

func init() {
	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, reqSize+respSize)
		},
	}

	req = make([]byte, respSize)
}

func main() {
	f, _ := os.Create("prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wgc sync.WaitGroup
	var count uint64

	addr := os.Args[1]
	n, _ := strconv.Atoi(os.Args[2])
	thr, _ := strconv.Atoi(os.Args[3])
	nclients, _ := strconv.Atoi(os.Args[3])
	nperthr := n / thr
	c, _ := squash.NewClient(addr)

	t0 := time.Now()

	clientWorker := func(wgc *sync.WaitGroup) {
		defer wgc.Done()

		var wg sync.WaitGroup
		tm := time.Now()
		for i := 0; i < thr; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < nperthr; i++ {
					p := c.NewConn()
					p.Write(req)
					buf := pool.Get()
					p.Read(buf.([]byte)[:respSize])
					pool.Put(buf)
					p.Close()
					x := atomic.AddUint64(&count, 1)
					if x == 1000000 {
						tx := time.Now()
						fmt.Println(float64(x)/tx.Sub(tm).Seconds(), "req/sec")
						tm = tx
						atomic.StoreUint64(&count, 0)
					}

				}
			}()
		}

		wg.Wait()
		tx := time.Now()
		fmt.Println(float64(count)/tx.Sub(tm).Seconds(), "req/sec")
	}

	for i := 0; i < nclients; i++ {
		wgc.Add(1)
		go clientWorker(&wgc)
	}

	wgc.Wait()

	fmt.Println(float64(n)/time.Since(t0).Seconds(), "req/sec")
}
