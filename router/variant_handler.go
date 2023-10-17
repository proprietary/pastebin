package router

import (
	badger "github.com/dgraph-io/badger/v4"
	"net/http"
)

type VariantHandler struct {
	Db *badger.DB
}

func (rh VariantHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// prefer to assume the client is a browser
	// use "curl"-like handler only if certain
	var clientIsTextTerminal bool = !userAgentIsBrowser(req) && userAgentIsCurlLike(req)
	var handler http.Handler = nil
	if clientIsTextTerminal {
		handler = ttyClientHandlerSingleton(rh.Db)
	} else {
		handler = browserClientHandlerSingleton(rh.Db)
	}
	handler.ServeHTTP(w, req)
}
