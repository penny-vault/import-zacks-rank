package zacks

import (
	"os"
	"strings"

	"github.com/penny-vault/import-zacks-rank/common"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Download() ([]byte, string) {
	page, context, browser, pw := common.StartPlaywright(viper.GetBool("playwright.headless"))

	// block a variety of domains that contain trackers and ads
	page.Route("**/*", func(route playwright.Route, request playwright.Request) {
		if strings.Contains(request.URL(), "google.com") ||
			strings.Contains(request.URL(), "facebook.com") ||
			strings.Contains(request.URL(), "adsystem.com") ||
			strings.Contains(request.URL(), "sitescout.com") ||
			strings.Contains(request.URL(), "ipredictive.com") ||
			strings.Contains(request.URL(), "eyeota.net") {
			err := route.Abort("failed")
			if err != nil {
				log.Error().Err(err).Msg("failed blocking route")
			}
			return
		}

		if request.ResourceType() == "image" {
			err := route.Abort("failed")
			if err != nil {
				log.Error().Err(err).Msg("failed blocking image")
			}
		}

		route.Continue()
	})

	// load the login page
	if _, err := page.Goto(LOGIN_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Error().Err(err).Msg("could not load login page")
		return []byte{}, ""
	}

	page.WaitForSelector("#login input[name=username]")

	page.Type("#login input[name=username]", viper.GetString("zacks.username"))
	page.Type("#login input[name=password]", viper.GetString("zacks.password"))
	page.Click("#login input[value=Login]")

	// For some reason page.WaitForNavigation just times out here
	// substituting 1 second wait for the login to complete
	// page.WaitForNavigation()
	page.WaitForTimeout(1000)

	log.Info().Msg("wait for navigation completed")

	if _, err := page.Goto(STOCK_SCREENER_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Error().Err(err).Msg("could not load stock screener page")
		return []byte{}, ""
	}

	iframe, err := page.WaitForSelector("#screenerContent")
	if err != nil {
		log.Error().Err(err).Msg("could not load login page")
		return []byte{}, ""
	}

	frame, err := iframe.ContentFrame()
	if err != nil {
		log.Error().Err(err).Msg("could not get screener content frame")
		return []byte{}, ""
	}

	// navigate to saved screens tab
	frame.WaitForSelector("#my-screen-tab")
	frame.Click("#my-screen-tab")

	// navigate to our saved screen
	frame.WaitForSelector("#btn_run_137005")
	frame.Click("#btn_run_137005")

	// wait for the screen to load
	frame.WaitForSelector("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5")

	zacksPdfFn := viper.GetString("zacks.pdf")
	if zacksPdfFn != "" {
		log.Info().Str("fn", zacksPdfFn).Msg("saving PDF")
		if _, err = page.PDF(playwright.PagePdfOptions{
			Path: playwright.String(zacksPdfFn),
		}); err != nil {
			log.Error().Err(err).Msg("could not save page to PDF")
		}
	}

	var data []byte
	var outputFilename string

	if download, err := page.ExpectDownload(func() error {
		return frame.Click("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5")
	}); err != nil {
		log.Error().Err(err).Msg("download failed")
	} else {
		if path, err := download.Path(); err != nil {
			log.Error().Err(err).Msg("download failed")
		} else {
			outputFilename = download.SuggestedFilename()
			data, err = os.ReadFile(path)
			if err != nil {
				log.Error().Err(err).Msg("reading data failed")
			}
		}
	}

	common.StopPlaywright(page, context, browser, pw)

	return data, outputFilename
}
