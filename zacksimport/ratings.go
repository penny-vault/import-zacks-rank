package zacksimport

import (
	"context"
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/viper"

	"github.com/jackc/pgx/v4"

	"github.com/rs/zerolog/log"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type Ticker struct {
	CompanyName   string
	TickerId      int
	Ticker        string
	CompositeFigi string
}

type ZacksRecord struct {
	CompanyName                               string  `csv:"Company Name" json:"company_name" parquet:"name=company_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ticker                                    string  `csv:"Ticker" json:"ticker" parquet:"name=ticker, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"ticker"`
	CompositeFigi                             string  `csv:"-" json:"composite_figi" parquet:"name=composite_figi, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"composite_figi"`
	Exchange                                  string  `csv:"Exchange" json:"exchange" parquet:"name=exchange, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"exchange"`
	EventDate                                 string  `csv:"-" json:"event_date" parquet:"name=event_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"event_date"`
	InSp500                                   bool    `csv:"S&P 500 - ETF" json:"in_sp500" parquet:"name=in_sp500, type=BOOLEAN" db:"in_sp500"`
	LastClose                                 float64 `csv:"Last Close" json:"last_close" parquet:"name=last_close, type=DOUBLE"`
	MonthOfFiscalYrEnd                        int     `csv:"Month of Fiscal Yr End" json:"month_of_fiscal_yr_end" parquet:"name=month_of_fiscal_yr_end, type=INT32" db:"month_of_fiscal_yr_end"`
	Optionable                                bool    `csv:"Optionable" json:"optionable" parquet:"name=optionable, type=BOOLEAN" db:"optionable"`
	Sector                                    string  `csv:"Sector" json:"sector" parquet:"name=sector, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"sector"`
	Industry                                  string  `csv:"Industry" json:"industry" parquet:"name=industry, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"industry"`
	SharesOutstandingMil                      float64 `csv:"Shares Outstanding (mil)" json:"shares_outstanding_mil" parquet:"name=shares_outstanding_mil, type=DOUBLE" db:"shares_outstanding_mil"`
	MarketCapMil                              float64 `csv:"Market Cap (mil)" json:"market_cap_mil" parquet:"name=market_cap_mil, type=DOUBLE" db:"market_cap_mil"`
	AvgVolume                                 int64   `csv:"Avg Volume" json:"avg_volume" parquet:"name=avg_volume, type=INT64" db:"avg_volume"`
	WkHigh52                                  float64 `csv:"52 Week High" json:"wk_high_52" parquet:"name=wk_high_52, type=DOUBLE" db:"wk_high_52"`
	WkLow52                                   float64 `csv:"52 Week Low" json:"wk_low_52" parquet:"name=wk_low_52, type=DOUBLE" db:"wk_low_52"`
	PriceAsPercentOf52wkHighLow               float32 `csv:"Price as a % of 52 Wk H-L Range" json:"price_as_percent_of_52wk_hl" parquet:"name=price_as_percent_of_52wk_hl, type=FLOAT" db:"price_as_percent_of_52wk_hl"`
	Beta                                      float32 `csv:"Beta" json:"beta" parquet:"name=beta, type=FLOAT" db:"beta"`
	PercentPriceChange1Wk                     float32 `csv:"% Price Change (1 Week)" json:"percent_price_change_1wk" parquet:"name=percent_price_change_1wk, type=FLOAT" db:"percent_price_change_1wk"`
	PercentPriceChange4Wk                     float32 `csv:"% Price Change (4 Weeks)" json:"percent_price_change_4wk" parquet:"name=percent_price_change_4wk, type=FLOAT" db:"percent_price_change_4wk"`
	PercentPriceChange12Wk                    float32 `csv:"% Price Change (12 Weeks)" json:"percent_price_change_12wk" parquet:"name=percent_price_change_12wk, type=FLOAT" db:"percent_price_change_12wk"`
	PercentPriceChangeYtd                     float32 `csv:"% Price Change (YTD)" json:"percent_price_change_ytd" parquet:"name=percent_price_change_ytd, type=FLOAT" db:"percent_price_change_ytd"`
	RelativePriceChange                       float32 `csv:"Relative Price Change" json:"relative_price_change" parquet:"name=relative_price_change, type=FLOAT" db:"relative_price_change"`
	ZacksRank                                 int     `csv:"Zacks Rank" json:"zacks_rank" parquet:"name=zacks_rank, type=INT32" db:"zacks_rank"`
	ZacksRankChangeIndicator                  int     `csv:"Zacks Rank Change Indicator" json:"zacks_rank_change_indicator" parquet:"name=zacks_rank_change_indicator, type=INT32" db:"zacks_rank_change_indicator"`
	ZacksIndustryRank                         int     `csv:"Zacks Industry Rank" json:"zacks_industry_rank" parquet:"name=zacks_industry_rank, type=INT32" db:"zacks_industry_rank"`
	ValueScore                                string  `csv:"Value Score" json:"value_score" parquet:"name=value_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"value_score"`
	GrowthScore                               string  `csv:"Growth Score" json:"growth_score" parquet:"name=growth_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"growth_score"`
	MomentumScore                             string  `csv:"Momentum Score" json:"momentum_score" parquet:"name=momentum_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"momentum_score"`
	VgmScore                                  string  `csv:"VGM Score" json:"vgm_score" parquet:"name=vgm_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"vgm_score"`
	CurrentAvgBrokerRec                       float32 `csv:"Current Avg Broker Rec" json:"current_avg_broker_rec" parquet:"name=current_avg_broker_rec, type=FLOAT" db:"current_avg_broker_rec"`
	NumBrokersInRating                        int     `csv:"# of Brokers in Rating" json:"num_brokers_in_rating" parquet:"name=num_brokers_in_rating, type=INT32" db:"num_brokers_in_rating"`
	NumRatingStrongBuyOrBuy                   int     `csv:"# Rating Strong Buy or Buy" json:"num_rating_strong_buy_or_buy" parquet:"name=num_rating_strong_buy_or_buy, type=INT32" db:"num_rating_strong_buy_or_buy"`
	PercentRatingStrongBuyOrBuy               float32 `csv:"% Rating Strong Buy or Buy" json:"percent_rating_strong_buy_or_buy" parquet:"name=percent_rating_strong_buy_or_buy, type=FLOAT" db:"percent_rating_strong_buy_or_buy"`
	NumRatingHold                             int     `csv:"# Rating Hold" json:"num_rating_hold" parquet:"name=num_rating_hold, type=INT32" db:"num_rating_hold"`
	NumRatingStrongSellOrSell                 int     `csv:"# Rating Strong Sell or Sell" json:"num_rating_strong_sell_or_sell" parquet:"name=num_rating_strong_sell_or_sell, type=INT32" db:"num_rating_strong_sell_or_sell"`
	PercentRatingStrongSellOrSell             float32 `csv:"% Rating Strong Sell or Sell" json:"percent_rating_strong_sell_or_sell" parquet:"name=percent_rating_strong_sell_or_sell, type=FLOAT" db:"percent_rating_strong_sell_or_sell"`
	PercentRatingChange4Wk                    float32 `csv:"% Rating Change - 4 Weeks" json:"percent_rating_change_4wk" parquet:"name=percent_rating_change_4wk, type=FLOAT" db:"percent_rating_change_4wk"`
	IndustryRankOfAbr                         int     `csv:"Industry Rank (of ABR)" json:"industry_rank_of_abr" parquet:"name=industry_rank_of_abr, type=INT32" db:"industry_rank_of_abr"`
	RankInIndustryOfAbr                       int     `csv:"Rank in Industry (of ABR)" json:"rank_in_industry_of_abr" parquet:"name=rank_in_industry_of_abr, type=INT32" db:"rank_in_industry_of_abr"`
	ChangeInAvgRec                            float32 `csv:"Change in Avg Rec " json:"change_in_avg_rec" parquet:"name=change_in_avg_rec, type=FLOAT" db:"change_in_avg_rec"`
	NumberRatingUpgrades                      int     `csv:"# Rating Upgrades" json:"number_rating_upgrades" parquet:"name=number_rating_upgrades, type=INT32" db:"number_rating_upgrades"`
	NumberRatingDowngrades                    int     `csv:"# Rating Downgrades " json:"number_rating_downgrades" parquet:"name=number_rating_downgrades, type=INT32" db:"number_rating_downgrades"`
	PercentRatingHold                         float32 `csv:"% Rating Hold" json:"percent_rating_hold" parquet:"name=percent_rating_hold, type=FLOAT" db:"percent_rating_hold"`
	PercentRatingUpgrades                     float32 `csv:"% Rating Upgrades " json:"percent_rating_upgrades" parquet:"name=percent_rating_upgrades, type=FLOAT" db:"percent_rating_upgrades"`
	PercentRatingDowngrades                   float32 `csv:"% Rating Downgrades " json:"percent_rating_downgrades" parquet:"name=percent_rating_downgrades, type=FLOAT" db:"percent_rating_downgrades"`
	AverageTargetPrice                        float64 `csv:"Average Target Price" json:"average_target_price" parquet:"name=average_target_price, type=DOUBLE" db:"average_target_price"`
	EarningsEsp                               float32 `csv:"Earnings ESP" json:"earnings_esp" parquet:"name=earnings_esp, type=FLOAT" db:"earnings_esp"`
	LastEpsSurprisePercent                    float32 `csv:"Last EPS Surprise (%)" json:"last_eps_surprise_percent" parquet:"name=last_eps_surprise_percent, type=FLOAT" db:"last_eps_surprise_percent"`
	PreviousEpsSurprisePercent                float32 `csv:"Previous EPS Surprise (%)" json:"previous_eps_surprise_percent" parquet:"name=previous_eps_surprise_percent, type=FLOAT" db:"previous_eps_surprise_percent"`
	AvgEpsSurpriseLast4Qtrs                   float32 `csv:"Avg EPS Surprise (Last 4 Qtrs)" json:"avg_eps_surprise_last_4_qtrs" parquet:"name=avg_eps_surprise_last_4_qtrs, type=FLOAT" db:"avg_eps_surprise_last_4_qtrs"`
	ActualEpsUsedInSurpriseDollarsPerShare    float32 `csv:"Actual EPS used in Surprise ($/sh)" json:"actual_eps_used_in_surprise_dollars_per_share" parquet:"name=actual_eps_used_in_surprise_dollars_per_share, type=FLOAT" db:"actual_eps_used_in_surprise_dollars_per_share"`
	LastQtrEps                                float32 `csv:"Last Qtr EPS" json:"last_qtr_eps" parquet:"name=last_qtr_eps, type=FLOAT" db:"last_qtr_eps"`
	LastReportedQtrDate                       string  `csv:"Last Reported Qtr (yyyymm)" json:"last_reported_qtr_date" parquet:"name=last_reported_qtr_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"last_reported_qtr_date"`
	LastYrEpsF0BeforeNri                      float32 `csv:"Last Yr's EPS (F0) Before NRI" json:"last_yr_eps_f0_before_nri" parquet:"name=last_yr_eps_f0_before_nri, type=FLOAT" db:"last_yr_eps_f0_before_nri"`
	TwelveMoTrailingEps                       float32 `csv:"12 Mo Trailing EPS" json:"twelve_mo_trailing_eps" parquet:"name=twelve_mo_trailing_eps, type=FLOAT" db:"twelve_mo_trailing_eps"`
	LastReportedFiscalYr                      string  `csv:"Last Reported Fiscal Yr  (yyyymm)" json:"last_reported_fiscal_yr" parquet:"name=last_reported_fiscal_yr, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"last_reported_fiscal_yr"`
	LastEpsReportDate                         string  `csv:"Last EPS Report Date (yyyymmdd)" json:"last_eps_report_date" parquet:"name=last_eps_report_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"last_eps_report_date"`
	NextEpsReportDate                         string  `csv:"Next EPS Report Date  (yyyymmdd)" json:"next_eps_report_date" parquet:"name=next_eps_report_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"next_eps_report_date"`
	PercentChangeQ0Est                        float32 `csv:"% Change Q0 Est. (4 weeks)" json:"percent_change_q0_est" parquet:"name=percent_change_q0_est, type=FLOAT" db:"percent_change_q0_est"`
	PercentChangeQ2Est                        float32 `csv:"% Change Q2 Est. (4 weeks)" json:"percent_change_q2_est" parquet:"name=percent_change_q2_est, type=FLOAT" db:"percent_change_q2_est"`
	PercentChangeF1Est                        float32 `csv:"% Change F1 Est. (4 weeks)" json:"percent_change_f1_est" parquet:"name=percent_change_f1_est, type=FLOAT" db:"percent_change_f1_est"`
	PercentChangeQ1Est                        float32 `csv:"% Change Q1 Est. (4 weeks)" json:"percent_change_q1_est" parquet:"name=percent_change_q1_est, type=FLOAT" db:"percent_change_q1_est"`
	PercentChangeF2Est                        float32 `csv:"% Change F2 Est. (4 weeks)" json:"percent_change_f2_est" parquet:"name=percent_change_f2_est, type=FLOAT" db:"percent_change_f2_est"`
	PercentChangeLtGrowthEst                  float32 `csv:"% Change LT Growth Est. (4 weeks)" json:"percent_change_lt_growth_est" parquet:"name=percent_change_lt_growth_est, type=FLOAT" db:"percent_change_lt_growth_est"`
	Q0ConsensusEstLastCompletedFiscalQtr      float32 `csv:"Q0 Consensus Est. (last completed fiscal Qtr)" json:"q0_consensus_est_last_completed_fiscal_qtr" parquet:"name=q0_consensus_est_last_completed_fiscal_qtr, type=FLOAT" db:"q0_consensus_est_last_completed_fiscal_qtr"`
	NumberOfAnalystsInQ0Consensus             int     `csv:"# of Analysts in Q0 Consensus" json:"number_of_analysts_in_q0_consensus" parquet:"name=number_of_analysts_in_q0_consensus, type=INT32" db:"number_of_analysts_in_q0_consensus"`
	Q1ConsensusEst                            float32 `csv:"Q1 Consensus Est. " json:"q1_consensus_est" parquet:"name=q1_consensus_est, type=FLOAT" db:"q1_consensus_est"`
	NumberOfAnalystsInQ1Consensus             int     `csv:"# of Analysts in Q1 Consensus" json:"number_of_analysts_in_q1_consensus" parquet:"name=number_of_analysts_in_q1_consensus, type=INT32" db:"number_of_analysts_in_q1_consensus"`
	StdevQ1Q1ConsensusRatio                   float32 `csv:"St. Dev. Q1 / Q1 Consensus" json:"stdev_q1_q1_consensus_ratio" parquet:"name=stdev_q1_q1_consensus_ratio, type=FLOAT" db:"stdev_q1_q1_consensus_ratio"`
	Q2ConsensusEstNextFiscalQtr               float32 `csv:"Q2 Consensus Est. (next fiscal Qtr)" json:"q2_consensus_est_next_fiscal_qtr" parquet:"name=q2_consensus_est_next_fiscal_qtr, type=FLOAT" db:"q2_consensus_est_next_fiscal_qtr"`
	NumberOfAnalystsInQ2Consensus             int     `csv:"# of Analysts in Q2 Consensus" json:"number_of_analysts_in_q2_consensus" parquet:"name=number_of_analysts_in_q2_consensus, type=INT32" db:"number_of_analysts_in_q2_consensus"`
	StdevQ2Q2ConsensusRatio                   float32 `csv:"St. Dev. Q2 / Q2 Consensus" json:"stdev_q2_q2_consensus_ratio" parquet:"name=stdev_q2_q2_consensus_ratio, type=FLOAT" db:"stdev_q2_q2_consensus_ratio"`
	F0ConsensusEst                            float32 `csv:"F0 Consensus Est." json:"f0_consensus_est" parquet:"name=f0_consensus_est, type=FLOAT" db:"f0_consensus_est"`
	NumberOfAnalystsInF0Consensus             float32 `csv:"# of Analysts in F0 Consensus" json:"number_of_analysts_in_f0_consensus" parquet:"name=number_of_analysts_in_f0_consensus, type=FLOAT" db:"number_of_analysts_in_f0_consensus"`
	F1ConsensusEst                            float32 `csv:"F1 Consensus Est." json:"f1_consensus_est" parquet:"name=f1_consensus_est, type=FLOAT" db:"f1_consensus_est"`
	NumberOfAnalystsInF1Consensus             int     `csv:"# of Analysts in F1 Consensus" json:"number_of_analysts_in_f1_consensus" parquet:"name=number_of_analysts_in_f1_consensus, type=INT32" db:"number_of_analysts_in_f1_consensus"`
	StdevF1F1ConsensusRatio                   float32 `csv:"St. Dev. F1 / F1 Consensus" json:"stdev_f1_f1_consensus_ratio" parquet:"name=stdev_f1_f1_consensus_ratio, type=FLOAT" db:"stdev_f1_f1_consensus_ratio"`
	F2ConsensusEst                            float32 `csv:"F2 Consensus Est." json:"f2_consensus_est" parquet:"name=f2_consensus_est, type=FLOAT" db:"f2_consensus_est"`
	NumberOfAnalystsInF2Consensus             int     `csv:"# of Analysts in F2 Consensus" json:"number_of_analysts_in_f2_consensus" parquet:"name=number_of_analysts_in_f2_consensus, type=INT32" db:"number_of_analysts_in_f2_consensus"`
	FiveYrHistEpsGrowth                       float32 `csv:"5 Yr. Hist. EPS Growth" json:"five_yr_hist_eps_growth" parquet:"name=five_yr_hist_eps_growth, type=FLOAT" db:"five_yr_hist_eps_growth"`
	LongTermGrowthConsensusEst                float32 `csv:"Long-Term Growth Consensus Est." json:"long_term_growth_consensus_est" parquet:"name=long_term_growth_consensus_est, type=FLOAT" db:"long_term_growth_consensus_est"`
	PercentChangeEps                          float32 `csv:"% Change EPS (F(-1)/F(-2))" json:"percent_change_eps" parquet:"name=percent_change_eps, type=FLOAT" db:"percent_change_eps"`
	LastYrsGrowth                             float32 `csv:"Last Yrs Growth (F[0] / F [-1])" json:"last_yrs_growth" parquet:"name=last_yrs_growth, type=FLOAT" db:"last_yrs_growth"`
	ThisYrsEstGrowth                          float32 `csv:"This Yr's Est.d Growth (F(1)/F(0))" json:"this_yrs_est_growth" parquet:"name=this_yrs_est_growth, type=FLOAT" db:"this_yrs_est_growth"`
	PercentRatioOfQ1Q0                        float32 `csv:"% Ratio of Q1/Q0" json:"percent_ratio_of_q1_q0" parquet:"name=percent_ratio_of_q1_q0, type=FLOAT" db:"percent_ratio_of_q1_q0"`
	PercentRatioOfQ1PriorYrQ1ActualQ          float32 `csv:"% Ratio of Q1/prior Yr Q1 Actual Q(-3)" json:"percent_ratio_of_q1_prior_yr_q1_actual_q" parquet:"name=percent_ratio_of_q1_prior_yr_q1_actual_q, type=FLOAT" db:"percent_ratio_of_q1_prior_yr_q1_actual_q"`
	SalesGrowth                               float32 `csv:"Sales Growth F(0)/F(-1)" json:"sales_growth" parquet:"name=sales_growth, type=FLOAT" db:"sales_growth"`
	FiveYrHistoricalSalesGrowth               float32 `csv:"5 Yr Historical Sales Growth" json:"five_yr_historical_sales_growth" parquet:"name=five_yr_historical_sales_growth, type=FLOAT" db:"five_yr_historical_sales_growth"`
	Q1ConsensusSalesEstMil                    float32 `csv:"Q(1) Consensus Sales Est. ($mil)" json:"q1_consensus_sales_est_mil" parquet:"name=q1_consensus_sales_est_mil, type=FLOAT" db:"q1_consensus_sales_est_mil"`
	F1ConsensusSalesEstMil                    float32 `csv:"F(1) Consensus Sales Est. ($mil)" json:"f1_consensus_sales_est_mil" parquet:"name=f1_consensus_sales_est_mil, type=FLOAT" db:"f1_consensus_sales_est_mil"`
	PeTrailing12Months                        float32 `csv:"P/E (Trailing 12 Months)" json:"pe_trailing_12_months" parquet:"name=pe_trailing_12_months, type=FLOAT" db:"pe_trailing_12_months"`
	PeF1                                      float32 `csv:"P/E (F1)" json:"pe_f1" parquet:"name=pe_f1, type=FLOAT" db:"pe_f1"`
	PeF2                                      float32 `csv:"P/E (F2)" json:"pe_f2" parquet:"name=pe_f2, type=FLOAT" db:"pe_f2"`
	PegRatio                                  float32 `csv:"PEG Ratio" json:"peg_ratio" parquet:"name=peg_ratio, type=FLOAT" db:"peg_ratio"`
	PriceToCashFlow                           float32 `csv:"Price/Cash Flow" json:"price_to_cash_flow" parquet:"name=price_to_cash_flow, type=FLOAT" db:"price_to_cash_flow"`
	PriceToSales                              float32 `csv:"Price/Sales" json:"price_to_sales" parquet:"name=price_to_sales, type=FLOAT" db:"price_to_sales"`
	PriceToBook                               float32 `csv:"Price/Book" json:"price_to_book" parquet:"name=price_to_book, type=FLOAT" db:"price_to_book"`
	CurrentRoeTtm                             float32 `csv:"Current ROE (TTM)" json:"current_roe_ttm" parquet:"name=current_roe_ttm, type=FLOAT" db:"current_roe_ttm"`
	CurrentRoiTtm                             float32 `csv:"Current ROI (TTM)" json:"current_roi_ttm" parquet:"name=current_roi_ttm, type=FLOAT" db:"current_roi_ttm"`
	Roi5YrAvg                                 float32 `csv:"ROI (5 Yr Avg)" json:"roi_5_yr_avg" parquet:"name=roi_5_yr_avg, type=FLOAT" db:"roi_5_yr_avg"`
	CurrentRoaTtm                             float32 `csv:"Current ROA (TTM)" json:"current_roa_ttm" parquet:"name=current_roa_ttm, type=FLOAT" db:"current_roa_ttm"`
	Roa5YrAvg                                 float32 `csv:"ROA (5 Yr Avg)" json:"roa_5_yr_avg" parquet:"name=roa_5_yr_avg, type=FLOAT" db:"roa_5_yr_avg"`
	MarketValueToNumberAnalysts               float32 `csv:"Market Value/# Analysts" json:"market_value_to_number_analysts" parquet:"name=market_value_to_number_analysts, type=FLOAT" db:"market_value_to_number_analysts"`
	AnnualSalesMil                            float32 `csv:"Annual Sales ($mil)" json:"annual_sales_mil" parquet:"name=annual_sales_mil, type=FLOAT" db:"annual_sales_mil"`
	CostOfGoodsSoldMil                        float32 `csv:"Cost of Goods Sold ($mil)" json:"cost_of_goods_sold_mil" parquet:"name=cost_of_goods_sold_mil, type=FLOAT" db:"cost_of_goods_sold_mil"`
	EbitdaMil                                 float32 `csv:"EBITDA ($mil)" json:"ebitda_mil" parquet:"name=ebitda_mil, type=FLOAT" db:"ebitda_mil"`
	EbitMil                                   float32 `csv:"EBIT ($mil)" json:"ebit_mil" parquet:"name=ebit_mil, type=FLOAT" db:"ebit_mil"`
	PretaxIncomeMil                           float32 `csv:"Pretax Income ($mil)" json:"pretax_income_mil" parquet:"name=pretax_income_mil, type=FLOAT" db:"pretax_income_mil"`
	NetIncomeMil                              float32 `csv:"Net Income  ($mil)" json:"net_income_mil" parquet:"name=net_income_mil, type=FLOAT" db:"net_income_mil"`
	CashFlowMil                               float32 `csv:"Cash Flow ($mil)" json:"cash_flow_mil" parquet:"name=cash_flow_mil, type=FLOAT" db:"cash_flow_mil"`
	NetIncomeGrowthF0FNeg1                    float32 `csv:"Net Income Growth F(0)/F(-1)" json:"net_income_growth_f0_f_neg1" parquet:"name=net_income_growth_f0_f_neg1, type=FLOAT" db:"net_income_growth_f0_f_neg1"`
	TwelveMoNetIncomeCurrentToLastPercent     float32 `csv:"12 Mo. Net Income Current/Last %" json:"twelve_mo_net_income_current_to_last_percent" parquet:"name=twelve_mo_net_income_current_to_last_percent, type=FLOAT" db:"twelve_mo_net_income_current_to_last_percent"`
	TwelveMoNetIncomeCurrent1qToLast1qPercent float32 `csv:"12 Mo. Net Income Current-1Q/Last-1Q %" json:"twelve_mo_net_income_current_1q_to_last_1q_percent" parquet:"name=twelve_mo_net_income_current_1q_to_last_1q_percent, type=FLOAT" db:"twelve_mo_net_income_current_1q_to_last_1q_percent"`
	DivYieldPercent                           float32 `csv:"Div. Yield %" json:"div_yield_percent" parquet:"name=div_yield_percent, type=FLOAT" db:"div_yield_percent"`
	FiveYrDivYieldPercent                     float32 `csv:"5 Yr Div. Yield %" json:"five_yr_div_yield_percent" parquet:"name=five_yr_div_yield_percent, type=FLOAT" db:"five_yr_div_yield_percent"`
	FiveYrHistDivGrowthPercent                float32 `csv:"5 Yr Hist. Div. Growth %" json:"five_yr_hist_div_growth_percent" parquet:"name=five_yr_hist_div_growth_percent, type=FLOAT" db:"five_yr_hist_div_growth_percent"`
	Dividend                                  float32 `csv:"Dividend " json:"dividend" parquet:"name=dividend, type=FLOAT" db:"dividend"`
	NetMarginPercent                          float32 `csv:"Net Margin %" json:"net_margin_percent" parquet:"name=net_margin_percent, type=FLOAT" db:"net_margin_percent"`
	Turnover                                  float32 `csv:"Turnover" json:"turnover" parquet:"name=turnover, type=FLOAT" db:"turnover"`
	OperatingMargin12MoPercent                float32 `csv:"Operating Margin 12 Mo %" json:"operating_margin_12_mo_percent" parquet:"name=operating_margin_12_mo_percent, type=FLOAT" db:"operating_margin_12_mo_percent"`
	InventoryTurnover                         float32 `csv:"Inventory Turnover" json:"inventory_turnover" parquet:"name=inventory_turnover, type=FLOAT" db:"inventory_turnover"`
	AssetUtilization                          float32 `csv:"Asset Utilization" json:"asset_utilization" parquet:"name=asset_utilization, type=FLOAT" db:"asset_utilization"`
	ReceivablesMil                            float32 `csv:"Receivables ($mil)" json:"receivables_mil" parquet:"name=receivables_mil, type=FLOAT" db:"receivables_mil"`
	IntangiblesMil                            float32 `csv:"Intangibles ($mil)" json:"intangibles_mil" parquet:"name=intangibles_mil, type=FLOAT" db:"intangibles_mil"`
	InventoryMil                              float32 `csv:"Inventory ($mil)" json:"inventory_mil" parquet:"name=inventory_mil, type=FLOAT" db:"inventory_mil"`
	CurrentAssetsMil                          float32 `csv:"Current Assets  ($mil)" json:"current_assets_mil" parquet:"name=current_assets_mil, type=FLOAT" db:"current_assets_mil"`
	CurrentLiabilitiesMil                     float32 `csv:"Current Liabilities ($mil)" json:"current_liabilities_mil" parquet:"name=current_liabilities_mil, type=FLOAT" db:"current_liabilities_mil"`
	LongTermDebtMil                           float32 `csv:"Long Term Debt ($mil)" json:"long_term_debt_mil" parquet:"name=long_term_debt_mil, type=FLOAT" db:"long_term_debt_mil"`
	PreferredEquityMil                        float32 `csv:"Preferred Equity ($mil)" json:"preferred_equity_mil" parquet:"name=preferred_equity_mil, type=FLOAT" db:"preferred_equity_mil"`
	CommonEquityMil                           float32 `csv:"Common Equity ($mil)" json:"common_equity_mil" parquet:"name=common_equity_mil, type=FLOAT" db:"common_equity_mil"`
	BookValue                                 float32 `csv:"Book Value" json:"book_value" parquet:"name=book_value, type=FLOAT" db:"book_value"`
	DebtToTotalCapital                        float32 `csv:"Debt/Total Capital" json:"debt_to_total_capital" parquet:"name=debt_to_total_capital, type=FLOAT" db:"debt_to_total_capital"`
	DebtToEquityRatio                         float32 `csv:"Debt/Equity Ratio" json:"debt_to_equity_ratio" parquet:"name=debt_to_equity_ratio, type=FLOAT" db:"debt_to_equity_ratio"`
	CurrentRatio                              float32 `csv:"Current Ratio" json:"current_ratio" parquet:"name=current_ratio, type=FLOAT" db:"current_ratio"`
	QuickRatio                                float32 `csv:"Quick Ratio" json:"quick_ratio" parquet:"name=quick_ratio, type=FLOAT" db:"quick_ratio"`
	CashRatio                                 float32 `csv:"Cash Ratio" json:"cash_ratio" parquet:"name=cash_ratio, type=FLOAT" db:"cash_ratio"`
}

func LoadRatings(ratingsFn string, limit int) []*ZacksRecord {
	log.Info().Str("RatingsFile", ratingsFn).Msg("loading ratings from file")
	fh, err := os.OpenFile(ratingsFn, os.O_RDONLY, 0644)
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Str("FileName", ratingsFn).Msg("failed to open file")
		return make([]*ZacksRecord, 0)
	}
	defer fh.Close()

	records := []*ZacksRecord{}

	if err := gocsv.UnmarshalFile(fh, &records); err != nil { // Load clients from file
		log.Error().Str("OriginalError", err.Error()).Str("FileName", ratingsFn).Msg("failed to open file")
		return make([]*ZacksRecord, 0)
	}

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	return records
}

