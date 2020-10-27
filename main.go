package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"strings"
)

var httpPort = flag.String("http", ":8080", "http port")
var httpsPort = flag.String("https", ":8090", "https port")
var httpRedirect = flag.Bool("http-redirect", false, "if true redirect http to https")
var crtFile = flag.String("crt", "server.crt", "certificate public key")
var keyFile = flag.String("key", "server.key", "certificate private key")
var files = flag.String("files", "static", "html files location")

// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
// snap install --classic certbot
// certbot certonly --standalone
// certbot renew
// ln -s /etc/letsencrypt/live/example.com/cert.pem server.crt
// ln -s /etc/letsencrypt/live/example.com/privkey.pem server.key
// /root/webserver -http=:80 -https=:443 -crt=/etc/letsencrypt/live/example.com/cert.pem -key=/etc/letsencrypt/live/example.com/privkey.pem -files=/root/static
// curl --verbose --insecure -L http://localhost:8080
// curl --verbose --insecure https://localhost:8090
func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go func() {
		// allow using self signed certificates
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if err := http.ListenAndServeTLS(*httpsPort, *crtFile, *keyFile, http.FileServer(http.Dir(*files))); err != nil {
			log.Println(err)
		}
	}()
	h := http.FileServer(http.Dir(*files))
	if *httpRedirect {
		h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+strings.Split(r.Host, ":")[0]+*httpsPort+r.RequestURI, http.StatusMovedPermanently)
		})
	}
	if err := http.ListenAndServe(*httpPort, h); err != nil {
		log.Fatalln(err)
	}
}
