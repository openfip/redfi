package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/openfip/redfi"
)

var (
	planPath = flag.String("plan", "", "path to the plan.json")
	server   = flag.String("redis", "127.0.0.1:6379", "address to the redis server, to proxy requests to")
	addr     = flag.String("addr", "127.0.0.1:8083", "address for the proxy to listen on")
)

func main() {
	flag.Parse()

	if len(*planPath) == 0 {
		flag.Usage()
		fmt.Println("\n[ERROR]: -plan can't be empty")
		os.Exit(1)
	}

	proxy, err := redfi.New(*planPath, *server, *addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	proxy.Start()
}
