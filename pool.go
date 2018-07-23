package gofastcgi

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)

type factory func(host string, port int, keepAlive bool) (*Client, error)

type ClientPool struct {
	sync.Mutex
	poolOptions PoolOptions
	pool        chan *Client
	closed      bool    // pool state
	factory     factory // create client factory
}

type PoolOptions struct {
	host        string
	port        int
	maxOpen     int // max client
	numOpen     int // now client count
	startOpen   int // start client
	maxLifetime time.Duration
}

// New Client Pool
func NewClientPool(poolOptions PoolOptions) (*ClientPool, error) {
	if poolOptions.host == "" || poolOptions.port == 0 {
		return nil, ErrInvalidConfig
	}
	if poolOptions.maxOpen <= 0 || poolOptions.startOpen > poolOptions.maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &ClientPool{
		poolOptions: poolOptions,
		factory:     NewClient,
		pool:        make(chan *Client, poolOptions.maxOpen),
	}

	for i := 0; i < poolOptions.startOpen; i++ {
		client, err := p.factory(poolOptions.host, poolOptions.port, true)
		if err != nil {
			continue
		}
		p.poolOptions.numOpen++
		p.pool <- client
	}
	return p, nil
}

func (p *ClientPool) Acquire() (*Client, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	for {
		client, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		// if set maxLifetime, and (client activityTime + maxLifetime) < now，then the client is expire
		if p.poolOptions.maxLifetime > 0 && client.GetActiveTime().Add(p.poolOptions.maxLifetime).Before(time.Now()) {
			p.Close(client)
			continue
		}
		client.SetActiveTime()
		return client, nil
	}
}

func (p *ClientPool) getOrCreate() (*Client, error) {
	select {
	case closer := <-p.pool:
		return closer, nil
	default:
	}
	p.Lock()
	if p.poolOptions.numOpen >= p.poolOptions.maxOpen {
		closer := <-p.pool
		p.Unlock()
		return closer, nil
	}
	// create client
	client, err := p.factory(p.poolOptions.host, p.poolOptions.port, true)
	if err != nil {
		p.Unlock()
		return nil, err
	}
	p.poolOptions.numOpen++
	p.Unlock()
	return client, nil
}

// release single client to pool
func (p *ClientPool) Release(client *Client) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	p.pool <- client
	p.Unlock()
	return nil
}

// close single client
func (p *ClientPool) Close(client *Client) error {
	p.Lock()
	client.Close()
	p.poolOptions.numOpen--
	p.Unlock()
	return nil
}

// close pool, release clients
func (p *ClientPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	close(p.pool)
	for client := range p.pool {
		client.Close()
		p.poolOptions.numOpen--
	}
	p.closed = true
	p.Unlock()
	return nil
}
