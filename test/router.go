package main

import (
	"context"
	"fmt"
	"net/http"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

func buildChain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](buildChain(f, m[1:cap(m)]...))
}

/*
func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/test", buildChain(endpoint, authentication, login)).Methods("GET")
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 7 * time.Second,
		Addr:         ":40500",
		Handler:      r,
	}
	server.ListenAndServe()
}
*/
//TOKEN asdf
const TOKEN = "TOKEN"

var authentication = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ... pre handler functionality
		fmt.Println("start auth")
		auth := r.Header.Get("Authentication")
		ctx := context.WithValue(r.Context(), TOKEN, auth)
		fmt.Println("end auth")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

var login = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login")
		token := r.Context().Value(TOKEN)
		fmt.Println(token)
		next.ServeHTTP(w, r)
	}
}

func endpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("server content")
	fmt.Fprintf(w, "CONNECTED")
}
