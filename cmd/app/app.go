package main

import (
	"ezp/pkg/cfg"
	"ezp/pkg/git"
	"ezp/pkg/srv"
	"flag"
	"fmt"
	"log"
	"net/netip"
)

func main() {
	host := flag.String("host", "127.0.0.1", "--host 127.0.0.1")
	port := flag.String("port", "2442", "--port 2442")
	conf := flag.String("conf", "", "--conf configs/prod.json")
	flag.Parse()

	addr, err := netip.ParseAddrPort(fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("bad addr: %v", err)
	}

	if *conf == "" {
		log.Fatalf("config file path not passed")
	} else if err := cfg.Load(*conf); err != nil {
		log.Fatalf("unable load config: %v", err)
	}

	go git.ReposWatcher()

	s := srv.New(addr)
	s.Run()
}
