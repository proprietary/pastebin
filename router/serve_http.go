package router

import (
	"github.com/proprietary/pastebin/text_store"
	"io"
	"log"
	"net/http"
	"time"
)

type RootHandler struct {
}

func (_ RootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// prefer to assume the client is a browser
	// use "curl"-like handler only if certain
	var clientIsTextTerminal bool = !userAgentIsBrowser(req) && userAgentIsCurlLike(req)
	if clientIsTextTerminal {
		var handler TtyClientHandler
		handler.ServeHTTP(w, req)
	} else {
		var handler BrowserClientHandler
		handler.ServeHTTP(w, req)
	}
}

type TtyClientHandler struct {
}

func (_ TtyClientHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
			slug, err := text_store.SavePastebin(db, body, time.Now().Add(time.Hour*24*365*2))
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

type BrowserClientHandler struct{}

func (_ BrowserClientHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Not Implemented"))
			return
		}
	case http.MethodPost:
		{
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Not Implemented"))
			return
		}
	default:
		{
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Not Implemented"))
			return
		}
	}
}
