package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	proxyAddr = "127.0.0.1:5050"
	apiServer = "http://127.0.0.1:8080"
)

var authToken string

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Printf("The proxy is running on http://%s → %s", proxyAddr, apiServer)
	log.Fatal(http.ListenAndServe(proxyAddr, nil))
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
