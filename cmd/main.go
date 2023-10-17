package main

import (
	"log"
	"net/http"
	// "html/template"
	"github.com/proprietary/pastebin/router"
	"github.com/proprietary/pastebin/text_store"
	// pb "github.com/proprietary/pastebin/pastebin_record"
	"flag"
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
	"os"
	"os/signal"
	"syscall"
)

var port *uint = flag.Uint("port", 50999, "port to run HTTP server")

func main() {
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
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
