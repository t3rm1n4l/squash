package main

import (
	//"fmt"
	"github.com/t3rm1n4l/squash"
	"os"
	"sync"
)

const (
	reqSize  = 100
	respSize = 100
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

func callback(p *squash.Conn) {
	buf := pool.Get()
	p.Read(buf.([]byte)[:reqSize])
	pool.Put(buf)

	p.Write(resp)
	p.Close()
}

func main() {
	addr := os.Args[1]
	s := squash.NewServer(addr, callback)
	s.Run()
}
