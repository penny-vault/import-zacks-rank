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
	"github.com/playwright-community/playwright-go"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func EnsureLoggedIn(page playwright.Page) {
	if _, err := page.Goto(HOMEPAGE_URL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(10000),
	}); err != nil {
		log.Error().Err(err).Msg("waiting for network idle on home page timed out")
	}

	locator := page.Locator("#user_menu > li.welcome_usn")
	if visible, err := locator.IsVisible(); visible {
		// already logged in
		log.Info().Msg("user is already logged in")
		return
	} else if err != nil {
		log.Error().Err(err).Msg("encountered error when checking if user logged in")
	}

	log.Info().Msg("need to log user in")

	// load the login page
	if _, err := page.Goto(LOGIN_URL); err != nil {
		log.Error().Err(err).Msg("could not load login page")
		return
	}

	if err := page.Locator("#login input[name=username]").Fill(viper.GetString("zacks.username")); err != nil {
		log.Error().Err(err).Msg("could not fill username")
		return
	}

	if err := page.Locator("#login input[name=password]").Fill(viper.GetString("zacks.password")); err != nil {
		log.Error().Err(err).Msg("could not fill password")
		return
	}

	if err := page.Locator("#login input[value=Login]").Click(); err != nil {
		log.Error().Err(err).Msg("could not click login button")
		return
	}
}
