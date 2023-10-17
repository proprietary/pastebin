package router

import (
	"github.com/proprietary/pastebin/text_store"
	"github.com/proprietary/pastebin/pastebin_record"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	var handler http.Handler = nil
	if clientIsTextTerminal {
		handler = NewTtyClientHandler()
	} else {
		handler = NewBrowserClientHandler()
	}
	handler.ServeHTTP(w, req)
}

type TtyClientHandler struct {
}

func NewTtyClientHandler() TtyClientHandler {
	return TtyClientHandler{}
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

type BrowserClientHandler struct{
	mux *http.ServeMux
}

func NewBrowserClientHandler() BrowserClientHandler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost && req.Method != http.MethodPut {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title: "Wrong method",
					Description: "This is an error page responding to an incorrect HTTP method.",
				},
				ErrorMessage: "Wrong HTTP method; only POST or PUT are allowed at `/create`",
				StatusCode: http.StatusMethodNotAllowed,
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
			Body: paste,
			TimeCreated: timestamppb.New(time.Now()),
			Filename: &filename,
			Expiration: timestamppb.New(expiration),
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
		db := text_store.OpenDb()
		defer db.Close()
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
	})
	mux.HandleFunc("/delete", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost && req.Method != http.MethodDelete {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title: "Wrong method",
					Description: "This is an error page responding to an incorrect HTTP method.",
				},
				ErrorMessage: "Wrong HTTP method; only POST or DELETE are allowed on `/delete`",
				StatusCode: http.StatusMethodNotAllowed,
			})
			return
		}
		slug := text_store.Slug(req.FormValue("slug"))
		if len(slug) == 0 {
			OurViews.renderErrorPage(w, &ErrorPage{
				Meta: Meta{
					Title: "Error: Bad request",
					Description: "This is an error page responding to a bad request missing a POST parameter.",
				},
				ErrorMessage: `Missing "slug" POST parameter in body`,
				StatusCode: http.StatusBadRequest,
			})
			return
		}
		db := text_store.OpenDb()
		defer db.Close()
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
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
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
			paste, err := text_store.FindPaste(db, slug)
			if err != nil {
				OurViews.renderErrorPage(w, &ErrorPage{
					Meta: Meta{
						Title: "Paste not found",
						Description: "paste not found",
					},
					ErrorMessage: fmt.Sprintf("Paste \"%s\" not found or expired", slug),
					StatusCode: http.StatusNotFound,
				})
				return
			}
			log.Println(paste.GetFilename(), paste.GetBody())
			page := ResultPage{
				Meta: Meta{
					Title: "Paste found",
					Description: "Paste found",
				},
				Paste: paste.GetBody(),
				Exp: paste.GetExpiration().AsTime(),
				Filename: paste.GetFilename(),
			}
			err = OurViews.renderResultPage(w, &page)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	})
	return BrowserClientHandler{
		mux: mux,
	}
}

func (c BrowserClientHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.mux.ServeHTTP(w, req)
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
	// correctly receive client IP if this service is running behind Cloudflare
	cfClientIp := req.Header.Get("CF-Connecting-IP")
	if len(cfClientIp) > 0 {
		var isTrustworthy bool = false
		// make sure this is not spoofed
		if addr.IsPrivate() || addr.IsLoopback() {
			isTrustworthy = true
		} else {
			for _, prefix := range knownCloudflareNetworks {
				if prefix.Contains(addr) {
					isTrustworthy = true
				}
			}
		}
		if isTrustworthy {
			parsedCfClientIp, err := netip.ParseAddr(cfClientIp)
			if err != nil {
				// basically should never happen
				return addr
			}
			return parsedCfClientIp
		} else {
			// client tried to forge a Cloudflare request
			log.Printf(`Client with IP=%q tried to forge a Cloudflare request by setting the header "CF-Connecting-IP"`,
			addr.String())
			return addr
		}
	}
	return addr
}

const DEFAULT_EXPIRATION_HOURS = 24 * 365 * 2
