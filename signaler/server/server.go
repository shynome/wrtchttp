package main

import (
	"flag"
	"net/http"

	"github.com/lainio/err2/try"
	"github.com/shynome/wrtchttp/signaler"
)

var args struct {
	Addr string
}

func init() {
	flag.StringVar(&args.Addr, "addr", ":1338", "listen addr")
}

func main() {

	flag.Parse()
	srv := signaler.New()
	// http://127.0.0.1:1338/signaler
	http.Handle("/signaler", srv)

	try.To(
		http.ListenAndServe(args.Addr, nil))
}
