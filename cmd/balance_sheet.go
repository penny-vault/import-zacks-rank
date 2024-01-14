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
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/penny-vault/import-zacks-rank/zacks"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			if rows, err := conn.Query(ctx, "SELECT distinct ticker FROM fundamentals WHERE event_date > $1 AND working_capital = 'NaN'::float8", time.Now().Add(-365*24*time.Hour)); err != nil {
				log.Fatal().Err(err).Msg("error querying database for tickers without working_capital")
			} else {
				var ticker string
				if err := rows.Scan(&ticker); err != nil {
					log.Fatal().Err(err).Msg("unable to scan query value into ticker string")
				}

				args = append(args, ticker)
			}
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
	rootCmd.AddCommand(balanceSheetCmd)
}
