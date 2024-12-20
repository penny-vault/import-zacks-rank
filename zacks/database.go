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
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func AddExclusion(ticker string) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, viper.GetString("database.url"))
	if err != nil {
		log.Error().Err(err).Msg("could not connect to database in AddExclusion")
	}
	defer conn.Close(ctx)

	// lookup the composite figi
	rows, err := conn.Query(ctx, "SELECT composite_figi FROM assets WHERE ticker=$1 AND active='t' LIMIT 1", ticker)
	if err != nil {
		log.Error().Err(err).Msg("could not query database from composite_figi in AddExclusion")
	}

	var figi string
	for rows.Next() {
		if err := rows.Scan(&figi); err != nil {
			log.Error().Err(err).Msg("could not scan into figi string")
		}
	}

	// save to exclusions table
	if _, err := conn.Exec(ctx, `INSERT INTO zacks_balance_sheet_exclusions ("ticker", "composite_figi") VALUES ($1, $2)`, ticker, figi); err != nil {
		log.Error().Err(err).Str("Ticker", ticker).Str("CompositeFIGI", figi).Msg("could not save exclusion to DB")
	}
}

func SaveToDB(records []*ZacksRecord) error {
	conn, err := pgx.Connect(context.Background(), viper.GetString("database.url"))
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to database")
		return err
	}
	defer conn.Close(context.Background())

	cnt := 0
	for _, r := range records {
		if r.CompositeFigi != "" {
			_, err = conn.Exec(context.Background(),
				`INSERT INTO zacks_financials (
				"ticker",
				"composite_figi",
				"event_date",
				"in_sp500",
				"month_of_fiscal_yr_end",
				"optionable",
				"sector",
				"industry",
				"shares_outstanding_mil",
				"market_cap_mil",
				"avg_volume",
				"wk_high_52",
				"wk_low_52",
				"price_as_percent_of_52wk_hl",
				"beta",
				"percent_price_change_1wk",
				"percent_price_change_4wk",
				"percent_price_change_12wk",
				"percent_price_change_ytd",
				"relative_price_change",
				"zacks_rank",
				"zacks_rank_change_indicator",
				"zacks_industry_rank",
				"value_score",
				"growth_score",
				"momentum_score",
				"vgm_score",
				"current_avg_broker_rec",
				"num_brokers_in_rating",
				"num_rating_strong_buy_or_buy",
				"percent_rating_strong_buy_or_buy",
				"num_rating_hold",
				"num_rating_strong_sell_or_sell",
				"percent_rating_strong_sell_or_sell",
				"percent_rating_change_4wk",
				"industry_rank_of_abr",
				"rank_in_industry_of_abr",
				"change_in_avg_rec",
				"number_rating_upgrades",
				"number_rating_downgrades",
				"percent_rating_hold",
				"percent_rating_upgrades",
				"percent_rating_downgrades",
				"average_target_price",
				"earnings_esp",
				"last_eps_surprise_percent",
				"previous_eps_surprise_percent",
				"avg_eps_surprise_last_4_qtrs",
				"actual_eps_used_in_surprise_dollars_per_share",
				"last_qtr_eps",
				"last_reported_qtr_date",
				"last_yr_eps_f0_before_nri",
				"twelve_mo_trailing_eps",
				"last_reported_fiscal_yr",
				"last_eps_report_date",
				"next_eps_report_date",
				"percent_change_q0_est",
				"percent_change_q2_est",
				"percent_change_f1_est",
				"percent_change_q1_est",
				"percent_change_f2_est",
				"percent_change_lt_growth_est",
				"q0_consensus_est_last_completed_fiscal_qtr",
				"number_of_analysts_in_q0_consensus",
				"q1_consensus_est",
				"number_of_analysts_in_q1_consensus",
				"stdev_q1_q1_consensus_ratio",
				"q2_consensus_est_next_fiscal_qtr",
				"number_of_analysts_in_q2_consensus",
				"stdev_q2_q2_consensus_ratio",
				"f0_consensus_est",
				"number_of_analysts_in_f0_consensus",
				"f1_consensus_est",
				"number_of_analysts_in_f1_consensus",
				"stdev_f1_f1_consensus_ratio",
				"f2_consensus_est",
				"number_of_analysts_in_f2_consensus",
				"five_yr_hist_eps_growth",
				"long_term_growth_consensus_est",
				"percent_change_eps",
				"last_yrs_growth",
				"this_yrs_est_growth",
				"percent_ratio_of_q1_q0",
				"percent_ratio_of_q1_prior_yr_q1_actual_q",
				"sales_growth",
				"five_yr_historical_sales_growth",
				"q1_consensus_sales_est_mil",
				"f1_consensus_sales_est_mil",
				"pe_trailing_12_months",
				"pe_f1",
				"pe_f2",
				"peg_ratio",
				"price_to_cash_flow",
				"price_to_sales",
				"price_to_book",
				"current_roe_ttm",
				"current_roi_ttm",
				"roi_5_yr_avg",
				"current_roa_ttm",
				"roa_5_yr_avg",
				"market_value_to_number_analysts",
				"annual_sales_mil",
				"cost_of_goods_sold_mil",
				"ebitda_mil",
				"ebit_mil",
				"pretax_income_mil",
				"net_income_mil",
				"cash_flow_mil",
				"net_income_growth_f0_f_neg1",
				"twelve_mo_net_income_current_to_last_percent",
				"twelve_mo_net_income_current_1q_to_last_1q_percent",
				"div_yield_percent",
				"five_yr_div_yield_percent",
				"five_yr_hist_div_growth_percent",
				"dividend",
				"net_margin_percent",
				"turnover",
				"operating_margin_12_mo_percent",
				"inventory_turnover",
				"asset_utilization",
				"receivables_mil",
				"intangibles_mil",
				"inventory_mil",
				"current_assets_mil",
				"current_liabilities_mil",
				"long_term_debt_mil",
				"preferred_equity_mil",
				"common_equity_mil",
				"book_value",
				"debt_to_total_capital",
				"debt_to_equity_ratio",
				"current_ratio",
				"quick_ratio",
				"cash_ratio"
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12,
				$13,
				$14,
				$15,
				$16,
				$17,
				$18,
				$19,
				$20,
				$21,
				$22,
				$23,
				$24,
				$25,
				$26,
				$27,
				$28,
				$29,
				$30,
				$31,
				$32,
				$33,
				$34,
				$35,
				$36,
				$37,
				$38,
				$39,
				$40,
				$41,
				$42,
				$43,
				$44,
				$45,
				$46,
				$47,
				$48,
				$49,
				$50,
				$51,
				$52,
				$53,
				$54,
				$55,
				$56,
				$57,
				$58,
				$59,
				$60,
				$61,
				$62,
				$63,
				$64,
				$65,
				$66,
				$67,
				$68,
				$69,
				$70,
				$71,
				$72,
				$73,
				$74,
				$75,
				$76,
				$77,
				$78,
				$79,
				$80,
				$81,
				$82,
				$83,
				$84,
				$85,
				$86,
				$87,
				$88,
				$89,
				$90,
				$91,
				$92,
				$93,
				$94,
				$95,
				$96,
				$97,
				$98,
				$99,
				$100,
				$101,
				$102,
				$103,
				$104,
				$105,
				$106,
				$107,
				$108,
				$109,
				$110,
				$111,
				$112,
				$113,
				$114,
				$115,
				$116,
				$117,
				$118,
				$119,
				$120,
				$121,
				$122,
				$123,
				$124,
				$125,
				$126,
				$127,
				$128,
				$129,
				$130,
				$131,
				$132,
				$133,
				$134
			) ON CONFLICT ON CONSTRAINT zacks_financials_pkey
			DO UPDATE SET
				ticker = EXCLUDED.ticker,
				composite_figi = EXCLUDED.composite_figi,
				event_date = EXCLUDED.event_date,
				in_sp500 = EXCLUDED.in_sp500,
				month_of_fiscal_yr_end = EXCLUDED.month_of_fiscal_yr_end,
				optionable = EXCLUDED.optionable,
				sector = EXCLUDED.sector,
				industry = EXCLUDED.industry,
				shares_outstanding_mil = EXCLUDED.shares_outstanding_mil,
				market_cap_mil = EXCLUDED.market_cap_mil,
				avg_volume = EXCLUDED.avg_volume,
				wk_high_52 = EXCLUDED.wk_high_52,
				wk_low_52 = EXCLUDED.wk_low_52,
				price_as_percent_of_52wk_hl = EXCLUDED.price_as_percent_of_52wk_hl,
				beta = EXCLUDED.beta,
				percent_price_change_1wk = EXCLUDED.percent_price_change_1wk,
				percent_price_change_4wk = EXCLUDED.percent_price_change_4wk,
				percent_price_change_12wk = EXCLUDED.percent_price_change_12wk,
				percent_price_change_ytd = EXCLUDED.percent_price_change_ytd,
				relative_price_change = EXCLUDED.relative_price_change,
				zacks_rank = EXCLUDED.zacks_rank,
				zacks_rank_change_indicator = EXCLUDED.zacks_rank_change_indicator,
				zacks_industry_rank = EXCLUDED.zacks_industry_rank,
				value_score = EXCLUDED.value_score,
				growth_score = EXCLUDED.growth_score,
				momentum_score = EXCLUDED.momentum_score,
				vgm_score = EXCLUDED.vgm_score,
				current_avg_broker_rec = EXCLUDED.current_avg_broker_rec,
				num_brokers_in_rating = EXCLUDED.num_brokers_in_rating,
				num_rating_strong_buy_or_buy = EXCLUDED.num_rating_strong_buy_or_buy,
				percent_rating_strong_buy_or_buy = EXCLUDED.percent_rating_strong_buy_or_buy,
				num_rating_hold = EXCLUDED.num_rating_hold,
				num_rating_strong_sell_or_sell = EXCLUDED.num_rating_strong_sell_or_sell,
				percent_rating_strong_sell_or_sell = EXCLUDED.percent_rating_strong_sell_or_sell,
				percent_rating_change_4wk = EXCLUDED.percent_rating_change_4wk,
				industry_rank_of_abr = EXCLUDED.industry_rank_of_abr,
				rank_in_industry_of_abr = EXCLUDED.rank_in_industry_of_abr,
				change_in_avg_rec = EXCLUDED.change_in_avg_rec,
				number_rating_upgrades = EXCLUDED.number_rating_upgrades,
				number_rating_downgrades = EXCLUDED.number_rating_downgrades,
				percent_rating_hold = EXCLUDED.percent_rating_hold,
				percent_rating_upgrades = EXCLUDED.percent_rating_upgrades,
				percent_rating_downgrades = EXCLUDED.percent_rating_downgrades,
				average_target_price = EXCLUDED.average_target_price,
				earnings_esp = EXCLUDED.earnings_esp,
				last_eps_surprise_percent = EXCLUDED.last_eps_surprise_percent,
				previous_eps_surprise_percent = EXCLUDED.previous_eps_surprise_percent,
				avg_eps_surprise_last_4_qtrs = EXCLUDED.avg_eps_surprise_last_4_qtrs,
				actual_eps_used_in_surprise_dollars_per_share = EXCLUDED.actual_eps_used_in_surprise_dollars_per_share,
				last_qtr_eps = EXCLUDED.last_qtr_eps,
				last_reported_qtr_date = EXCLUDED.last_reported_qtr_date,
				last_yr_eps_f0_before_nri = EXCLUDED.last_yr_eps_f0_before_nri,
				twelve_mo_trailing_eps = EXCLUDED.twelve_mo_trailing_eps,
				last_reported_fiscal_yr = EXCLUDED.last_reported_fiscal_yr,
				last_eps_report_date = EXCLUDED.last_eps_report_date,
				next_eps_report_date = EXCLUDED.next_eps_report_date,
				percent_change_q0_est = EXCLUDED.percent_change_q0_est,
				percent_change_q2_est = EXCLUDED.percent_change_q2_est,
				percent_change_f1_est = EXCLUDED.percent_change_f1_est,
				percent_change_q1_est = EXCLUDED.percent_change_q1_est,
				percent_change_f2_est = EXCLUDED.percent_change_f2_est,
				percent_change_lt_growth_est = EXCLUDED.percent_change_lt_growth_est,
				q0_consensus_est_last_completed_fiscal_qtr = EXCLUDED.q0_consensus_est_last_completed_fiscal_qtr,
				number_of_analysts_in_q0_consensus = EXCLUDED.number_of_analysts_in_q0_consensus,
				q1_consensus_est = EXCLUDED.q1_consensus_est,
				number_of_analysts_in_q1_consensus = EXCLUDED.number_of_analysts_in_q1_consensus,
				stdev_q1_q1_consensus_ratio = EXCLUDED.stdev_q1_q1_consensus_ratio,
				q2_consensus_est_next_fiscal_qtr = EXCLUDED.q2_consensus_est_next_fiscal_qtr,
				number_of_analysts_in_q2_consensus = EXCLUDED.number_of_analysts_in_q2_consensus,
				stdev_q2_q2_consensus_ratio = EXCLUDED.stdev_q2_q2_consensus_ratio,
				f0_consensus_est = EXCLUDED.f0_consensus_est,
				number_of_analysts_in_f0_consensus = EXCLUDED.number_of_analysts_in_f0_consensus,
				f1_consensus_est = EXCLUDED.f1_consensus_est,
				number_of_analysts_in_f1_consensus = EXCLUDED.number_of_analysts_in_f1_consensus,
				stdev_f1_f1_consensus_ratio = EXCLUDED.stdev_f1_f1_consensus_ratio,
				f2_consensus_est = EXCLUDED.f2_consensus_est,
				number_of_analysts_in_f2_consensus = EXCLUDED.number_of_analysts_in_f2_consensus,
				five_yr_hist_eps_growth = EXCLUDED.five_yr_hist_eps_growth,
				long_term_growth_consensus_est = EXCLUDED.long_term_growth_consensus_est,
				percent_change_eps = EXCLUDED.percent_change_eps,
				last_yrs_growth = EXCLUDED.last_yrs_growth,
				this_yrs_est_growth = EXCLUDED.this_yrs_est_growth,
				percent_ratio_of_q1_q0 = EXCLUDED.percent_ratio_of_q1_q0,
				percent_ratio_of_q1_prior_yr_q1_actual_q = EXCLUDED.percent_ratio_of_q1_prior_yr_q1_actual_q,
				sales_growth = EXCLUDED.sales_growth,
				five_yr_historical_sales_growth = EXCLUDED.five_yr_historical_sales_growth,
				q1_consensus_sales_est_mil = EXCLUDED.q1_consensus_sales_est_mil,
				f1_consensus_sales_est_mil = EXCLUDED.f1_consensus_sales_est_mil,
				pe_trailing_12_months = EXCLUDED.pe_trailing_12_months,
				pe_f1 = EXCLUDED.pe_f1,
				pe_f2 = EXCLUDED.pe_f2,
				peg_ratio = EXCLUDED.peg_ratio,
				price_to_cash_flow = EXCLUDED.price_to_cash_flow,
				price_to_sales = EXCLUDED.price_to_sales,
				price_to_book = EXCLUDED.price_to_book,
				current_roe_ttm = EXCLUDED.current_roe_ttm,
				current_roi_ttm = EXCLUDED.current_roi_ttm,
				roi_5_yr_avg = EXCLUDED.roi_5_yr_avg,
				current_roa_ttm = EXCLUDED.current_roa_ttm,
				roa_5_yr_avg = EXCLUDED.roa_5_yr_avg,
				market_value_to_number_analysts = EXCLUDED.market_value_to_number_analysts,
				annual_sales_mil = EXCLUDED.annual_sales_mil,
				cost_of_goods_sold_mil = EXCLUDED.cost_of_goods_sold_mil,
				ebitda_mil = EXCLUDED.ebitda_mil,
				ebit_mil = EXCLUDED.ebit_mil,
				pretax_income_mil = EXCLUDED.pretax_income_mil,
				net_income_mil = EXCLUDED.net_income_mil,
				cash_flow_mil = EXCLUDED.cash_flow_mil,
				net_income_growth_f0_f_neg1 = EXCLUDED.net_income_growth_f0_f_neg1,
				twelve_mo_net_income_current_to_last_percent = EXCLUDED.twelve_mo_net_income_current_to_last_percent,
				twelve_mo_net_income_current_1q_to_last_1q_percent = EXCLUDED.twelve_mo_net_income_current_1q_to_last_1q_percent,
				div_yield_percent = EXCLUDED.div_yield_percent,
				five_yr_div_yield_percent = EXCLUDED.five_yr_div_yield_percent,
				five_yr_hist_div_growth_percent = EXCLUDED.five_yr_hist_div_growth_percent,
				dividend = EXCLUDED.dividend,
				net_margin_percent = EXCLUDED.net_margin_percent,
				turnover = EXCLUDED.turnover,
				operating_margin_12_mo_percent = EXCLUDED.operating_margin_12_mo_percent,
				inventory_turnover = EXCLUDED.inventory_turnover,
				asset_utilization = EXCLUDED.asset_utilization,
				receivables_mil = EXCLUDED.receivables_mil,
				intangibles_mil = EXCLUDED.intangibles_mil,
				inventory_mil = EXCLUDED.inventory_mil,
				current_assets_mil = EXCLUDED.current_assets_mil,
				current_liabilities_mil = EXCLUDED.current_liabilities_mil,
				long_term_debt_mil = EXCLUDED.long_term_debt_mil,
				preferred_equity_mil = EXCLUDED.preferred_equity_mil,
				common_equity_mil = EXCLUDED.common_equity_mil,
				book_value = EXCLUDED.book_value,
				debt_to_total_capital = EXCLUDED.debt_to_total_capital,
				debt_to_equity_ratio = EXCLUDED.debt_to_equity_ratio,
				current_ratio = EXCLUDED.current_ratio,
				quick_ratio = EXCLUDED.quick_ratio,
				cash_ratio = EXCLUDED.cash_ratio
			`,
				r.Ticker, r.CompositeFigi, r.EventDate,
				r.InSp500, r.MonthOfFiscalYrEnd, r.Optionable, r.Sector,
				r.Industry, r.SharesOutstandingMil, r.MarketCapMil,
				r.AvgVolume, r.WkHigh52, r.WkLow52,
				r.PriceAsPercentOf52wkHighLow, r.Beta, r.PercentPriceChange1Wk,
				r.PercentPriceChange4Wk, r.PercentPriceChange12Wk,
				r.PercentPriceChangeYtd, r.RelativePriceChange, r.ZacksRank,
				r.ZacksRankChangeIndicator, r.ZacksIndustryRank, r.ValueScore,
				r.GrowthScore, r.MomentumScore, r.VgmScore, r.CurrentAvgBrokerRec,
				r.NumBrokersInRating, r.NumRatingStrongBuyOrBuy,
				r.PercentRatingStrongBuyOrBuy, r.NumRatingHold,
				r.NumRatingStrongSellOrSell, r.PercentRatingStrongSellOrSell,
				r.PercentRatingChange4Wk, r.IndustryRankOfAbr,
				r.RankInIndustryOfAbr, r.ChangeInAvgRec, r.NumberRatingUpgrades,
				r.NumberRatingDowngrades, r.PercentRatingHold,
				r.PercentRatingUpgrades, r.PercentRatingDowngrades,
				r.AverageTargetPrice, r.EarningsEsp, r.LastEpsSurprisePercent,
				r.PreviousEpsSurprisePercent, r.AvgEpsSurpriseLast4Qtrs,
				r.ActualEpsUsedInSurpriseDollarsPerShare, r.LastQtrEps,
				r.LastReportedQtrDate, r.LastYrEpsF0BeforeNri,
				r.TwelveMoTrailingEps, r.LastReportedFiscalYr,
				r.LastEpsReportDate, r.NextEpsReportDate, r.PercentChangeQ0Est,
				r.PercentChangeQ2Est, r.PercentChangeF1Est, r.PercentChangeQ1Est,
				r.PercentChangeF2Est, r.PercentChangeLtGrowthEst,
				r.Q0ConsensusEstLastCompletedFiscalQtr,
				r.NumberOfAnalystsInQ0Consensus, r.Q1ConsensusEst,
				r.NumberOfAnalystsInQ1Consensus, r.StdevQ1Q1ConsensusRatio,
				r.Q2ConsensusEstNextFiscalQtr, r.NumberOfAnalystsInQ2Consensus,
				r.StdevQ2Q2ConsensusRatio, r.F0ConsensusEst,
				r.NumberOfAnalystsInF0Consensus, r.F1ConsensusEst,
				r.NumberOfAnalystsInF1Consensus, r.StdevF1F1ConsensusRatio,
				r.F2ConsensusEst, r.NumberOfAnalystsInF2Consensus,
				r.FiveYrHistEpsGrowth, r.LongTermGrowthConsensusEst,
				r.PercentChangeEps, r.LastYrsGrowth, r.ThisYrsEstGrowth,
				r.PercentRatioOfQ1Q0, r.PercentRatioOfQ1PriorYrQ1ActualQ,
				r.SalesGrowth, r.FiveYrHistoricalSalesGrowth,
				r.Q1ConsensusSalesEstMil, r.F1ConsensusSalesEstMil,
				r.PeTrailing12Months, r.PeF1, r.PeF2, r.PegRatio,
				r.PriceToCashFlow, r.PriceToSales, r.PriceToBook,
				r.CurrentRoeTtm, r.CurrentRoiTtm, r.Roi5YrAvg, r.CurrentRoaTtm,
				r.Roa5YrAvg, r.MarketValueToNumberAnalysts, r.AnnualSalesMil,
				r.CostOfGoodsSoldMil, r.EbitdaMil, r.EbitMil, r.PretaxIncomeMil,
				r.NetIncomeMil, r.CashFlowMil, r.NetIncomeGrowthF0FNeg1,
				r.TwelveMoNetIncomeCurrentToLastPercent,
				r.TwelveMoNetIncomeCurrent1qToLast1qPercent, r.DivYieldPercent,
				r.FiveYrDivYieldPercent, r.FiveYrHistDivGrowthPercent,
				r.Dividend, r.NetMarginPercent, r.Turnover,
				r.OperatingMargin12MoPercent, r.InventoryTurnover,
				r.AssetUtilization, r.ReceivablesMil, r.IntangiblesMil,
				r.InventoryMil, r.CurrentAssetsMil,
				r.CurrentLiabilitiesMil, r.LongTermDebtMil,
				r.PreferredEquityMil, r.CommonEquityMil, r.BookValue,
				r.DebtToTotalCapital, r.DebtToEquityRatio,
				r.CurrentRatio, r.QuickRatio, r.CashRatio)
			if err != nil {
				log.Warn().Err(err).Str("CompositeFigi", r.CompositeFigi).Str("Ticker", r.Ticker).Int("ZacksRank", r.ZacksRank).Msg("insert into db failed")
				return err
			} else {
				cnt++
			}
		}
	}

	log.Info().Int("NumRecords", cnt).Msg("records saved to DB")
	return nil
}

