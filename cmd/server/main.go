package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	port        = flag.Uint("port", 8955, "port to listend or connect to rpc calls")
	isServer    = flag.Bool("server", false, "activates server mode")
	useJson     = flag.Bool("json", false, "whether it should use json-rpc")
	serverSleep = flag.Duration("server.sleep", 0, "time for the server to sleep on requests")
)

func handleSignals() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	log.Println("signal recieved ", sig)
}

func must(err error) {
	if err == nil {
		return
	}

	log.Panic(err)
}

func runServer() {
	server := &Server{
		UseHttp: *http,
		UseJson: *useJson,
		Sleep:   *serverSleep,
		Port:    *port,
	}
	defer server.Close()

	go func() {
		handleSignals()
		server.Close()
		os.Exit(0)
	}()

	must(server.Start())

	return
}

func runClient() {
	client := &Client{
		UseHttp: *http,
		useJson: *useJson,
		Port:    *port,
	}
	defer client.Close()

	must(client.Init())

	response, err := client.Execute("ciro")
	must(err)

	log.Println(response)
}

func main() {
	flag.Parse()

	if *isServer {
		log.Println("starting server")
		log.Printf("will listen on port %d\n", *port)

		runServer()

		return
	}

	log.Println("starting client")
	log.Printf("will connect to port %d\n", *port)

	runClient()

	return
}
