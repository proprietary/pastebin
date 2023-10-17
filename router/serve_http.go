package router

import (
	"github.com/proprietary/pastebin/text_store"
	"io"
	"fmt"
	"log"
	"net/http"
	"net/netip"
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
			switch req.URL.Path {
			case "/":
				page := CreatePage{
					Meta: Meta{
						Title: "Create a new paste",
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
				db := text_store.OpenDb()
				defer db.Close()
				text, err := text_store.LookupPastebin(db, slug)
				if err != nil {
					// TODO(zds): show not found page
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				page := ResultPage{
					Meta: Meta{
						Title: "Paste found",
						Description: "Paste found",
					},
					Paste: text,
					Exp: time.Now(),
				}
				err = OurViews.renderResultPage(w, &page)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
	case http.MethodPost, http.MethodPut:
		{
			// get form fields
			paste := req.FormValue("paste")
			// filename := req.FormValue("filename")
			// ip := getClientIp(req)
			// validate paste
			// persist paste to store
			db := text_store.OpenDb()
			defer db.Close()
			expiration, err := time.Parse("2006-01-02", req.FormValue("expiration"))
			if err != nil {
				// TODO(zds): render error page
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slug, err := text_store.SavePastebin(db, []byte(paste), expiration)
			if err != nil {
				// TODO(zds): show error page
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// render result page
			w.Header().Set("Location", fmt.Sprintf("%s/%s", "", slug))
			w.WriteHeader(http.StatusFound)
			return
		}
	default:
		{
			// TODO(zds): show error page
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("Not Implemented"))
			return
		}
	}
}

func knownCloudflare() []netip.Prefix {
	// See: https://www.cloudflare.com/ips/
	cloudflareIps := [...]string{
		"173.245.48.0/20",
		"103.21.244.0/22",
		"103.22.200.0/22",
		"103.31.4.0/22",
		"141.101.64.0/18",
		"108.162.192.0/18",
		"190.93.240.0/20",
		"188.114.96.0/20",
		"197.234.240.0/22",
		"198.41.128.0/17",
		"162.158.0.0/15",
		"104.16.0.0/13",
		"104.24.0.0/14",
		"172.64.0.0/13",
		"131.0.72.0/22",
		"2400:cb00::/32",
		"2606:4700::/32",
		"2803:f800::/32",
		"2405:b500::/32",
		"2405:8100::/32",
		"2a06:98c0::/29",
		"2c0f:f248::/32",
	}
	dst := make([]netip.Prefix, len(cloudflareIps))
	for i, ipString := range cloudflareIps {
		dst[i] = netip.MustParsePrefix(ipString)
	}
	return dst
}

var knownCloudflareNetworks []netip.Prefix

func init() {
	knownCloudflareNetworks = knownCloudflare()
}


func getClientIp(req *http.Request) netip.Addr {
	// normal client
	addrPort := netip.MustParseAddrPort(req.RemoteAddr)
	addr := addrPort.Addr()
	cfClientIp := req.Header.Get("CF-Connecting-IP")
	if len(cfClientIp) > 0 {
		// proxied by Cloudflare
		// make sure this is spoofed
		for _, prefix := range knownCloudflareNetworks {
			if prefix.Contains(addr) {
				addr, err := netip.ParseAddr(cfClientIp)
				if err != nil {
					return addr
				}
				return addr
			}
		}
		// client tried to forge a Cloudflare request
		log.Printf(`Client with IP=%q tried to forge a Cloudflare request by setting the header "CF-Connecting-IP"`,
			addr.String())
		return addr
	}
	return addr
}
