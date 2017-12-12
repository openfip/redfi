package redfi

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pool "gopkg.in/fatih/pool.v2"
)

type Proxy struct {
	server   string
	plan     *Plan
	bind     string
	connPool pool.Pool
}

func factory(server string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return net.Dial("tcp", server)

	}
}

func New(planPath, server, addr string) (*Proxy, error) {
	p, err := pool.NewChannelPool(5, 30, factory(server))
	if err != nil {
		return nil, err
	}

	// parse the failures plan
	plan, err := Parse(planPath)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		server:   server,
		connPool: p,
		plan:     plan,
		bind:     addr,
	}, nil
}

func (p *Proxy) Start() error {
	fmt.Println("RedFI: a Redis Fault-Injection Proxy")
	ln, err := net.Listen("tcp", p.bind)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go p.handle(conn)
	}
}

func (p *Proxy) handle(conn net.Conn) {
	var wg sync.WaitGroup

	targetConn, err := p.connPool.Get()
	if err != nil {
		log.Fatal("failed to get a connection from connPool")
	}

	wg.Add(2)
	go func() {
		p.faulter(targetConn, conn)
		wg.Done()
	}()
	go func() {
		p.pipe(conn, targetConn)
		wg.Done()
	}()
	wg.Wait()

	log.Println("Close connection", conn.Close())
}

func (p *Proxy) pipe(dst, src net.Conn) {
	buf := make([]byte, 32<<10)

	for {
		n, err := src.Read(buf)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		// @TODO(kl): check if written is less than what's in buf
		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}

	}
}

func (p *Proxy) faulter(dst, src net.Conn) {
	buf := make([]byte, 32<<10)

	for {
		n, err := src.Read(buf)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		rule := p.plan.SelectRule(src.RemoteAddr().String(), buf)

		if rule != nil {
			if rule.Delay > 0 {
				time.Sleep(time.Duration(rule.Delay) * time.Millisecond)
			}

			if rule.Drop {
				err = src.Close()
				if err != nil {
					log.Println("encountered error while closing srcConn", err)
				}
				break
			}

			if rule.ReturnEmpty {
				_, err = dst.Write([]byte("$-1\r\n"))
				if err != nil {
					log.Println(err)
				}
				continue
			}
		}

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}

	}
}
