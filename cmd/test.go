// Copyright 2022
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"regexp"

	"github.com/penny-vault/import-zacks-rank/zacks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test downloading zacks ratings",
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var outputFilename string
		var err error

		data, outputFilename, err = zacks.Download()
		// after multiple retries check if the download succeeded
		if err != nil {
			os.Exit(1)
		}

		// parse date from filename :(
		// if that doesn't work use the current date
		regex := regexp.MustCompile(`zacks_custom_screen_(\d{4}-\d{2}-\d{2})`)
		match := regex.FindAllStringSubmatch(outputFilename, -1)
		var dateStr string
		if len(match) > 0 {
			dateStr = match[0][1]
		} else {
			log.Error().Str("FileName", outputFilename).Msg("cannot extract date from filename, expecting zacks_custom_screen_YYYY-MM-DD")
			return
		}

		ratings := zacks.LoadRatings(data, dateStr, viper.GetInt("limit"))
		log.Info().Int("NumRatings", len(ratings)).Msg("loaded ratings")
		if len(ratings) == 0 {
			log.Fatal().Msg("no ratings returned")
		}
		zacks.EnrichWithFigi(ratings)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
