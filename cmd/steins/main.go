package main

import (
	"github.com/YKMeIz/steins/internal/docker"
	"github.com/YKMeIz/steins/internal/podman"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

var mode string

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func init() {
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
}

func main() {
	virtualHosts, err := podman.GetVirtualHosts()
	if err != nil {
		log.Println("error from podman:", err.Error())
		docker.NetworkInit()
		virtualHosts, err = docker.GetVirtualHosts()
		if err != nil {
			panic(err)
		}
		mode = "docker"
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

		var (
			ip string
			ok bool
		)

		if mode == "docker" {
			ip, ok = virtualHosts.LoadWithReacquire(strings.Split(r.Host, ":")[0])
		} else {
			ip, ok = podman.LoadWithReacquire(virtualHosts, strings.Split(r.Host, ":")[0])
		}

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		director := func(req *http.Request) {
			req.Header = r.Header
			req.URL = r.URL
			req.URL.Scheme = "http"
			req.URL.Host = ip
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", mux))
}
