package agw

import (
	"context"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
)

type contextKey string

const (
	//ContextKeyBody for json body parse
	ContextKeyBody = contextKey("body")
)

func EnableCORS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func ParseJSONBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sj, err := simplejson.NewFromReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyBody, sj)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