func (balanceSheetList BalanceSheetList) SaveToDB(ctx context.Context, conn *pgx.Conn) {
	// build a list of all active records that have composite figi's
	tickerMap := make(map[string]*Ticker)

	// build figi map
	rows, err := conn.Query(context.Background(), "SELECT ticker, name, composite_figi FROM assets WHERE active='t' AND composite_figi IS NOT NULL")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve tickers from database")
	}

	for rows.Next() {
		var ticker Ticker
		err := rows.Scan(&ticker.Ticker, &ticker.CompanyName, &ticker.CompositeFigi)
		if err != nil {
			log.Error().Err(err).Msg("Failed to retrieve ticker row from database")
		}
		tickerMap[ticker.Ticker] = &ticker
	}

	// save each balance sheet to database
	for _, r := range balanceSheetList {
		if ticker, ok := tickerMap[r.Ticker]; ok {
			r.CompositeFigi = ticker.CompositeFigi
			if _, err := conn.Exec(ctx, "UPDATE fundamentals SET curr_assets=$1, curr_liabilities=$2, working_capital=$3 WHERE composite_figi=$4 AND calendar_date=$5 AND dim=$6", r.TotalCurrentAssets, r.TotalCurrentLiabilities, r.TotalCurrentAssets-r.TotalCurrentLiabilities, r.CompositeFigi, r.CalendarDate, r.Dimension); err != nil {
				log.Error().Err(err).Str("Ticker", r.Ticker).Str("CompositeFIGI", r.CompositeFigi).Msg("error updating database")
			}
		}
	}
}
