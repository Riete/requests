package requests

import (
	"sync"
)

// RequestPool both RequestPool and SessionPool use sync.Pool to cache and reuse Request when you want to do bulk
// http requests with separate http.Client, default use WithClient(NewClient) Option when creating pool
// Option can be provided to Get to set Request properties, don't forget to call Put when http request done
type RequestPool struct {
	p sync.Pool
}

func (rp *RequestPool) Get(options ...Option) *Request {
	item := rp.p.Get()
	r := item.(*Request)
	for _, option := range options {
		option(r)
	}
	return r
}

// Put to avoid connection leaks, call CloseIdleConnections before put back to pool
func (rp *RequestPool) Put(r *Request) {
	r.CloseIdleConnections()
	rp.p.Put(r)
}

func NewRequestPool() *RequestPool {
	return &RequestPool{
		p: sync.Pool{
			New: func() any {
				return NewRequest(WithClient(NewClient()))
			},
		},
	}
}

type SessionPool struct {
	*RequestPool
}

func NewSessionPool(options ...Option) *SessionPool {
	return &SessionPool{
		&RequestPool{
			p: sync.Pool{
				New: func() any {
					return NewSession(WithClient(NewClient()))
				},
			},
		},
	}
}
