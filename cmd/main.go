package main

import (
	"io"
	"log"
	"net/http"
	// "html/template"
	"github.com/proprietary/pastebin/text_store"
)

func main() {
	var handler MyHandler
	err := http.ListenAndServe(":50999", handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type MyHandler struct {
}

func (_ MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var ua = req.Header.Get("user-agent")
	if len(ua) == 0 {
		http.Error(w, "No user agent", http.StatusBadRequest)
		return
	}
	switch req.Method {
	case http.MethodGet:
		{
			slug := text_store.Slug(req.URL.Path[1:])
			db := text_store.OpenDb()
			defer db.Close()
			text, err := text_store.LookupPastebin(db, slug)
			if err != nil {
				log.Println(err)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			w.Write([]byte(text))
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
	case http.MethodPost, http.MethodPut:
		{
			body, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(w, "Invalid text posted", http.StatusBadRequest)
				return
			}
			log.Println(string(body))
			db := text_store.OpenDb()
			defer db.Close()
			slug, err := text_store.SavePastebin(db, body)
			if err != nil {
				log.Println("Fail to save pastebin:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(slug))
			w.WriteHeader(http.StatusOK)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}
