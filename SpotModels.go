package goghostex

import (
	"errors"
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
	Sequence  uint64 // The increasing sequence, cause the http return sequence is not sure.
	Date      string
	AskList   DepthRecords // Ascending order
	BidList   DepthRecords // Descending order
}

// check the depth data is right
func (depth *Depth) Check() error {
	AskCount := len(depth.AskList)
	BidCount := len(depth.BidList)

	if BidCount < 10 || AskCount < 10 {
		return errors.New("The ask_list or bid_list not enough! ")
	}

	for i := 1; i < AskCount; i++ {
		pre := depth.AskList[i-1]
		last := depth.AskList[i]
		if pre.Price >= last.Price {
			return errors.New("The ask_list is not ascending ordered! ")
		}
	}

	for i := 1; i < BidCount; i++ {
		pre := depth.BidList[i-1]
		last := depth.BidList[i]
		if pre.Price <= last.Price {
			return errors.New("The bid_list is not descending ordered! ")
		}
	}

	return nil
}

/*
	models about trade
*/

type Order struct {
	// cid is important, when the order api return wrong, you can find it in unfinished api
	Cid            string
	Price          float64
	Amount         float64
	AvgPrice       float64
	DealAmount     float64
	Fee            float64
	OrderId        string
	OrderTimestamp uint64
	OrderDate      string
	Status         TradeStatus
	Currency       CurrencyPair
	Side           TradeSide
	//0:NORMAL,1:ONLY_MAKER,2:FOK,3:IOC
	OrderType OrderType
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
