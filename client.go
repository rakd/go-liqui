package liqui

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Client ...
type Client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	throttle    <-chan time.Time
	httpTimeout time.Duration
	debug       bool
}

var (
	// Technically 6 req/s allowed, but we're being nice / playing it safe.
	reqInterval = 200 * time.Millisecond
)

// NewClient return a new LIQUI HTTP client
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		httpClient:  &http.Client{},
		throttle:    time.Tick(reqInterval),
		httpTimeout: 30 * time.Second,
		debug:       false,
	}
}

// NewClientWithCustomTimeout returns a new LIQUI HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Client {
	return &Client{
		apiKey:      apiKey,
		apiSecret:   apiSecret,
		httpClient:  &http.Client{},
		throttle:    time.Tick(reqInterval),
		httpTimeout: timeout,
		debug:       false,
	}
}

func (c Client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c Client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {

		if c.debug {
			c.dumpRequest(req)
		}

		resp, err := c.httpClient.Do(req)

		if c.debug {
			c.dumpResponse(resp)
		}

		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from LIQUI API")
	}
}

/*
func generateHmacSha256(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
*/
func generateHmacSha512(text, key string) string {
	hasher := hmac.New(sha512.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *Client) makeReq(method, resource string, payload map[string]string, authNeeded bool, respCh chan<- []byte, errCh chan<- error) {
	body := []byte{}
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(resource, "http") {
		rawurl = resource
	} else {
		rawurl = fmt.Sprintf("%s/%s", APIBaseURL, resource)
	}

	formValues := url.Values{}
	nonce := fmt.Sprintf("%d", time.Now().Unix())
	for key, value := range payload {
		formValues.Add(key, value)
		//formValues.Set(key, value)
	}
	formValues.Add("nonce", nonce)
	//formValues.Set("nonce", nonce)
	formData := formValues.Encode()

	req, err := http.NewRequest(method, rawurl, strings.NewReader(formData))
	if err != nil {
		respCh <- body
		//errCh <- errors.New("You need to set API Key and API Secret to call this method")
		errCh <- err
		return
	}
	log.Printf("rawurl:%s", rawurl)

	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			respCh <- body
			errCh <- errors.New("You need to set API Key and API Secret to call this method")
			return
		}
		log.Printf("key:%s", c.apiKey)
		log.Printf("secret:%s", c.apiSecret[0:5]+"...")
		log.Printf("secret:%s", c.apiSecret)
		log.Printf("nonce:%s", nonce)
		//log.Printf("rawurl:%s", rawurl)
		log.Printf("formData:%s", formData)
		//sig := generateHmacSha512(rawurl+formData, c.apiSecret)
		sig := generateHmacSha512(formData, c.apiSecret)
		log.Printf("sig:%s", sig)
		req.Header.Add("Key", c.apiKey)
		//req.Header.Add("ACCESS-NONCE", nonce)
		req.Header.Add("Sign", sig)
	}
	//req.Header.Add("User-Agent", userAgent)

	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Header.Add("Accept", "application/json")

	resp, err := c.doTimeoutRequest(connectTimer, req)

	if err != nil {
		log.Printf("123err=%v", err)
		respCh <- body
		errCh <- err
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}
	//log.Printf("body:%s", string(body))
	if resp.StatusCode != 200 {
		respCh <- body
		errCh <- errors.New(resp.Status)
		return
	}

	respCh <- body
	errCh <- nil
	close(respCh)
	close(errCh)
}

// do prepare and process HTTP request to LIQUI API
func (c *Client) do(method, resource string, payload map[string]string, authNeeded bool) (response []byte, err error) {

	respCh := make(chan []byte)
	errCh := make(chan error)
	<-c.throttle
	go c.makeReq(method, resource, payload, authNeeded, respCh, errCh)
	response = <-respCh
	err = <-errCh
	return
}
