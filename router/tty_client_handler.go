package router

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/proprietary/pastebin/pastebin_record"
	"github.com/proprietary/pastebin/text_store"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			w.WriteHeader(http.StatusOK)
			return
		}
	case http.MethodPost, http.MethodPut:
		{
			var buf bytes.Buffer
			_, err := io.Copy(&buf, req.Body)
			if err != nil {
				http.Error(w, "Invalid text", http.StatusBadRequest)
				return
			}
			body := buf.String()
			if !utf8.ValidString(string(body)) {
				http.Error(w, "Text must be valid UTF-8", http.StatusBadRequest)
				return
			}
			clientIp := getClientIp(req)
			now := time.Now()
			expiration := now.Add(time.Hour*DEFAULT_EXPIRATION_HOURS)
			mimeType := req.Header.Get("Content-Type")
			record := pastebin_record.PastebinRecord{
				Body: body,
				Expiration: timestamppb.New(expiration),
				TimeCreated: timestamppb.New(now),
				MimeType: &mimeType,
			}
			record.Creator = &pastebin_record.IPAddress{
				Ip: clientIp.AsSlice(),
			}
			if clientIp.Is4() {
				record.Creator.Version = pastebin_record.IPAddressVersion_V4
			} else {
				record.Creator.Version = pastebin_record.IPAddressVersion_V6
			}
			slug, err := text_store.StoreNewPaste(c.db, &record)
			if err != nil {
				log.Println("Fail to save pastebin:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Write([]byte(slug))
			w.Header().Add("Content-Type", "text/plain")
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
