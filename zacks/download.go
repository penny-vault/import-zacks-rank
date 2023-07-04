package zacks

import (
	"os"
	"strings"

	"github.com/penny-vault/import-zacks-rank/common"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func ensureLoggedIn(page playwright.Page) {
	if _, err := page.Goto(HOMEPAGE_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(10000),
	}); err != nil {
		log.Error().Err(err).Msg("waiting for network idle on home page timed out")
	}

	elems, err := page.QuerySelectorAll("#user_menu > li.welcome_usn")
	if err != nil {
		log.Error().Err(err).Msg("could not determine if user is logged in; error when calling QuerySelectorAll")
	}
	if len(elems) > 0 {
		// already logged in
		log.Info().Msg("user is already logged in")
		return
	}

	log.Info().Msg("need to log user in")

	// load the login page
	if _, err = page.Goto(LOGIN_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Error().Err(err).Msg("could not load login page")
		return
	}

	if _, err = page.WaitForSelector("#login input[name=username]"); err != nil {
		log.Error().Err(err).Msg("could not find username input field")
		return
	}

	if err = page.Type("#login input[name=username]", viper.GetString("zacks.username")); err != nil {
		log.Error().Err(err).Msg("could not type username")
		return
	}

	if err = page.Type("#login input[name=password]", viper.GetString("zacks.password")); err != nil {
		log.Error().Err(err).Msg("could not type password")
		return
	}

	if err = page.Click("#login input[value=Login]"); err != nil {
		log.Error().Err(err).Msg("could not click login button")
		return
	}

	// For some reason page.WaitForNavigation just times out here
	// substituting 1 second wait for the login to complete
	// page.WaitForNavigation()
	page.WaitForTimeout(1000)
}

// Download authenticates with the zacks webpage and downloads the results of the stock screen
// it returns the downloaded bytes, filename, and any errors that occur
func Download() (fileData []byte, outputFilename string, err error) {
	page, context, browser, pw := common.StartPlaywright(viper.GetBool("playwright.headless"))

	// block a variety of domains that contain trackers and ads
	page.Route("**/*", func(route playwright.Route) {
		request := route.Request()
		if strings.Contains(request.URL(), "google.com") ||
			strings.Contains(request.URL(), "googletagservices.com") ||
			strings.Contains(request.URL(), "googlesyndication.com") ||
			strings.Contains(request.URL(), "facebook.com") ||
			strings.Contains(request.URL(), "moatpixel.com") ||
			strings.Contains(request.URL(), "moatads.com") ||
			strings.Contains(request.URL(), "adsystem.com") ||
			strings.Contains(request.URL(), "amazon-adsystem.com") ||
			strings.Contains(request.URL(), "adnxs.com") ||
			strings.Contains(request.URL(), "lijit.com") ||
			strings.Contains(request.URL(), "3lift.com") ||
			strings.Contains(request.URL(), "doubleclick.net") ||
			strings.Contains(request.URL(), "bidswitch.net") ||
			strings.Contains(request.URL(), "casalemedia.com") ||
			strings.Contains(request.URL(), "yahoo.com") ||
			strings.Contains(request.URL(), "sitescout.com") ||
			strings.Contains(request.URL(), "ipredictive.com") ||
			strings.Contains(request.URL(), "uat5-b.investingchannel.com") ||
			strings.Contains(request.URL(), "eyeota.net") {
			log.Info().Str("URL", request.URL()).Msg("blocking request")
			err := route.Abort("failed")
			if err != nil {
				log.Error().Err(err).Msg("failed blocking route")
			}
			return
		}

		/*
			if request.ResourceType() == "image" {
				err := route.Abort("failed")
				if err != nil {
					log.Error().Err(err).Msg("failed blocking image")
				}
			}
		*/

		route.Continue()
	})

	ensureLoggedIn(page)

	log.Info().Msg("Load stock screener page")

	if _, err = page.Goto(STOCK_SCREENER_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Error().Err(err).Msg("could not load stock screener page")
		return
	}

	iframe, err := page.WaitForSelector("#screenerContent")
	if err != nil {
		log.Error().Err(err).Msg("could not load screener page")
		return
	}

	frame, err := iframe.ContentFrame()
	if err != nil {
		log.Error().Err(err).Msg("could not get screener content frame")
		return
	}

	log.Info().Msg("navigate to saved screens tab")

	// navigate to saved screens tab
	if _, err = frame.WaitForSelector("#my-screen-tab"); err != nil {
		log.Error().Err(err).Msg("wait for screener tabs failed")
		return
	}
	if err = frame.Click("#my-screen-tab"); err != nil {
		log.Error().Err(err).Msg("click tab button failed")
		return
	}

	log.Info().Msg("run the saved stock screen")

	// navigate to our saved screen
	if _, err = frame.WaitForSelector("#btn_run_137005"); err != nil {
		log.Error().Err(err).Msg("wait for run button failed")
		return
	}

	log.Info().Msg("clicking run button")

	if err = frame.Click("#btn_run_137005"); err != nil {
		log.Error().Err(err).Msg("click run button failed")
		return
	}

	log.Info().Msg("button clicked")

	// wait for the screen to finish running
	if _, err = frame.WaitForSelector("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5"); err != nil {
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
		return frame.Click("#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5")
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
