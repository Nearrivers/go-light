package request

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/nearrivers/go-light/app_errors"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

// Effectue une requête HTTP qui ne nécessite pas de corps dans la requête
func NewHueBodylessRequest(method HttpMethod, uri string) (*http.Response, error) {
	var bridgeIp string = os.Getenv("BRIDGE_IP")
	client := http.Client{}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(string(method), fmt.Sprintf("https://%s/%s", bridgeIp, uri), nil)

	if err != nil {
		return &http.Response{}, app_errors.RuntimeError{Err: err}
	}

	req.Header.Set("hue-application-key", os.Getenv("HUE_BRIDGE_USERNAME"))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "*/*")

	return client.Do(req)
}

// Effectue une requête HTTP qui nécessite un corps dans la requête
func NewHueRequestWithBody(method HttpMethod, uri string, body io.Reader) (*http.Response, error) {
	var bridgeIp string = os.Getenv("BRIDGE_IP")
	client := http.Client{}
	req, err := http.NewRequest(string(method), fmt.Sprintf("https://%s/%s", bridgeIp, uri), body)

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("hue-application-key", os.Getenv("HUE_BRIDGE_USERNAME"))
	req.Header.Set("Content-type", "application/json")

	return client.Do(req)
}
