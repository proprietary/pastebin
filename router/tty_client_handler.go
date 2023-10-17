package router

import (
	badger "github.com/dgraph-io/badger/v4"
	"github.com/proprietary/pastebin/text_store"
	"io"
	"log"
	"net/http"
	"time"
)

type TtyClientHandler struct {
	db *badger.DB
}

func NewTtyClientHandler(db *badger.DB) TtyClientHandler {
	return TtyClientHandler{
		db: db,
	}
}

func (c TtyClientHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			slug := text_store.Slug(req.URL.Path[1:])
			text, err := text_store.LookupPastebin(c.db, slug)
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
			slug, err := text_store.SavePastebin(c.db, body, time.Now().Add(time.Hour*24*365*2))
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

func ttyClientHandlerSingleton(db *badger.DB) *TtyClientHandler {
	if ttyClientHandler == nil {
		ttyClientHandler = new(TtyClientHandler)
		*ttyClientHandler = NewTtyClientHandler(db)
	}
	return ttyClientHandler
}

var ttyClientHandler *TtyClientHandler = nil
