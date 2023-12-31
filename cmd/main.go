package main

import (
	"flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/proprietary/pastebin/router"
	"github.com/proprietary/pastebin/text_store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var port *uint = flag.Uint("port", 50999, "port to run HTTP server")

func main() {
	flag.Parse()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	var db *badger.DB = text_store.OpenDb()
	defer db.Close()
	go func() {
		<-signals
		log.Println("Shutting down...")
		db.Close()
		os.Exit(0)
	}()
	var handler router.VariantHandler = router.VariantHandler{
		Db: db,
	}
	log.Printf("Listening on port %d...", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
