package common

import (
	"strings"

	"github.com/go-rod/stealth"
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// StealthPage creates a new playwright page with stealth js loaded to prevent bot detection
func StealthPage(context *playwright.BrowserContext) playwright.Page {
	page, err := (*context).NewPage()
	if err != nil {
		log.Error().Err(err).Msg("could not create page")
	}

	if err = page.AddInitScript(playwright.PageAddInitScriptOptions{
		Script: playwright.String(stealth.JS),
	}); err != nil {
		log.Error().Err(err).Msg("could not load stealth mode")
	}

	return page
}

// BuildUserAgent dynamically determines the user agent and removes the headless identifier
func BuildUserAgent(browser *playwright.Browser) string {
	context, err := (*browser).NewContext()
	if err != nil {
		log.Error().Err(err).Msg("could not create context for building user agent")
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Error().Err(err).Msg("could not create page BuildUserAgent")
	}

	resp, err := page.Goto("https://playwright.dev", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		log.Error().Err(err).Str("Url", "https://playwright.dev").Msg("could not load page")
	}

	headers, err := resp.Request().AllHeaders()
	if err != nil {
		log.Error().Err(err).Msg("could not load request headers")
	}

	userAgent := headers["user-agent"]
	userAgent = strings.Replace(userAgent, "Headless", "", -1)
	return userAgent
}

// StartPlaywright starts the playwright server and browser, it then creates a new context and page with the stealth extensions loaded
func StartPlaywright(headless bool) (page playwright.Page, context playwright.BrowserContext, browser playwright.Browser, pw *playwright.Playwright) {
	pw, err := playwright.Run()
	if err != nil {
		log.Error().Err(err).Msg("could not launch playwright")
	}

	browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		log.Error().Err(err).Msg("could not launch Chromium")
	}

	log.Info().Bool("Headless", headless).Str("ExecutablePath", pw.Chromium.ExecutablePath()).Str("BrowserVersion", browser.Version()).Msg("starting playwright")

	// calculate user-agent
	userAgent := viper.GetString("user_agent")
	if userAgent == "" {
		userAgent = BuildUserAgent(&browser)
	}
	log.Info().Str("UserAgent", userAgent).Msg("using user-agent")

	// create context
	context, err = browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(userAgent),
	})
	if err != nil {
		log.Error().Msg("could not create browser context")
	}

	// get a page
	page = StealthPage(&context)

	return
}

func StopPlaywright(page playwright.Page, context playwright.BrowserContext, browser playwright.Browser, pw *playwright.Playwright) {
	log.Info().Msg("closing context")
	if err := context.Close(); err != nil {
		log.Error().Err(err).Msg("error encountered when closing context")
	}

	log.Info().Msg("closing browser")
	if err := browser.Close(); err != nil {
		log.Error().Err(err).Msg("error encountered when closing browser")
	}

	log.Info().Msg("stopping playwright")
	if err := pw.Stop(); err != nil {
		log.Error().Err(err).Msg("error encountered when stopping playwright")
	}
}
