package requests

import (
	"sync"
)

// RequestPool both RequestPool and SessionPool use sync.Pool to cache and reuse Request when you want to do bulk
// http requests with separate http.Client, default use WithClient(NewClient) RequestOption when creating pool
// RequestOption can be provided to Get to set Request properties, don't forget to call Put when http request done
// options is default RequestOption
type RequestPool struct {
	p       sync.Pool
	options []RequestOption
}

func (rp *RequestPool) SetOptions(options ...RequestOption) {
	rp.options = options
}

func (rp *RequestPool) Get(options ...RequestOption) *Request {
	item := rp.p.Get()
	r := item.(*Request)
	options = append(rp.options, options...)
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

func NewRequestPool(options ...RequestOption) *RequestPool {
	return &RequestPool{
		p: sync.Pool{
			New: func() any {
				return NewRequest(WithClient(NewClient()))
			},
		},
		options: options,
	}
}

type SessionPool struct {
	*RequestPool
}

func NewSessionPool(options ...RequestOption) *SessionPool {
	return &SessionPool{
		&RequestPool{
			p: sync.Pool{
				New: func() any {
					return NewSession(WithClient(NewClient()))
				},
			},
			options: options,
		},
	}
}
