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
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/penny-vault/import-zacks-rank/zacks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	lookback  int
	maxAssets int
)

var balanceSheetCmd = &cobra.Command{
	Use:   "balance-sheet <tickers> ...",
	Args:  cobra.MinimumNArgs(0),
	Short: "load balance sheet from zacks",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		conn, err := pgx.Connect(ctx, viper.GetString("database.url"))
		if err != nil {
			log.Error().Err(err).Msg("Could not connect to database")
		}
		defer conn.Close(ctx)

		if len(args) == 0 {
			// get exclusions
			exclusion := make(map[string]bool)
			if rows, err := conn.Query(ctx, "SELECT distinct composite_figi FROM zacks_balance_sheet_exclusions"); err != nil {
				log.Fatal().Err(err).Msg("error querying database for tickers without working_capital")
			} else {
				var figi string

				for rows.Next() {
					if err := rows.Scan(&figi); err != nil {
						log.Fatal().Err(err).Msg("unable to scan query exclusion into figi string")
					}

					exclusion[figi] = true
				}
			}

			// get tickers
			if rows, err := conn.Query(ctx, "SELECT distinct ticker, composite_figi FROM fundamentals WHERE event_date > $1 AND working_capital = 'NaN'::float8 AND dim='As-Reported-Quarterly'", time.Now().Add(-90*24*time.Hour)); err != nil {
				log.Fatal().Err(err).Msg("error querying database for tickers without working_capital")
			} else {
				var (
					ticker string
					figi   string
				)

				cnt := 0
				for rows.Next() {
					cnt += 1
					if err := rows.Scan(&ticker, &figi); err != nil {
						log.Fatal().Err(err).Msg("unable to scan query value into ticker string")
					}

					if _, ok := exclusion[figi]; !ok {
						args = append(args, ticker)
					}
				}

				log.Info().Int("Count", cnt).Int("LenArgs", len(args)).Msg("found assets with missing working capital in database")
			}

			// shuffle the list
			for i := range args {
				j := rand.Intn(i + 1)
				args[i], args[j] = args[j], args[i]
			}

			log.Info().Int("Count", len(args)).Int("lookback", lookback).Msg("found records missing working_capital")

			// limit run to maxAssets items
			if len(args) > maxAssets {
				args = args[:maxAssets]
			}
		}

		if len(args) == 0 {
			log.Error().Msg("No assets to lookup")
			os.Exit(0)
		}

		if balanceSheets, err := zacks.BalanceSheet(args); err == nil {
			log.Info().Int("Count", len(balanceSheets)).Msg("saving balance sheets to database")
			balanceSheets.SaveToDB(ctx, conn)
			if err := balanceSheets.SaveToParquet("balance_sheet_info.parquet"); err != nil {
				log.Error().Err(err).Msg("failed to save to parquet")
			}
		} else {
			log.Error().Err(err).Msg("caught error when parsing balance sheet")
		}
	},
}

func init() {
	balanceSheetCmd.LocalFlags().IntVar(&lookback, "lookback", 30, "Number of days to lookback")
	balanceSheetCmd.LocalFlags().IntVar(&maxAssets, "max-assets", 25, "Maximum number of discovered assets to include")

	rootCmd.AddCommand(balanceSheetCmd)
}
