package main

import (
	//"fmt"
	"github.com/t3rm1n4l/squash"
	"net"
	"os"
	"sync"
	//	"time"
)

const (
	reqSize  = 200
	respSize = 200
)

var (
	pool sync.Pool
	resp []byte
)

func init() {
	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, reqSize+respSize)
		},
	}

	resp = make([]byte, respSize)
}

func callback(p net.Conn) {
	buf := pool.Get()
	p.Read(buf.([]byte)[:reqSize])
	pool.Put(buf)
	// time.Sleep(time.Millisecond * 20)
	p.Write(resp)
	p.Close()
}

func main() {
	addr := os.Args[1]
	s := squash.NewServer(addr, callback)
	s.Run()
}