func isValidExchange(record *ZacksRecord) bool {
	return (record.Exchange != "OTCQX" &&
		record.Exchange != "OTCQB" &&
		record.Exchange != "OTC Markets" &&
		record.Exchange != "Grey Market" &&
		record.Exchange != "Pink No Info" &&
		record.Exchange != "Pink Current Info")
}

func EnrichWithFigi(records []*ZacksRecord) []*ZacksRecord {
	conn, err := pgx.Connect(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Could not connect to database")
	}
	defer conn.Close(context.Background())

	// build a list of all active records that have SA composite figi's
	saIdMap := make(map[int]*Ticker)
	rows, err := conn.Query(context.Background(), "SELECT ticker, seeking_alpha_id, composite_figi FROM tickers_v1 WHERE active='t' AND seeking_alpha_id IS NOT NULL AND composite_figi IS NOT NULL")
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Failed to retrieve tickers from database")
	}

	for rows.Next() {
		var ticker Ticker
		err := rows.Scan(&ticker.Ticker, &ticker.TickerId, &ticker.CompositeFigi)
		if err != nil {
			log.Error().Str("OriginalError", err.Error()).Msg("Failed to retrieve ticker row from database")
		}
		saIdMap[ticker.TickerId] = &ticker
	}

	return records
}

