package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	. "github.com/strengthening/goghostex"
)

const (
	API_KEY       = ""
	API_SECRETKEY = ""
)

/**
 *
 * The func of market unit test step is:
 * 1. Get the BNBBTC ticker
 * 2. Get the BNBBTC depth
 * 3. Get the BNBBTC 1d 1m kline
 *
 **/

func TestSpot_GetTicker(t *testing.T) {
	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1090"}, nil
				},
			},
		},
		ApiKey:        API_KEY,
		ApiSecretKey:  API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	b := New(config)
	if ticker, resp, err := b.Spot.GetTicker(CurrencyPair{
		CurrencyTarget: Currency{"btc", ""},
		CurrencyBasis:  Currency{"usdt", ""},
	}); err != nil {
		t.Error(err)
		return
	} else {
		body, _ := json.Marshal(*ticker)
		fmt.Println(string(body))
		fmt.Println(string(resp))
	}
}

func TestSpot_GetDepth(t *testing.T) {

	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1090"}, nil
				},
			},
		},
		ApiKey:        API_KEY,
		ApiSecretKey:  API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	b := New(config)
	if depth, _, err := b.Spot.GetDepth(
		50,
		CurrencyPair{BTC, USDT}); err != nil {
		t.Error(err)
		return
	} else {
		body, _ := json.Marshal(*depth)
		fmt.Println(string(body))
		//fmt.Println(string(resp))
	}
}

func TestSpot_GetKlineRecords(t *testing.T) {
	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1090"}, nil
				},
			},
		},
		ApiKey:        API_KEY,
		ApiSecretKey:  API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	b := New(config)
	if klines, resp, err := b.Spot.GetKlineRecords(
		CurrencyPair{
			CurrencyTarget: Currency{"btc", ""},
			CurrencyBasis:  Currency{"usdt", ""},
		},
		KLINE_PERIOD_1MIN,
		50,
		int(time.Now().Add(-2*24*time.Hour).UnixNano()),
	); err != nil {
		t.Error(err)
		return
	} else {
		body, _ := json.Marshal(klines)
		fmt.Println(string(body))
		fmt.Println(string(resp))
	}
}

/*
 *
 * The func of order unit test step is:
 * 1. Get BNBBTC ticker.
 * 2. Get the account, and find have the enough crypto.
 * 2. Order the Limit Sell/Buy without deal.
 * 3. Get the unfinished orders info, and find the order in step 1.
 * 4. Get the order info.
 * 5. Cancel the Limit Order
 *
 */

func TestSpot_GetAccount(t *testing.T) {

	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1090"}, nil
				},
			},
		},
		ApiKey:        API_KEY,
		ApiSecretKey:  API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	b := New(config)
	if account, _, err := b.Spot.GetAccount(); err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println(*account)
	}
}

func TestSpot_LimitSell(t *testing.T) {
	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1090"}, nil
				},
			},
		},
		ApiKey:        API_KEY,
		ApiSecretKey:  API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	b := New(config)
	order := &Order{
		Currency: CurrencyPair{BNB, BTC},
		Price:    0.0030061,
		Amount:   1,
		Side:     SELL,
	}

	if resp, err := b.Spot.LimitSell(order); err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("LIMIT SELL ~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(string(resp))
	}

	for i := 0; i < 3; i++ {
		if resp, err := b.Spot.GetOneOrder(order); err != nil {
			t.Error(err)
			return
		} else {
			fmt.Println("GET ORDER ~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Println(string(resp))
		}

		body, _ := json.Marshal(*order)
		fmt.Println(string(body))
	}

	if _, resp, err := b.Spot.GetUnFinishOrders(
		CurrencyPair{BNB, BTC},
	); err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("GET UNFINISHED ORDER ~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(string(resp))
	}

	if resp, err := b.Spot.CancelOrder(order); err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("CANCEL ORDER ~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(string(resp))
	}

}
