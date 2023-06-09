package helpers

import (
	"api/bitcoin-api/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	httpsBinance       = "https://api.binance.com"
	httpsBinanceRoute  = "/api/v3/avgPrice"
	convertionCurrency = "BTCUAH"
	invalidPrice       = -1
)

func RequestPriceBinance() (float64, error) {
	params := url.Values{}
	params.Add("symbol", convertionCurrency)

	u, err := url.ParseRequestURI(httpsBinance)
	u.Path = httpsBinanceRoute
	u.RawQuery = params.Encode()
	finalUrl := fmt.Sprintf("%v", u)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while parsing the URL: %v\n", err)
		return invalidPrice, err
	}

	exchangeRate, err := http.Get(finalUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while requesting GET from the %s: %v\n",
			httpsBinance+httpsBinanceRoute, err)
		return invalidPrice, err
	}

	defer exchangeRate.Body.Close()

	body, err := ioutil.ReadAll(exchangeRate.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while reading the request body: %v\n", err)
		return invalidPrice, err
	}

	var exchangeRateObj models.ExchangeRate
	if err := json.Unmarshal(body, &exchangeRateObj); err != nil {
		fmt.Println(err.Error())
		return invalidPrice, err
	}

	result, _ := strconv.ParseFloat(exchangeRateObj.Price, 64)
	return result, nil
}
