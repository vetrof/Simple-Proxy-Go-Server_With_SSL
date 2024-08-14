package main

import (
	"io"
	"log"
	"net/http"
)

func handleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received request: %s %s", req.Method, req.URL.String())
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf("Sent response: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
}

func main() {
	http.HandleFunc("/", handleRequestAndRedirect)

	// Запуск HTTP сервера
	go func() {
		log.Println("Starting HTTP proxy server on :8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	// Запуск HTTPS сервера
	log.Println("Starting HTTPS proxy server on :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil))

	// На продакшене
	//log.Println("Starting HTTPS proxy server on :8443")
	//log.Fatal(http.ListenAndServeTLS(":8443", "/etc/letsencrypt/live/yourdomain/fullchain.pem", "/etc/letsencrypt/live/yourdomain/privkey.pem", nil))

}
