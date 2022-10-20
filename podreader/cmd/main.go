package main

import (
	"io"
	"log"
	"net/http"
	"podreader/formatter"
	"podreader/readers"
	"strings"
)

func getStatus(w http.ResponseWriter, r *http.Request, reader *readers.Reader) {
	// Remove all "/" from the url in an attempt to get the
	// desired namesapce. There shouldn't be any "/" other than
	// the preceding and ending "/" so any request with more
	// than those 2 uses of "/" are invalid anyway and will
	// return a 404.
	namespace := strings.Replace(r.URL.Path, "/", "", -1)
	log.Printf("Recevied request for namespace: %s", namespace)
	status, code := reader.GetStatus(namespace)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// If we're returning an error, just write the error directly
	if code != 200 {
		w.WriteHeader(code)
		io.WriteString(w, status)
		return
	}

	formattedStatus, err := formatter.FormatStatus(status)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, formattedStatus)
		return
	}
	io.WriteString(w, formattedStatus)
}

func main() {
	reader := readers.NewReader()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		getStatus(w, r, reader)
	})

	log.Print("Pod Reader starting...")
	if err := http.ListenAndServeTLS(":42069", "tls/tls.crt", "tls/tls.key", nil); err != nil {
		log.Fatalf("Got error: %v\n", err)
	}
}
