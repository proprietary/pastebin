package main

import (
	"log"
	"net/http"
	// "html/template"
	// "github.com/proprietary/pastebin/text_store"
	"github.com/proprietary/pastebin/router"
	// pb "github.com/proprietary/pastebin/pastebin_record"
)

func main() {
	var handler router.RootHandler
	err := http.ListenAndServe(":50999", handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
