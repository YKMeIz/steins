package main

import (
	"github.com/YKMeIz/steins/internal/docker"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	virtualHosts, err := docker.GetVirtualHosts()
	if err != nil {
		log.Panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// redirect www to non-www
		if strings.HasPrefix(r.Host, "www") {
			target := "https://" + r.Host[4:] + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusFound)
			return
		}

		ip, ok := virtualHosts.Load(strings.Split(r.Host, ":")[0])

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		director := func(req *http.Request) {
			req.Header = r.Header
			req.URL = r.URL
			req.URL.Scheme = "http"
			req.URL.Host = ip.(string)
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", mux))
}
