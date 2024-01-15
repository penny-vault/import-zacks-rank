/*
Copyright 2022

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package zacks

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/penny-vault/import-zacks-rank/common"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
)

func BalanceSheet(tickers []string) (BalanceSheetList, error) {
	page, context, browser, pw := common.StartPlaywright(viper.GetBool("playwright.headless"))

	result := make([]*BalanceSheetRecord, 0, len(tickers)*5)

	bar := progressbar.NewOptions(len(tickers),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetDescription("Preparing to download ..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	completed := 0
	for _, ticker := range tickers {
		bar.Describe(ticker)

		zacksTicker := strings.ReplaceAll(ticker, "/", ".")

		if _, err := page.Goto(fmt.Sprintf("https://www.zacks.com/stock/quote/%s/balance-sheet?icid=quote-stock_overview-quote_nav_tracking-zcom-left_subnav_quote_navbar-balance_sheet", zacksTicker), playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
			Timeout:   playwright.Float(20000),
		}); err != nil {
			log.Error().Err(err).Msg("waiting for network idle on home page timed out")
		}

		page.SetDefaultTimeout(1000)

		// slow things down a bit so we don't over-whelm zacks.com
		time.Sleep(5 * time.Second)

		// Annual Income Statement

		// get section header
		annual := make(map[string]*BalanceSheetRecord, 5)
		colMap := make(map[int]string, 5)
		if err := parseHeader("#annual_income_statement", ticker, "As-Reported-Annual", page, annual, colMap); err != nil {
			// add to database
			AddExclusion(ticker)
			continue
		}

		if len(colMap) > 0 {
			// current assets
			parseRow("#annual_income_statement", "Total Current Assets", "TotalCurrentAssets", page, annual, colMap)
			// current liabilities
			parseRow("#annual_income_statement", "Total Current Liabilities", "TotalCurrentLiabilities", page, annual, colMap)
		}

		allNaN := true

		// add all ARY dimension to return val
		for _, v := range annual {
			result = append(result, v)
			allNaN = (math.IsNaN(v.TotalCurrentAssets) && math.IsNaN(v.TotalCurrentLiabilities)) && allNaN
		}

		// Quarterly Income Statement

		if err := page.GetByRole("tablist").GetByRole("link", playwright.LocatorGetByRoleOptions{
			Name: "Quarterly Balance Sheet",
		}).Click(); err != nil {
			log.Error().Err(err).Msg("could not switch to quarterly balance sheet")
		}

		// get section header
		quarterly := make(map[string]*BalanceSheetRecord, 5)
		colMap = make(map[int]string, 5)
		parseHeader("#quarterly_income_statement", ticker, "As-Reported-Quarterly", page, quarterly, colMap)

		if len(colMap) > 0 {
			// current assets
			parseRow("#quarterly_income_statement", "Total Current Assets", "TotalCurrentAssets", page, quarterly, colMap)
			// current liabilities
			parseRow("#quarterly_income_statement", "Total Current Liabilities", "TotalCurrentLiabilities", page, quarterly, colMap)
		}

		// add all ARQ dimension to return val
		for _, v := range quarterly {
			result = append(result, v)
			allNaN = (math.IsNaN(v.TotalCurrentAssets) && math.IsNaN(v.TotalCurrentLiabilities)) && allNaN
		}

		bar.Add(1)
		completed += 1

		if allNaN {
			AddExclusion(ticker)
		}

		// every 50 tickers restart playwright
		if completed > 50 {
			common.StopPlaywright(page, context, browser, pw)
			page, context, browser, pw = common.StartPlaywright(viper.GetBool("playwright.headless"))
			completed = 0
		}
	}

	common.StopPlaywright(page, context, browser, pw)
	return result, nil
}

func parseHeader(selector string, ticker string, dim string, page playwright.Page, table map[string]*BalanceSheetRecord, colMap map[int]string) error {
	if row, err := page.Locator(selector).GetByRole("rowgroup").First().TextContent(); err != nil {
		log.Error().Err(err).Msg("could not get row header")
		return err
	} else {
		row = strings.TrimSpace(row)
		cols := strings.Split(row, "\n")
		for idx, heading := range cols {
			colName := strings.Trim(heading, " \t")
			table[colName] = &BalanceSheetRecord{
				Ticker:       ticker,
				CalendarDate: colName,
				Dimension:    dim,
				DownloadDate: time.Now(),
			}
			colMap[idx] = colName
		}
	}

	return nil
}

func parseRow(selector string, rowLabel string, fieldName string, page playwright.Page, table map[string]*BalanceSheetRecord, colMap map[int]string) {
	if row, err := page.Locator(selector).GetByRole("row", playwright.LocatorGetByRoleOptions{Name: rowLabel}).TextContent(); err != nil {
		log.Error().Err(err).Str("dimension", selector).Str("rowLabel", rowLabel).Msg("could not get row")
	} else {
		row = strings.Trim(row, " \t")
		cols := strings.Split(row, "\n")[2:]
		for idx, field := range cols {
			val := strings.TrimSpace(field)
			val = strings.ReplaceAll(val, ",", "")

			if val == rowLabel {
				continue
			}

			if _, ok := colMap[idx]; !ok {
				continue
			}

			if _, ok := table[colMap[idx]]; !ok {
				continue
			}

			if val == "NA" || val == "" {
				reflect.ValueOf(table[colMap[idx]]).Elem().FieldByName(fieldName).Set(reflect.ValueOf(math.NaN()))
				continue
			}

			if floatVal, err := strconv.ParseFloat(val, 64); err != nil {
				log.Error().Err(err).Str("inputVal", val).Str("column", colMap[idx]).Msg("could not convert value to float")
			} else {
				floatVal *= 1e6
				if floatVal < 0 {
					floatVal = math.NaN()
				}
				reflect.ValueOf(table[colMap[idx]]).Elem().FieldByName(fieldName).Set(reflect.ValueOf(floatVal))
			}
		}
	}
}
