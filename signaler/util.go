package signaler

import (
	"log"
	"net/http"
)

func httpf(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			log.Println("err:", err)
			http.Error(w, "server error", 500)
		}
	}
}
