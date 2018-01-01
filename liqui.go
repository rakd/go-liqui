package liqui

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	// APIBaseURL liqui API endpoint
	APIBaseURL        = "https://api.liqui.io/api/3"
	APITradingBaseURL = "https://api.liqui.io/tapi"
)

func init() {
	log.SetFlags(log.Lshortfile)

}

// New returns an instantiated liqui struct
func New(apiKey, apiSecret string) *Liqui {
	client := NewClient(apiKey, apiSecret)
	return &Liqui{client}
}

// NewWithCustomTimeout returns an instantiated liqui struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Liqui {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Liqui{client}
}

// Liqui represent a liqui client
type Liqui struct {
	client *Client
}

// SetDebug represent a liqui client
func (b *Liqui) SetDebug(d bool) {
	b.client.debug = d
}

// GetBalance ..
func (b *Liqui) GetBalance() (res AccountInfoResponse, r []byte, err error) {
	payload := map[string]string{}
	payload["method"] = "getInfo"
	r, err = b.client.do("POST", APITradingBaseURL, payload, true)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &res); err != nil {
		er := ErrorResponse{}

		if err = json.Unmarshal(r, &er); err == nil {
			err = fmt.Errorf("Error: success=%d, code=%d, %s", er.Success, er.Code, er.Error)
			return
		}
		log.Printf("body=%s", string(r))
		return
	}

	// ok
	return
}

// GetInfo ..
func (b *Liqui) GetTicker() (res TickerResponse, r []byte, err error) {
	r, err = b.client.do("GET", fmt.Sprintf("info"), nil, false)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &res); err != nil {

		er := ErrorResponse{}
		if err = json.Unmarshal(r, &er); err == nil {
			err = fmt.Errorf("Error: success=%v, %s", er.Success, er.Error)
			return
		}
		log.Printf("body=%s", string(r))
		return
	}

	// ok
	return
}
