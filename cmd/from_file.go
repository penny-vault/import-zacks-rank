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
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/penny-vault/import-zacks-rank/backblaze"
	"github.com/penny-vault/import-zacks-rank/zacks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inputFile string

var fileCmd = &cobra.Command{
	Use:   "file",
	Args:  cobra.ExactArgs(1),
	Short: "load zacks rank from file",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		outputFilename := args[0]

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

		// read data from file
		fh, err := os.Open(outputFilename)
		if err != nil {
			log.Fatal().Str("Filename", outputFilename).Err(err).Msg("could not read input file")
		}

		data, err := io.ReadAll(fh)
		if err != nil {
			log.Fatal().Err(err).Msg("could not read input file")
		}

		// parse file
		ratings := zacks.LoadRatings(data, dateStr, viper.GetInt("limit"))
		log.Info().Int("NumRatings", len(ratings)).Msg("loaded ratings")
		if len(ratings) == 0 {
			log.Fatal().Msg("no ratings returned")
		}
		zacks.EnrichWithFigi(ratings)

		// Save data as parquet to a temporary directory
		tmpdir, err := os.MkdirTemp(os.TempDir(), "import-zacks")
		if err != nil {
			log.Error().Err(err).Msg("could not create tempdir")
		}
		dateStr = strings.ReplaceAll(dateStr, "-", "")
		parquetFn := fmt.Sprintf("%s/zacks-%s.parquet", tmpdir, dateStr)
		log.Info().Str("FileName", parquetFn).Msg("writing zacks ratings data to parquet")
		zacks.SaveToParquet(ratings, parquetFn)

		// Save to database
		zacks.SaveToDB(ratings)

		// Upload to backblaze
		year := string(dateStr[:4])
		log.Info().Str("Year", year).Str("Bucket", viper.GetString("backblaze.bucket")).Msg("data")
		backblaze.UploadToBackBlaze(parquetFn, viper.GetString("backblaze.bucket"), year)

		// Cleanup after ourselves
		os.RemoveAll(tmpdir)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
