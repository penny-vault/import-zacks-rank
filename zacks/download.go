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
	"os"

	"github.com/penny-vault/import-zacks-rank/common"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Download authenticates with the zacks webpage and downloads the results of the stock screen
// it returns the downloaded bytes, filename, and any errors that occur
func Download() (fileData []byte, outputFilename string, err error) {
	page, context, browser, pw := common.StartPlaywright(viper.GetBool("playwright.headless"))

	EnsureLoggedIn(page)

	log.Info().Msg("Load stock screener page")

	if _, err = page.Goto(STOCK_SCREENER_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Error().Err(err).Msg("could not load stock screener page")
		return
	}

	frame := page.FrameLocator("#screenerContent")

	log.Info().Msg("navigate to saved screens tab")

	if err = frame.Locator("#my-screen-tab").Click(); err != nil {
		log.Error().Err(err).Msg("click tab button failed")
		return
	}

	log.Info().Msg("run the saved stock screen")

	// navigate to our saved screen

	log.Info().Msg("clicking run button")

	if err = frame.Locator("#btn_run_137005").Click(); err != nil {
		log.Error().Err(err).Msg("click run button failed")
		return
	}

	log.Info().Msg("button clicked")

	// wait for the screen to finish running
	if err = frame.Locator("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5").WaitFor(); err != nil {
		log.Error().Err(err).Msg("wait for 'csv' download selector failed")
		return
	}

	zacksPdfFn := viper.GetString("zacks.pdf")
	if zacksPdfFn != "" {
		log.Info().Str("fn", zacksPdfFn).Msg("saving PDF")
		if _, err = page.PDF(playwright.PagePdfOptions{
			Path: playwright.String(zacksPdfFn),
		}); err != nil {
			log.Error().Err(err).Msg("could not save page to PDF")
		}
	}

	var download playwright.Download
	if download, err = page.ExpectDownload(func() error {
		return frame.Locator("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5").Click()
	}); err != nil {
		log.Error().Err(err).Msg("download failed")
	}

	var path string
	if path, err = download.Path(); err != nil {
		log.Error().Err(err).Msg("download failed")
	} else {
		outputFilename = download.SuggestedFilename()
		fileData, err = os.ReadFile(path)
		if err != nil {
			log.Error().Err(err).Msg("reading data failed")
			return
		}
	}

	common.StopPlaywright(page, context, browser, pw)
	return
}
