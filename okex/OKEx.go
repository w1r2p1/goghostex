package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/strengthening/goghostex"
)

const (
	/*
	  http headers
	*/
	OK_ACCESS_KEY        = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN       = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP  = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE = "OK-ACCESS-PASSPHRASE"

	/**
	  paging params
	*/
	OK_FROM  = "OK-FROM"
	OK_TO    = "OK-TO"
	OK_LIMIT = "OK-LIMIT"

	CONTENT_TYPE = "Content-Type"
	ACCEPT       = "Accept"
	COOKIE       = "Cookie"
	LOCALE       = "locale="

	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"

	/*
	  i18n: internationalization
	*/
	ENGLISH            = "en_US"
	SIMPLIFIED_CHINESE = "zh_CN"
	//zh_TW || zh_HK
	TRADITIONAL_CHINESE = "zh_HK"

	/*
	  http methods
	*/
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"

	/*
	 others
	*/
	ResultDataJsonString = "resultDataJsonString"
	ResultPageJsonString = "resultPageJsonString"

	BTC_USD_SWAP = "BTC-USD-SWAP"
	LTC_USD_SWAP = "LTC-USD-SWAP"
	ETH_USD_SWAP = "ETH-USD-SWAP"
	ETC_USD_SWAP = "ETC-USD-SWAP"
	BCH_USD_SWAP = "BCH-USD-SWAP"
	BSV_USD_SWAP = "BSV-USD-SWAP"
	EOS_USD_SWAP = "EOS-USD-SWAP"
	XRP_USD_SWAP = "XRP-USD-SWAP"

	/*Rest Endpoint*/
	Endpoint              = "https://www.okex.com"
	GET_ACCOUNTS          = "/api/swap/v3/accounts"
	PLACE_ORDER           = "/api/swap/v3/order"
	CANCEL_ORDER          = "/api/swap/v3/cancel_order/%s/%s"
	GET_ORDER             = "/api/swap/v3/orders/%s/%s"
	GET_POSITION          = "/api/swap/v3/%s/position"
	GET_DEPTH             = "/api/swap/v3/instruments/%s/depth?size=%d"
	GET_TICKER            = "/api/swap/v3/instruments/%s/ticker"
	GET_UNFINISHED_ORDERS = "/api/swap/v3/orders/%s?status=%d&from=%d&limit=%d"
)

type OKEx struct {
	config     *APIConfig
	OKExSpot   *OKExSpot
	OKExFuture *OKExFuture
	//OKExSwap   *OKExSwap
	//OKExWallet *OKExWallet
	//OKExMargin *OKExMargin
}

func NewOKEx(config *APIConfig) *OKEx {
	okex := &OKEx{config: config}
	okex.OKExSpot = &OKExSpot{okex}
	okex.OKExFuture = &OKExFuture{OKEx: okex, Locker: new(sync.Mutex)}
	//okex.OKExWallet = &OKExWallet{okex}
	//okex.OKExMargin = &OKExMargin{okex}
	return okex
}

func (ok *OKEx) GetExchangeName() string {
	return OKEX
}

func (ok *OKEx) DoRequest(httpMethod, uri, reqBody string, response interface{}) ([]byte, error) {
	url := ok.config.Endpoint + uri
	sign, timestamp := ok.doParamSign(httpMethod, uri, reqBody)
	resp, err := NewHttpRequest(ok.config.HttpClient, httpMethod, url, reqBody, map[string]string{
		CONTENT_TYPE: APPLICATION_JSON_UTF8,
		ACCEPT:       APPLICATION_JSON,
		//COOKIE:               LOCALE + "en_US",
		OK_ACCESS_KEY:        ok.config.ApiKey,
		OK_ACCESS_PASSPHRASE: ok.config.ApiPassphrase,
		OK_ACCESS_SIGN:       sign,
		OK_ACCESS_TIMESTAMP:  fmt.Sprint(timestamp)})
	if err != nil {
		return nil, err
	} else {
		return resp, json.Unmarshal(resp, &response)
	}
}

func (ok *OKEx) adaptOrderState(state int) TradeStatus {
	switch state {
	case -2:
		return ORDER_FAIL
	case -1:
		return ORDER_CANCEL
	case 0:
		return ORDER_UNFINISH
	case 1:
		return ORDER_PART_FINISH
	case 2:
		return ORDER_FINISH
	case 3:
		return ORDER_UNFINISH
	case 4:
		return ORDER_CANCEL_ING
	}
	return ORDER_UNFINISH
}

/*
 Get a http request body is a json string and a byte array.
*/
func (ok *OKEx) BuildRequestBody(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		//log.Println(err)
		return "", nil, errors.New("json convert string error")
	}

	jsonBody := string(data)
	binBody := bytes.NewReader(data)

	return jsonBody, binBody, nil
}

func (ok *OKEx) doParamSign(httpMethod, uri, requestBody string) (string, string) {
	timestamp := ok.IsoTime()
	preText := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), uri, requestBody)
	//log.Println("preHash", preText)
	sign, _ := GetParamHmacSHA256Base64Sign(ok.config.ApiSecretKey, preText)
	return sign, timestamp
}

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func (ok *OKEx) IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func (ok *OKEx) LimitBuy(amount, price string, currency CurrencyPair) (*Order, []byte, error) {
	return ok.OKExSpot.LimitBuy(amount, price, currency)
}

func (ok *OKEx) LimitSell(amount, price string, currency CurrencyPair) (*Order, []byte, error) {
	return ok.OKExSpot.LimitSell(amount, price, currency)
}

func (ok *OKEx) MarketBuy(amount, price string, currency CurrencyPair) (*Order, []byte, error) {
	return ok.OKExSpot.MarketBuy(amount, price, currency)
}

func (ok *OKEx) MarketSell(amount, price string, currency CurrencyPair) (*Order, []byte, error) {
	return ok.OKExSpot.MarketSell(amount, price, currency)
}

func (ok *OKEx) CancelOrder(orderId string, currency CurrencyPair) (bool, []byte, error) {
	return ok.OKExSpot.OKExSpot.CancelOrder(orderId, currency)
}

func (ok *OKEx) GetOneOrder(orderId string, currency CurrencyPair) (*Order, []byte, error) {
	return ok.OKExSpot.GetOneOrder(orderId, currency)
}

func (ok *OKEx) GetUnfinishOrders(currency CurrencyPair) ([]Order, []byte, error) {
	return ok.OKExSpot.GetUnfinishOrders(currency)
}

func (ok *OKEx) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	return ok.OKExSpot.GetOrderHistorys(currency, currentPage, pageSize)
}

func (ok *OKEx) GetAccount() (*Account, []byte, error) {
	return ok.OKExSpot.GetAccount()
}

func (ok *OKEx) GetTicker(currency CurrencyPair) (*Ticker, []byte, error) {
	return ok.OKExSpot.GetTicker(currency)
}

func (ok *OKEx) GetDepth(size int, currency CurrencyPair) (*Depth, []byte, error) {
	return ok.OKExSpot.GetDepth(size, currency)
}

func (ok *OKEx) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, []byte, error) {
	return ok.OKExSpot.GetKlineRecords(currency, period, size, since)
}

func (ok *OKEx) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	return ok.OKExSpot.GetTrades(currencyPair, since)
}