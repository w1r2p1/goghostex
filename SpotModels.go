package goghostex

import (
	"net/http"
	"time"
)

/*
	models about account
*/
type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[Currency]SubAccount
}

type SubAccount struct {
	Currency     Currency
	Amount       float64
	ForzenAmount float64
	LoanAmount   float64
}

type MarginAccount struct {
	Sub              map[Currency]MarginSubAccount
	LiquidationPrice float64
	RiskRate         float64
	MarginRatio      float64
}

type MarginSubAccount struct {
	Balance     float64
	Frozen      float64
	Available   float64
	CanWithdraw float64
	Loan        float64
	LendingFee  float64
}

/**
 * models about market
 **/

type Kline struct {
	Pair      CurrencyPair
	Timestamp int64
	Date      string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Vol       float64
}

type Ticker struct {
	Pair      CurrencyPair `json:"-"`
	Last      float64      `json:"last"`
	Buy       float64      `json:"buy"`
	Sell      float64      `json:"sell"`
	High      float64      `json:"high"`
	Low       float64      `json:"low"`
	Vol       float64      `json:"vol"`
	Timestamp uint64       `json:"timestamp"` // unit:ms
	Date      string       `json:"date"`      // date: format yyyy-mm-dd HH:MM:SS, the timezone define in apiconfig
}

// record
type Trade struct {
	Tid       int64        `json:"tid"`
	Type      TradeSide    `json:"type"`
	Amount    float64      `json:"amount,string"`
	Price     float64      `json:"price,string"`
	Timestamp uint64       `json:"timestamp"`
	Pair      CurrencyPair `json:"omitempty"`
}

type DepthRecord struct {
	Price  float64
	Amount float64
}

type DepthRecords []DepthRecord

func (dr DepthRecords) Len() int {
	return len(dr)
}

func (dr DepthRecords) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthRecords) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	//ContractType string //for future
	Pair      CurrencyPair
	Timestamp uint64
	Date      string
	AskList   DepthRecords // Descending order
	BidList   DepthRecords // Descending order
}

/*
	models about trade
*/

type Order struct {
	Price          float64
	Amount         float64
	AvgPrice       float64
	DealAmount     float64
	Fee            float64
	Cid            string //define by yourself
	OrderId        string
	OrderTimestamp uint64
	OrderDate      string
	Status         TradeStatus
	Currency       CurrencyPair
	Side           TradeSide
	//Type           string    // taker maker
	OrderType OrderType //0:NORMAL,1:ONLY_MAKER,2:FOK,3:IOC
}

/**
 *
 * models about API config
 *
 **/
type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string //for okex.com v3 api
	ClientId      string //for bitstamp.net , huobi.pro
	Location      *time.Location
}
