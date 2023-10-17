package router

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/proprietary/pastebin/pastebin_record"
	"github.com/proprietary/pastebin/text_store"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net/http"
	"net/netip"
	"time"
)

type BrowserClientHandler struct {
	mux *http.ServeMux
	db  *badger.DB
}

func NewBrowserClientHandler(db *badger.DB) BrowserClientHandler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", handleCreate(db))
	mux.HandleFunc("/delete", handleDelete(db))
	mux.HandleFunc("/", handleRoot(db))
	return BrowserClientHandler{
		mux: mux,
		db:  db,
	}
}

func (c BrowserClientHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.mux.ServeHTTP(w, req)
}

func browserClientHandlerSingleton(db *badger.DB) *BrowserClientHandler {
	if browserClientHandler == nil {
		browserClientHandler = new(BrowserClientHandler)
		*browserClientHandler = NewBrowserClientHandler(db)
	}
	return browserClientHandler
}

func handleCreate(db *badger.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost && req.Method != http.MethodPut {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title:       "Wrong method",
					Description: "This is an error page responding to an incorrect HTTP method.",
				},
				ErrorMessage: "Wrong HTTP method; only POST or PUT are allowed at `/create`",
				StatusCode:   http.StatusMethodNotAllowed,
			})
			return
		}
		// get form fields
		paste := req.FormValue("paste")
		filename := req.FormValue("filename")
		expirationInput := req.FormValue("expiration")
		expiration := time.Now().Add(time.Hour * DEFAULT_EXPIRATION_HOURS)
		if len(expirationInput) > 0 {
			maybeExpiration, err := time.Parse("2006-01-02", req.FormValue("expiration"))
			if err != nil {
				log.Printf("Could not parse expiration from form %q; defaulting to now + 2y",
					expirationInput)
			} else {
				expiration = maybeExpiration
			}
		}
		ip := getClientIp(req)
		// validate paste
		record := pastebin_record.PastebinRecord{
			Body:        paste,
			TimeCreated: timestamppb.New(time.Now()),
			Filename:    &filename,
			Expiration:  timestamppb.New(expiration),
		}
		// TODO(zds): add mime type and syntax highlighting fields
		// record client IP ad the creator of this paste
		record.Creator = &pastebin_record.IPAddress{
			Ip: ip.AsSlice(),
		}
		if ip.Is4() {
			record.Creator.Version = pastebin_record.IPAddressVersion_V4
		} else {
			record.Creator.Version = pastebin_record.IPAddressVersion_V6
		}
		// persist paste to store
		slug, err := text_store.StoreNewPaste(db, &record)
		if err != nil {
			// TODO(zds): show error page
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// return result page
		w.Header().Set("Location", fmt.Sprintf("%s/%s", "", slug))
		w.WriteHeader(http.StatusFound)
		return
	}
}

func handleDelete(db *badger.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost && req.Method != http.MethodDelete {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title:       "Wrong method",
					Description: "This is an error page responding to an incorrect HTTP method.",
				},
				ErrorMessage: "Wrong HTTP method; only POST or DELETE are allowed on `/delete`",
				StatusCode:   http.StatusMethodNotAllowed,
			})
			return
		}
		slug := text_store.Slug(req.FormValue("slug"))
		if len(slug) == 0 {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title:       "Error: Bad request",
					Description: "This is an error page responding to a bad request missing a POST parameter.",
				},
				ErrorMessage: `Missing "slug" POST parameter in body`,
				StatusCode:   http.StatusBadRequest,
			})
			return
		}
		pastebinRecord, err := text_store.FindPaste(db, slug)
		if err != nil {
			log.Println(err)
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}
		owner, ok := netip.AddrFromSlice(pastebinRecord.GetCreator().GetIp())
		if !ok {
			log.Printf("Found record with invalid IP slug=%q", slug)
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}
		clientIp := getClientIp(req)
		if clientIp != owner {
			// TODO(zds): Show message that you do not own this paste
			log.Printf("attempt by %v to delete a paste (%q) owned by %v", clientIp, slug, owner)
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusFound)
			return
		}
		err = text_store.DeletePaste(db, slug)
		if err != nil {
			log.Printf(`Error in deleting paste %q: %v`, slug, err.Error())
		}
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
		return
	}
}

func handleRoot(db *badger.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/":
			page := CreatePage{
				Meta: Meta{
					Title:       "Create a new paste",
					Description: "Create a new text snippet saved as a link on the internet",
				},
			}
			err := OurViews.renderCreatePage(w, &page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		default:
			slug := text_store.Slug(req.URL.Path[1:])
			paste, err := text_store.FindPaste(db, slug)
			if err != nil {
				OurViews.renderErrorPage(w, &ErrorPage{
					Meta: Meta{
						Title:       "Paste not found",
						Description: "paste not found",
					},
					ErrorMessage: fmt.Sprintf("Paste \"%s\" not found or expired", slug),
					StatusCode:   http.StatusNotFound,
				})
				return
			}
			log.Println(paste.GetFilename(), paste.GetBody())
			page := ResultPage{
				Meta: Meta{
					Title:       "Paste found",
					Description: "Paste found",
				},
				Paste:    paste.GetBody(),
				Exp:      paste.GetExpiration().AsTime(),
				Filename: paste.GetFilename(),
				Slug:     string(slug),
			}
			err = OurViews.renderResultPage(w, &page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
}

var browserClientHandler *BrowserClientHandler = nil
