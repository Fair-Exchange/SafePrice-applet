package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"./icon"

	"github.com/getlantern/systray"
)

var validPrecisionValues = []time.Duration{
	1 * time.Minute,
	5 * time.Minute,
	15 * time.Minute,
	30 * time.Minute,
	1 * time.Hour,
	2 * time.Hour,
	4 * time.Hour,
	6 * time.Hour,
	12 * time.Hour,
	24 * time.Hour,
	72 * time.Hour,
	168 * time.Hour,
}

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	systray.SetIcon(icon.Icon)
	systray.SetTitle("SafePrice")
	systray.SetTooltip("Safecoin price")

	var nosafetrade bool
	var nocrex24 bool
	var updatetime time.Duration
	var precision time.Duration
	flag.BoolVar(&nosafetrade, "nosafetrade", false, "Disable SafeTrade price check")
	flag.BoolVar(&nocrex24, "nocrex24", false, "Disable CREX24 price check")
	flag.DurationVar(&updatetime, "interval", 5*time.Minute, "Set different update interval")
	flag.DurationVar(&precision, "precision", 15*time.Minute, "Percentage change precision (higher precision needs more bandwidth and more CPU usage when updated)\nSupported values: 1m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 12h, 1g, 3g, 7g")

	flag.Parse()

	valid := false
	for _, value := range validPrecisionValues {
		if precision == value {
			valid = true
			break
		}
	}
	if !valid {
		fmt.Println("[ERR] Invalid precision value")
		os.Exit(1)
	}

	if !nosafetrade {
		mSafetrade := systray.AddMenuItem("SafeTrade:", "Update SafeTrade price")

		go func() {
			type tickers struct {
				Ticker struct {
					Last string
				}
			}
			type k [][]float64
			for {
				resp, err := http.Get("https://safe.trade/api/v2/tickers/safebtc")
				if err != nil {
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				var api tickers
				if err = json.Unmarshal(body, &api); err != nil {
					continue
				}

				resp, err = http.Get(fmt.Sprintf("https://safe.trade/api/v2/k?market=safebtc&period=%.f&limit=%.f", precision.Minutes(), 1440/precision.Minutes()+1))
				if err != nil {
					continue
				}
				body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				var kdata k
				if err = json.Unmarshal(body, &kdata); err != nil {
					continue
				}

				last, _ := strconv.ParseFloat(api.Ticker.Last, 64)
				mSafetrade.SetTitle(fmt.Sprintf("Safetrade:\t%.8f\tBTC\t%+.2f%%", last, 100-100*kdata[0][3]/last))
				select {
				case <-mSafetrade.ClickedCh:
				case <-time.After(updatetime):
					continue
				}
			}
		}()
	}

	if !nocrex24 {
		mCREX24 := systray.AddMenuItem("CREX24:", "Update CREX24 price")

		go func() {
			type tickers []struct {
				Last          float32
				PercentChange float32
			}
			for {
				resp, err := http.Get("https://api.crex24.com/v2/public/tickers?instrument=SAFE-BTC")
				if err != nil {
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				var api tickers
				if err = json.Unmarshal(body, &api); err != nil {
					continue
				}

				mCREX24.SetTitle(fmt.Sprintf("CREX24:\t%.8f\tBTC\t%+.2f%%", api[0].Last, api[0].PercentChange))
				select {
				case <-mCREX24.ClickedCh:
				case <-time.After(updatetime):
					continue
				}
			}
		}()
	}

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {}
