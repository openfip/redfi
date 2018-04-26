package redfi

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	pool "gopkg.in/fatih/pool.v2"
)

type Proxy struct {
	server   string
	plan     *Plan
	addr     string
	connPool pool.Pool
	api      *API
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

	plan := NewPlan()
	if len(planPath) > 0 {
		// parse the failures plan
		plan, err = Parse(planPath)
		if err != nil {
			return nil, err
		}
	}

	return &Proxy{
		server:   server,
		connPool: p,
		plan:     plan,
		addr:     addr,
		api:      NewAPI(plan),
	}, nil
}

func (p *Proxy) startAPI() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// RESTy routes for "rules" resource
	r.Route("/rules", func(r chi.Router) {
		r.Get("/", p.api.listRules) // GET /rules

		r.Post("/", p.api.createRule) // POST /rules

		// Subrouters:
		r.Route("/{ruleName}", func(r chi.Router) {
			r.Get("/", p.api.getRule)       // GET /rules/drop_20
			r.Delete("/", p.api.deleteRule) // DELETE /rules/drop_20
		})
	})

	// @TODO(kl): get api port from cli
	fmt.Println("API is listening on :8081")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Proxy) Start() error {
	ln, err := net.Listen("tcp", p.addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nRedFI is listening on", p.addr)
	fmt.Println("Don't forget to point your client to that address.")

	ctr, err := newController(p.plan)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		fmt.Println("\nRedFI Controller is listening on", configAddr)
		err := ctr.Start()
		if err != nil {
			log.Fatal("encountered err while starting controller", err)
		}
	}()

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