func SaveToDB(records []*ZacksRecord) {
	conn, err := pgx.Connect(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Could not connect to database")
	}
	defer conn.Close(context.Background())

	for _, r := range records {
		fmt.Println(r)
		/*
			conn.Exec(context.Background(),
				`INSERT INTO seeking_alpha_v1 (
				"ticker",
				"composite_figi",
				"event_date",
				"market_cap_mil",
				"quant_rating",
				"growth_grade",
				"profitability_grade",
				"value_grade",
				"eps_revisions_grade",
				"authors_rating_pro",
				"sell_side_rating"
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
				$11
			) ON CONFLICT ON CONSTRAINT seeking_alpha_v1_pkey
			DO UPDATE SET
				market_cap_mil = EXCLUDED.market_cap_mil,
				quant_rating = EXCLUDED.quant_rating,
				growth_grade = EXCLUDED.growth_grade,
				profitability_grade = EXCLUDED.profitability_grade,
				value_grade = EXCLUDED.value_grade,
				eps_revisions_grade = EXCLUDED.eps_revisions_grade,
				authors_rating_pro = EXCLUDED.authors_rating_pro,
				sell_side_rating = EXCLUDED.sell_side_rating;
			`,
				r.Ticker, r.CompositeFigi, r.Date, r.MarketCap/1e6,
				r.QuantRating, r.GrowthCategory, r.ProfitabilityCategory,
				r.ValueCategory, r.EpsRevisionsCategory,
				r.AuthorsRatingPro, r.SellSideRating)
		*/
	}

	log.Info().Int("NumRecords", len(records)).Msg("records saved to DB")
}

func SaveToParquet(records []*ZacksRecord, fn string) error {
	var err error

	fh, err := local.NewLocalFileWriter(fn)
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Str("FileName", fn).Msg("cannot create local file")
		return err
	}
	defer fh.Close()

	pw, err := writer.NewParquetWriter(fh, new(ZacksRecord), 4)
	if err != nil {
		log.Error().
			Str("OriginalError", err.Error()).
			Msg("Parquet write failed")
		return err
	}

	pw.RowGroupSize = 128 * 1024 * 1024 // 128M
	pw.PageSize = 8 * 1024              // 8k
	pw.CompressionType = parquet.CompressionCodec_GZIP

	for _, r := range records {
		if err = pw.Write(r); err != nil {
			log.Error().
				Str("OriginalError", err.Error()).
				Str("EventDate", r.EventDate).Str("Ticker", r.Ticker).
				Str("CompositeFigi", r.CompositeFigi).
				Msg("Parquet write failed for record")
		}
	}

	if err = pw.WriteStop(); err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Parquet write failed")
		return err
	}

	log.Info().Int("NumRecords", len(records)).Msg("Parquet write finished")
	return nil
}
