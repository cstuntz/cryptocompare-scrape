package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
)

//APIRes Top level response from the API
type APIRes struct {
	Raw     RawCoinTypes     `json:"RAW"`
	Display DisplayCoinTypes `json:"DISPLAY"`
}

//CoinTypes Raw Data for each of the crypto coins
type RawCoinTypes struct {
	ETH RawCurrencyCodes `json:"ETH"`
	BTC RawCurrencyCodes `json:"BTC"`
}

//CurrencyCodes Currency conversions for crypto currencies
type RawCurrencyCodes struct {
	USD RawData `json:"USD"` // US Dollar
	EUR RawData `json:"EUR"` // Euro
	GBP RawData `json:"GBP"` // Great Britain Pound
	CNY RawData `json:"CNY"` // Chinese Yuan
	JPY RawData `json:"JPY"` // Japanese Yen
}

//RawData Raw data for each conversion
type RawData struct {
	Type            string  `json:"TYPE"`
	Market          string  `json:"MARKET"`
	FromSymbol      string  `json:"FROMSYMBOL"`
	ToSymbol        string  `json:"TOSYMBOL"`
	Flags           string  `json:"FLAGS"`
	Price           float32 `json:"PRICE"`
	LastUpdate      int     `json:"LASTUPDATE"`
	LastVolume      float32 `json:"LASTVOLUME"`
	LastVolumeTo    float32 `json:"LASTVOLUMETO"`
	LastTradeID     float32 `json:"LASTTRADEID"`
	Volume24Hour    float32 `json:"VOLUME24HOUR"`
	Volume24HourTo  float32 `json:"VOLUME24HOURTO"`
	Open24Hour      float32 `json:"OPEN24HOUR"`
	High24Hour      float32 `json:"HIGH24HOUR"`
	Low24Hour       float32 `json:"LOW24HOUR"`
	LastMarket      string  `json:"LASTMARKET"`
	Change24Hour    float32 `json:"CHANGE24HOUR"`
	ChangePct24Hour float32 `json:"CHANGEPCT24HOUR"`
	Supply          float32 `json:"SUPPLY"`
	MarketCap       float32 `json:"MKTCAP"`
}

//CoinTypes Raw Data for each of the crypto coins
type DisplayCoinTypes struct {
	ETH DisplayCurrencyCodes `json:"ETH"`
	BTC DisplayCurrencyCodes `json:"BTC"`
}

//CurrencyCodes Currency conversions for crypto currencies
type DisplayCurrencyCodes struct {
	USD DisplayData `json:"USD"` // US Dollar
	EUR DisplayData `json:"EUR"` // Euro
	GBP DisplayData `json:"GBP"` // Great Britain Pound
	CNY DisplayData `json:"CNY"` // Chinese Yuan
	JPY DisplayData `json:"JPY"` // Japanese Yen
}

//DisplayData Data formatted for display purposes
type DisplayData struct {
	FromSymbol      string  `json:"FROMSYMBOL"`
	ToSymbol        string  `json:"TOSYMBOL"`
	Market          string  `json:"MARKET"`
	Price           string  `json:"PRICE"`
	LastUpdate      string  `json:"LASTUPDATE"`
	LastVolume      string  `json:"LASTVOLUME"`
	LastVolumeTo    string  `json:"LASTVOLUMETO"`
	LastTradeID     float32 `json:"LASTTRADEID"`
	Volume24Hour    string  `json:"VOLUME24HOUR"`
	Volume24HourTo  string  `json:"VOLUME24HOURTO"`
	Open24Hour      string  `json:"OPEN24HOUR"`
	High24Hour      string  `json:"HIGH24HOUR"`
	Low24Hour       string  `json:"LOW24HOUR"`
	LastMarket      string  `json:"LASTMARKET"`
	Change24Hour    string  `json:"CHANGE24HOUR"`
	ChangePct24Hour string  `json:"CHANGEPCT24HOUR"`
	Supply          string  `json:"SUPPLY"`
	MarketCap       string  `json:"MKTCAP"`
}

func main() {

	cryptoCompareURL := os.Getenv("CRYPTOURL")
	mongoURL := os.Getenv("MONGODBURL")
	mongoDBName := os.Getenv("MONGODBNAME")

	//Connect to MongoDB
	mdb, err := mgo.Dial(mongoURL)
	failOnErr("Failed to connect to mongodb", err)
	defer mdb.Close()

	//Call the cryptocompare API
	body := callURL(cryptoCompareURL)

	//Parse the body of the API call to a struct
	var p APIRes
	parseErr := json.Unmarshal(body, &p)
	failOnErr("Failed to parse json", parseErr)

	//Add info to the MongoDB
	c := mdb.DB(mongoDBName).C("eth_raw")
	insertErr := c.Insert(p.Raw.ETH.USD)
	failOnErr("Failed to insert to MongoBD", insertErr)

	//prettyPrint(p)
	//oneLinePrint(p)

}

func callURL(url string) []byte {

	res, err := http.Get(url)
	failOnErr("Failed to get the URL", err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	failOnErr("Failed to parse URL body", err)

	return body

}

func prettyPrint(p APIRes) {
	out, err := json.MarshalIndent(p, "", "  ")
	failOnErr("Failed to marshal json", err)
	fmt.Println(string(out))
}

func oneLinePrint(p APIRes) {
	out, err := json.Marshal(p)
	failOnErr("Failed to marshal json", err)
	fmt.Println(string(out))
}

func pushAll(p APIRes) {

}

func failOnErr(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s\n", msg, err)
	}
}
