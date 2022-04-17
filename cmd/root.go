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

	"github.com/penny-vault/import-zacks-rank/zacksimport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "import-zacks-rank [dirpath]",
	Short: "Import CSV ratings downloaded from Zacks stock screener",
	// Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		ratings := zacksimport.LoadRatings(args[0], viper.GetInt("limit"))

		fmt.Printf("%+v\n", ratings[0])

		/*
			zacksimport.EnrichWithFigi(ratings)
			zacksimport.SaveToDB(ratings)

			// Save data as parquet to a temporary directory
			tmpdir, err := os.MkdirTemp(os.TempDir(), "import-zacks")
			if err != nil {
				log.Error().Str("OriginalError", err.Error()).Msg("could not create tempdir")
			}

			parquetFn := fmt.Sprintf("%s/zacks-%s.parquet", tmpdir, ratings[0].Date.Format("20060102"))
			log.Info().Str("FileName", parquetFn).Msg("writing zacks ratings data to parquet")
			zacksimport.SaveToParquet(ratings, parquetFn)

			// Upload to backblaze
			backblaze.UploadToBackBlaze(parquetFn, viper.GetString("backblaze_bucket"), ratings[0].Date.Format("2006"))
		*/

		// Cleanup after ourselves
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

	// Persistent flags that are global to application
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.import-zacks-rank.toml)")

	// Add flags
	rootCmd.Flags().StringP("database_url", "d", "host=localhost port=5432", "DSN for database connection")
	viper.BindPFlag("database_url", rootCmd.Flags().Lookup("database_url"))

	rootCmd.Flags().Uint32P("limit", "l", 0, "limit results to N")
	viper.BindPFlag("limit", rootCmd.Flags().Lookup("limit"))

	rootCmd.Flags().StringP("backblaze_bucket", "b", "zacks-investment", "Backblaze bucket name")
	viper.BindPFlag("backblaze_bucket", rootCmd.Flags().Lookup("backblaze_bucket"))

	rootCmd.Flags().String("backblaze_application_id", "<not-set>", "Backblaze application id")
	viper.BindPFlag("backblaze_application_id", rootCmd.Flags().Lookup("backblaze_application_id"))

	rootCmd.Flags().String("backblaze_application_key", "<not-set>", "Backblaze application key")
	viper.BindPFlag("backblaze_application_key", rootCmd.Flags().Lookup("backblaze_application_key"))
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
		viper.AddConfigPath("/etc/import-zacks-rank/") // path to look for the config file in
		viper.AddConfigPath(fmt.Sprintf("%s/.import-zacks-rank", home))
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("import-zacks-rank")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
