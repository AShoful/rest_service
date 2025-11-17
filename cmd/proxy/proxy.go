package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	proxyAddr = "127.0.0.1:5050"
	apiServer = "http://127.0.0.1:8080"
)

var authToken string

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxyHandler)

	srv := &http.Server{
		Addr:    proxyAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("Proxy is running on http://%s → %s", proxyAddr, apiServer)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down proxy server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Proxy stopped gracefully.")
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/auth/sign-in" && r.Method == http.MethodPost {
		handleSignIn(w, r)
		return
	}
	handleProxy(w, r)
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	resp, err := http.Post(apiServer+"/auth/sign-in", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "API request error /auth/sign-in", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading API response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var data map[string]interface{}
		if err := json.Unmarshal(respBody, &data); err == nil {
			if token, ok := data["token"].(string); ok {
				authToken = token
				log.Printf("Token received and saved: %s", authToken)
			}
		}
	}

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	targetURL := apiServer + r.URL.Path

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error requesting target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
