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

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/penny-vault/import-zacks-rank/backblaze"
	"github.com/penny-vault/import-zacks-rank/zacks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "import-zacks-rank",
	Short: "Download and import ratings from Zacks stock screener",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var outputFilename string
		var err error

		for ii := 0; ii < viper.GetInt("zacks.max_retries"); ii++ {
			data, outputFilename, err = zacks.Download()
			if err == nil {
				break
			}
		}
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLog)

	// Persistent flags that are global to application
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.import-zacks-rank.toml)")
	rootCmd.PersistentFlags().Bool("log-json", false, "print logs as json to stderr")
	viper.BindPFlag("log.json", rootCmd.PersistentFlags().Lookup("log-json"))

	// Add flags
	rootCmd.Flags().StringP("database_url", "d", "host=localhost port=5432", "DSN for database connection")
	viper.BindPFlag("database.url", rootCmd.Flags().Lookup("database_url"))

	rootCmd.Flags().Uint32P("limit", "l", 0, "limit results to N")
	viper.BindPFlag("limit", rootCmd.Flags().Lookup("limit"))

	rootCmd.Flags().StringP("backblaze_bucket", "b", "zacks-investment", "Backblaze bucket name")
	viper.BindPFlag("backblaze.bucket", rootCmd.Flags().Lookup("backblaze_bucket"))

	rootCmd.Flags().String("backblaze_application_id", "<not-set>", "Backblaze application id")
	viper.BindPFlag("backblaze.application_id", rootCmd.Flags().Lookup("backblaze_application_id"))

	rootCmd.Flags().String("backblaze_application_key", "<not-set>", "Backblaze application key")
	viper.BindPFlag("backblaze.application_key", rootCmd.Flags().Lookup("backblaze_application_key"))

	rootCmd.Flags().String("zacks-pdf", "", "Save page to PDF for debug purposes")
	viper.BindPFlag("zacks.pdf", rootCmd.Flags().Lookup("zacks-pdf"))

	rootCmd.Flags().Int("max-retries", 3, "maximum number of times to retry if download fails")
	viper.BindPFlag("zacks.max_retries", rootCmd.Flags().Lookup("max-retries"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".import-zacks-rank" (without extension).
		viper.AddConfigPath("/etc/") // path to look for the config file in
		viper.AddConfigPath(fmt.Sprintf("%s/.config", home))
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("import-zacks-rank")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Str("ConfigFile", viper.ConfigFileUsed()).Msg("Loaded config file")
	} else {
		log.Error().Err(err).Msg("error reading config file")
	}
}

func initLog() {
	if !viper.GetBool("log.json") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
