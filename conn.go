package squash

type ConnId uint32
type Reqtype int

const (
	WriteReq Reqtype = iota
	FlushReq
)

type Request struct {
	id  ConnId
	typ Reqtype
}

type Conn struct {
	id  ConnId
	mux *ConnMux
	rch chan bool
	wch chan bool
}

func (p *Conn) Read(bs []byte) error {
	<-p.rch
	err := p.mux.read(bs)
	p.rch <- true
	return err
}

func (p *Conn) Write(bs []byte) error {
	p.mux.reqWrite <- Request{id: p.id, typ: WriteReq}
	<-p.wch
	err := p.mux.write(bs)
	p.wch <- true
	return err
}

func (p *Conn) Close() error {
	return p.mux.delConn(p.id)
}
