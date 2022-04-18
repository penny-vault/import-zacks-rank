package zacksimport

import (
	"context"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type Ticker struct {
	CompanyName   string `db:"name"`
	Ticker        string `db:"ticker"`
	CompositeFigi string `db:"composite_figi"`
}

type ZacksRecord struct {
	CompanyName                               string    `csv:"Company Name" json:"company_name" parquet:"name=company_name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Ticker                                    string    `csv:"Ticker" json:"ticker" parquet:"name=ticker, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"ticker,omitempty"`
	CompositeFigi                             string    `csv:"-" json:"composite_figi" parquet:"name=composite_figi, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"composite_figi,omitempty"`
	Exchange                                  string    `csv:"Exchange" json:"exchange" parquet:"name=exchange, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	EventDateStr                              string    `csv:"-" json:"event_date" parquet:"name=event_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	EventDate                                 time.Time `csv:"-" json:"-"`
	InSp500                                   bool      `csv:"S&P 500 - ETF" json:"in_sp500" parquet:"name=in_sp500, type=BOOLEAN" db:"in_sp500"`
	LastClose                                 float64   `csv:"Last Close" json:"last_close" parquet:"name=last_close, type=DOUBLE"`
	MonthOfFiscalYrEnd                        int       `csv:"Month of Fiscal Yr End" json:"month_of_fiscal_yr_end" parquet:"name=month_of_fiscal_yr_end, type=INT32" db:"month_of_fiscal_yr_end,omitempty"`
	Optionable                                bool      `csv:"Optionable" json:"optionable" parquet:"name=optionable, type=BOOLEAN" db:"optionable"`
	Sector                                    string    `csv:"Sector" json:"sector" parquet:"name=sector, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"sector,omitempty"`
	Industry                                  string    `csv:"Industry" json:"industry" parquet:"name=industry, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"industry,omitempty"`
	SharesOutstandingMil                      float64   `csv:"Shares Outstanding (mil)" json:"shares_outstanding_mil" parquet:"name=shares_outstanding_mil, type=DOUBLE" db:"shares_outstanding_mil,omitempty"`
	MarketCapMil                              float64   `csv:"Market Cap (mil)" json:"market_cap_mil" parquet:"name=market_cap_mil, type=DOUBLE" db:"market_cap_mil,omitempty"`
	AvgVolume                                 int64     `csv:"Avg Volume" json:"avg_volume" parquet:"name=avg_volume, type=INT64" db:"avg_volume,omitempty"`
	WkHigh52                                  float64   `csv:"52 Week High" json:"wk_high_52" parquet:"name=wk_high_52, type=DOUBLE" db:"wk_high_52,omitempty"`
	WkLow52                                   float64   `csv:"52 Week Low" json:"wk_low_52" parquet:"name=wk_low_52, type=DOUBLE" db:"wk_low_52,omitempty"`
	PriceAsPercentOf52wkHighLow               float32   `csv:"Price as a % of 52 Wk H-L Range" json:"price_as_percent_of_52wk_hl" parquet:"name=price_as_percent_of_52wk_hl, type=FLOAT" db:"price_as_percent_of_52wk_hl,omitempty"`
	Beta                                      float32   `csv:"Beta" json:"beta" parquet:"name=beta, type=FLOAT" db:"beta,omitempty"`
	PercentPriceChange1Wk                     float32   `csv:"% Price Change (1 Week)" json:"percent_price_change_1wk" parquet:"name=percent_price_change_1wk, type=FLOAT" db:"percent_price_change_1wk,omitempty"`
	PercentPriceChange4Wk                     float32   `csv:"% Price Change (4 Weeks)" json:"percent_price_change_4wk" parquet:"name=percent_price_change_4wk, type=FLOAT" db:"percent_price_change_4wk,omitempty"`
	PercentPriceChange12Wk                    float32   `csv:"% Price Change (12 Weeks)" json:"percent_price_change_12wk" parquet:"name=percent_price_change_12wk, type=FLOAT" db:"percent_price_change_12wk,omitempty"`
	PercentPriceChangeYtd                     float32   `csv:"% Price Change (YTD)" json:"percent_price_change_ytd" parquet:"name=percent_price_change_ytd, type=FLOAT" db:"percent_price_change_ytd,omitempty"`
	RelativePriceChange                       float32   `csv:"Relative Price Change" json:"relative_price_change" parquet:"name=relative_price_change, type=FLOAT" db:"relative_price_change,omitempty"`
	ZacksRank                                 int       `csv:"Zacks Rank" json:"zacks_rank" parquet:"name=zacks_rank, type=INT32" db:"zacks_rank,omitempty"`
	ZacksRankChangeIndicator                  int       `csv:"Zacks Rank Change Indicator" json:"zacks_rank_change_indicator" parquet:"name=zacks_rank_change_indicator, type=INT32" db:"zacks_rank_change_indicator,omitempty"`
	ZacksIndustryRank                         int       `csv:"Zacks Industry Rank" json:"zacks_industry_rank" parquet:"name=zacks_industry_rank, type=INT32" db:"zacks_industry_rank,omitempty"`
	ValueScore                                string    `csv:"Value Score" json:"value_score" parquet:"name=value_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"value_score,omitempty"`
	GrowthScore                               string    `csv:"Growth Score" json:"growth_score" parquet:"name=growth_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"growth_score,omitempty"`
	MomentumScore                             string    `csv:"Momentum Score" json:"momentum_score" parquet:"name=momentum_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"momentum_score,omitempty"`
	VgmScore                                  string    `csv:"VGM Score" json:"vgm_score" parquet:"name=vgm_score, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" db:"vgm_score,omitempty"`
	CurrentAvgBrokerRec                       float32   `csv:"Current Avg Broker Rec" json:"current_avg_broker_rec" parquet:"name=current_avg_broker_rec, type=FLOAT" db:"current_avg_broker_rec,omitempty"`
	NumBrokersInRating                        int       `csv:"# of Brokers in Rating" json:"num_brokers_in_rating" parquet:"name=num_brokers_in_rating, type=INT32" db:"num_brokers_in_rating,omitempty"`
	NumRatingStrongBuyOrBuy                   int       `csv:"# Rating Strong Buy or Buy" json:"num_rating_strong_buy_or_buy" parquet:"name=num_rating_strong_buy_or_buy, type=INT32" db:"num_rating_strong_buy_or_buy,omitempty"`
	PercentRatingStrongBuyOrBuy               float32   `csv:"% Rating Strong Buy or Buy" json:"percent_rating_strong_buy_or_buy" parquet:"name=percent_rating_strong_buy_or_buy, type=FLOAT" db:"percent_rating_strong_buy_or_buy,omitempty"`
	NumRatingHold                             int       `csv:"# Rating Hold" json:"num_rating_hold" parquet:"name=num_rating_hold, type=INT32" db:"num_rating_hold,omitempty"`
	NumRatingStrongSellOrSell                 int       `csv:"# Rating Strong Sell or Sell" json:"num_rating_strong_sell_or_sell" parquet:"name=num_rating_strong_sell_or_sell, type=INT32" db:"num_rating_strong_sell_or_sell,omitempty"`
	PercentRatingStrongSellOrSell             float32   `csv:"% Rating Strong Sell or Sell" json:"percent_rating_strong_sell_or_sell" parquet:"name=percent_rating_strong_sell_or_sell, type=FLOAT" db:"percent_rating_strong_sell_or_sell,omitempty"`
	PercentRatingChange4Wk                    float32   `csv:"% Rating Change - 4 Weeks" json:"percent_rating_change_4wk" parquet:"name=percent_rating_change_4wk, type=FLOAT" db:"percent_rating_change_4wk,omitempty"`
	IndustryRankOfAbr                         int       `csv:"Industry Rank (of ABR)" json:"industry_rank_of_abr" parquet:"name=industry_rank_of_abr, type=INT32" db:"industry_rank_of_abr,omitempty"`
	RankInIndustryOfAbr                       int       `csv:"Rank in Industry (of ABR)" json:"rank_in_industry_of_abr" parquet:"name=rank_in_industry_of_abr, type=INT32" db:"rank_in_industry_of_abr,omitempty"`
	ChangeInAvgRec                            float32   `csv:"Change in Avg Rec " json:"change_in_avg_rec" parquet:"name=change_in_avg_rec, type=FLOAT" db:"change_in_avg_rec,omitempty"`
	NumberRatingUpgrades                      int       `csv:"# Rating Upgrades" json:"number_rating_upgrades" parquet:"name=number_rating_upgrades, type=INT32" db:"number_rating_upgrades,omitempty"`
	NumberRatingDowngrades                    int       `csv:"# Rating Downgrades " json:"number_rating_downgrades" parquet:"name=number_rating_downgrades, type=INT32" db:"number_rating_downgrades,omitempty"`
	PercentRatingHold                         float32   `csv:"% Rating Hold" json:"percent_rating_hold" parquet:"name=percent_rating_hold, type=FLOAT" db:"percent_rating_hold,omitempty"`
	PercentRatingUpgrades                     float32   `csv:"% Rating Upgrades " json:"percent_rating_upgrades" parquet:"name=percent_rating_upgrades, type=FLOAT" db:"percent_rating_upgrades,omitempty"`
	PercentRatingDowngrades                   float32   `csv:"% Rating Downgrades " json:"percent_rating_downgrades" parquet:"name=percent_rating_downgrades, type=FLOAT" db:"percent_rating_downgrades,omitempty"`
	AverageTargetPrice                        float64   `csv:"Average Target Price" json:"average_target_price" parquet:"name=average_target_price, type=DOUBLE" db:"average_target_price,omitempty"`
	EarningsEsp                               float32   `csv:"Earnings ESP" json:"earnings_esp" parquet:"name=earnings_esp, type=FLOAT" db:"earnings_esp,omitempty"`
	LastEpsSurprisePercent                    float32   `csv:"Last EPS Surprise (%)" json:"last_eps_surprise_percent" parquet:"name=last_eps_surprise_percent, type=FLOAT" db:"last_eps_surprise_percent,omitempty"`
	PreviousEpsSurprisePercent                float32   `csv:"Previous EPS Surprise (%)" json:"previous_eps_surprise_percent" parquet:"name=previous_eps_surprise_percent, type=FLOAT" db:"previous_eps_surprise_percent,omitempty"`
	AvgEpsSurpriseLast4Qtrs                   float32   `csv:"Avg EPS Surprise (Last 4 Qtrs)" json:"avg_eps_surprise_last_4_qtrs" parquet:"name=avg_eps_surprise_last_4_qtrs, type=FLOAT" db:"avg_eps_surprise_last_4_qtrs,omitempty"`
	ActualEpsUsedInSurpriseDollarsPerShare    float32   `csv:"Actual EPS used in Surprise ($/sh)" json:"actual_eps_used_in_surprise_dollars_per_share" parquet:"name=actual_eps_used_in_surprise_dollars_per_share, type=FLOAT" db:"actual_eps_used_in_surprise_dollars_per_share,omitempty"`
	LastQtrEps                                float32   `csv:"Last Qtr EPS" json:"last_qtr_eps" parquet:"name=last_qtr_eps, type=FLOAT" db:"last_qtr_eps,omitempty"`
	LastReportedQtrDateStr                    string    `csv:"Last Reported Qtr (yyyymm)" json:"last_reported_qtr_date" parquet:"name=last_reported_qtr_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	LastReportedQtrDate                       time.Time `csv:"-" json:"-" db:"last_reported_qtr_date,omitempty"`
	LastYrEpsF0BeforeNri                      float32   `csv:"Last Yr's EPS (F0) Before NRI" json:"last_yr_eps_f0_before_nri" parquet:"name=last_yr_eps_f0_before_nri, type=FLOAT" db:"last_yr_eps_f0_before_nri,omitempty"`
	TwelveMoTrailingEps                       float32   `csv:"12 Mo Trailing EPS" json:"twelve_mo_trailing_eps" parquet:"name=twelve_mo_trailing_eps, type=FLOAT" db:"twelve_mo_trailing_eps,omitempty"`
	LastReportedFiscalYrStr                   string    `csv:"Last Reported Fiscal Yr  (yyyymm)" json:"last_reported_fiscal_yr" parquet:"name=last_reported_fiscal_yr, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	LastReportedFiscalYr                      time.Time `csv:"-" json:"-" db:"last_reported_fiscal_yr,omitempty"`
	LastEpsReportDateStr                      string    `csv:"Last EPS Report Date (yyyymmdd)" json:"last_eps_report_date" parquet:"name=last_eps_report_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	LastEpsReportDate                         time.Time `csv:"-" json:"-" db:"last_eps_report_date,omitempty"`
	NextEpsReportDateStr                      string    `csv:"Next EPS Report Date  (yyyymmdd)" json:"next_eps_report_date" parquet:"name=next_eps_report_date, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	NextEpsReportDate                         time.Time `csv:"-" json:"-" db:"next_eps_report_date,omitempty"`
	PercentChangeQ0Est                        float32   `csv:"% Change Q0 Est. (4 weeks)" json:"percent_change_q0_est" parquet:"name=percent_change_q0_est, type=FLOAT" db:"percent_change_q0_est,omitempty"`
	PercentChangeQ2Est                        float32   `csv:"% Change Q2 Est. (4 weeks)" json:"percent_change_q2_est" parquet:"name=percent_change_q2_est, type=FLOAT" db:"percent_change_q2_est,omitempty"`
	PercentChangeF1Est                        float32   `csv:"% Change F1 Est. (4 weeks)" json:"percent_change_f1_est" parquet:"name=percent_change_f1_est, type=FLOAT" db:"percent_change_f1_est,omitempty"`
	PercentChangeQ1Est                        float32   `csv:"% Change Q1 Est. (4 weeks)" json:"percent_change_q1_est" parquet:"name=percent_change_q1_est, type=FLOAT" db:"percent_change_q1_est,omitempty"`
	PercentChangeF2Est                        float32   `csv:"% Change F2 Est. (4 weeks)" json:"percent_change_f2_est" parquet:"name=percent_change_f2_est, type=FLOAT" db:"percent_change_f2_est,omitempty"`
	PercentChangeLtGrowthEst                  float32   `csv:"% Change LT Growth Est. (4 weeks)" json:"percent_change_lt_growth_est" parquet:"name=percent_change_lt_growth_est, type=FLOAT" db:"percent_change_lt_growth_est,omitempty"`
	Q0ConsensusEstLastCompletedFiscalQtr      float32   `csv:"Q0 Consensus Est. (last completed fiscal Qtr)" json:"q0_consensus_est_last_completed_fiscal_qtr" parquet:"name=q0_consensus_est_last_completed_fiscal_qtr, type=FLOAT" db:"q0_consensus_est_last_completed_fiscal_qtr,omitempty"`
	NumberOfAnalystsInQ0Consensus             int       `csv:"# of Analysts in Q0 Consensus" json:"number_of_analysts_in_q0_consensus" parquet:"name=number_of_analysts_in_q0_consensus, type=INT32" db:"number_of_analysts_in_q0_consensus,omitempty"`
	Q1ConsensusEst                            float32   `csv:"Q1 Consensus Est. " json:"q1_consensus_est" parquet:"name=q1_consensus_est, type=FLOAT" db:"q1_consensus_est,omitempty"`
	NumberOfAnalystsInQ1Consensus             int       `csv:"# of Analysts in Q1 Consensus" json:"number_of_analysts_in_q1_consensus" parquet:"name=number_of_analysts_in_q1_consensus, type=INT32" db:"number_of_analysts_in_q1_consensus,omitempty"`
	StdevQ1Q1ConsensusRatio                   float32   `csv:"St. Dev. Q1 / Q1 Consensus" json:"stdev_q1_q1_consensus_ratio" parquet:"name=stdev_q1_q1_consensus_ratio, type=FLOAT" db:"stdev_q1_q1_consensus_ratio,omitempty"`
	Q2ConsensusEstNextFiscalQtr               float32   `csv:"Q2 Consensus Est. (next fiscal Qtr)" json:"q2_consensus_est_next_fiscal_qtr" parquet:"name=q2_consensus_est_next_fiscal_qtr, type=FLOAT" db:"q2_consensus_est_next_fiscal_qtr,omitempty"`
	NumberOfAnalystsInQ2Consensus             int       `csv:"# of Analysts in Q2 Consensus" json:"number_of_analysts_in_q2_consensus" parquet:"name=number_of_analysts_in_q2_consensus, type=INT32" db:"number_of_analysts_in_q2_consensus,omitempty"`
	StdevQ2Q2ConsensusRatio                   float32   `csv:"St. Dev. Q2 / Q2 Consensus" json:"stdev_q2_q2_consensus_ratio" parquet:"name=stdev_q2_q2_consensus_ratio, type=FLOAT" db:"stdev_q2_q2_consensus_ratio,omitempty"`
	F0ConsensusEst                            float32   `csv:"F0 Consensus Est." json:"f0_consensus_est" parquet:"name=f0_consensus_est, type=FLOAT" db:"f0_consensus_est,omitempty"`
	NumberOfAnalystsInF0Consensus             float32   `csv:"# of Analysts in F0 Consensus" json:"number_of_analysts_in_f0_consensus" parquet:"name=number_of_analysts_in_f0_consensus, type=FLOAT" db:"number_of_analysts_in_f0_consensus,omitempty"`
	F1ConsensusEst                            float32   `csv:"F1 Consensus Est." json:"f1_consensus_est" parquet:"name=f1_consensus_est, type=FLOAT" db:"f1_consensus_est,omitempty"`
	NumberOfAnalystsInF1Consensus             int       `csv:"# of Analysts in F1 Consensus" json:"number_of_analysts_in_f1_consensus" parquet:"name=number_of_analysts_in_f1_consensus, type=INT32" db:"number_of_analysts_in_f1_consensus,omitempty"`
	StdevF1F1ConsensusRatio                   float32   `csv:"St. Dev. F1 / F1 Consensus" json:"stdev_f1_f1_consensus_ratio" parquet:"name=stdev_f1_f1_consensus_ratio, type=FLOAT" db:"stdev_f1_f1_consensus_ratio,omitempty"`
	F2ConsensusEst                            float32   `csv:"F2 Consensus Est." json:"f2_consensus_est" parquet:"name=f2_consensus_est, type=FLOAT" db:"f2_consensus_est,omitempty"`
	NumberOfAnalystsInF2Consensus             int       `csv:"# of Analysts in F2 Consensus" json:"number_of_analysts_in_f2_consensus" parquet:"name=number_of_analysts_in_f2_consensus, type=INT32" db:"number_of_analysts_in_f2_consensus,omitempty"`
	FiveYrHistEpsGrowth                       float32   `csv:"5 Yr. Hist. EPS Growth" json:"five_yr_hist_eps_growth" parquet:"name=five_yr_hist_eps_growth, type=FLOAT" db:"five_yr_hist_eps_growth,omitempty"`
	LongTermGrowthConsensusEst                float32   `csv:"Long-Term Growth Consensus Est." json:"long_term_growth_consensus_est" parquet:"name=long_term_growth_consensus_est, type=FLOAT" db:"long_term_growth_consensus_est,omitempty"`
	PercentChangeEps                          float32   `csv:"% Change EPS (F(-1)/F(-2))" json:"percent_change_eps" parquet:"name=percent_change_eps, type=FLOAT" db:"percent_change_eps,omitempty"`
	LastYrsGrowth                             float32   `csv:"Last Yrs Growth (F[0] / F [-1])" json:"last_yrs_growth" parquet:"name=last_yrs_growth, type=FLOAT" db:"last_yrs_growth,omitempty"`
	ThisYrsEstGrowth                          float32   `csv:"This Yr's Est.d Growth (F(1)/F(0))" json:"this_yrs_est_growth" parquet:"name=this_yrs_est_growth, type=FLOAT" db:"this_yrs_est_growth,omitempty"`
	PercentRatioOfQ1Q0                        float32   `csv:"% Ratio of Q1/Q0" json:"percent_ratio_of_q1_q0" parquet:"name=percent_ratio_of_q1_q0, type=FLOAT" db:"percent_ratio_of_q1_q0,omitempty"`
	PercentRatioOfQ1PriorYrQ1ActualQ          float32   `csv:"% Ratio of Q1/prior Yr Q1 Actual Q(-3)" json:"percent_ratio_of_q1_prior_yr_q1_actual_q" parquet:"name=percent_ratio_of_q1_prior_yr_q1_actual_q, type=FLOAT" db:"percent_ratio_of_q1_prior_yr_q1_actual_q,omitempty"`
	SalesGrowth                               float32   `csv:"Sales Growth F(0)/F(-1)" json:"sales_growth" parquet:"name=sales_growth, type=FLOAT" db:"sales_growth,omitempty"`
	FiveYrHistoricalSalesGrowth               float32   `csv:"5 Yr Historical Sales Growth" json:"five_yr_historical_sales_growth" parquet:"name=five_yr_historical_sales_growth, type=FLOAT" db:"five_yr_historical_sales_growth,omitempty"`
	Q1ConsensusSalesEstMil                    float32   `csv:"Q(1) Consensus Sales Est. ($mil)" json:"q1_consensus_sales_est_mil" parquet:"name=q1_consensus_sales_est_mil, type=FLOAT" db:"q1_consensus_sales_est_mil,omitempty"`
	F1ConsensusSalesEstMil                    float32   `csv:"F(1) Consensus Sales Est. ($mil)" json:"f1_consensus_sales_est_mil" parquet:"name=f1_consensus_sales_est_mil, type=FLOAT" db:"f1_consensus_sales_est_mil,omitempty"`
	PeTrailing12Months                        float32   `csv:"P/E (Trailing 12 Months)" json:"pe_trailing_12_months" parquet:"name=pe_trailing_12_months, type=FLOAT" db:"pe_trailing_12_months,omitempty"`
	PeF1                                      float32   `csv:"P/E (F1)" json:"pe_f1" parquet:"name=pe_f1, type=FLOAT" db:"pe_f1,omitempty"`
	PeF2                                      float32   `csv:"P/E (F2)" json:"pe_f2" parquet:"name=pe_f2, type=FLOAT" db:"pe_f2,omitempty"`
	PegRatio                                  float32   `csv:"PEG Ratio" json:"peg_ratio" parquet:"name=peg_ratio, type=FLOAT" db:"peg_ratio,omitempty"`
	PriceToCashFlow                           float32   `csv:"Price/Cash Flow" json:"price_to_cash_flow" parquet:"name=price_to_cash_flow, type=FLOAT" db:"price_to_cash_flow,omitempty"`
	PriceToSales                              float32   `csv:"Price/Sales" json:"price_to_sales" parquet:"name=price_to_sales, type=FLOAT" db:"price_to_sales,omitempty"`
	PriceToBook                               float32   `csv:"Price/Book" json:"price_to_book" parquet:"name=price_to_book, type=FLOAT" db:"price_to_book,omitempty"`
	CurrentRoeTtm                             float32   `csv:"Current ROE (TTM)" json:"current_roe_ttm" parquet:"name=current_roe_ttm, type=FLOAT" db:"current_roe_ttm,omitempty"`
	CurrentRoiTtm                             float32   `csv:"Current ROI (TTM)" json:"current_roi_ttm" parquet:"name=current_roi_ttm, type=FLOAT" db:"current_roi_ttm,omitempty"`
	Roi5YrAvg                                 float32   `csv:"ROI (5 Yr Avg)" json:"roi_5_yr_avg" parquet:"name=roi_5_yr_avg, type=FLOAT" db:"roi_5_yr_avg,omitempty"`
	CurrentRoaTtm                             float32   `csv:"Current ROA (TTM)" json:"current_roa_ttm" parquet:"name=current_roa_ttm, type=FLOAT" db:"current_roa_ttm,omitempty"`
	Roa5YrAvg                                 float32   `csv:"ROA (5 Yr Avg)" json:"roa_5_yr_avg" parquet:"name=roa_5_yr_avg, type=FLOAT" db:"roa_5_yr_avg,omitempty"`
	MarketValueToNumberAnalysts               float32   `csv:"Market Value/# Analysts" json:"market_value_to_number_analysts" parquet:"name=market_value_to_number_analysts, type=FLOAT" db:"market_value_to_number_analysts,omitempty"`
	AnnualSalesMil                            float32   `csv:"Annual Sales ($mil)" json:"annual_sales_mil" parquet:"name=annual_sales_mil, type=FLOAT" db:"annual_sales_mil,omitempty"`
	CostOfGoodsSoldMil                        float32   `csv:"Cost of Goods Sold ($mil)" json:"cost_of_goods_sold_mil" parquet:"name=cost_of_goods_sold_mil, type=FLOAT" db:"cost_of_goods_sold_mil,omitempty"`
	EbitdaMil                                 float32   `csv:"EBITDA ($mil)" json:"ebitda_mil" parquet:"name=ebitda_mil, type=FLOAT" db:"ebitda_mil,omitempty"`
	EbitMil                                   float32   `csv:"EBIT ($mil)" json:"ebit_mil" parquet:"name=ebit_mil, type=FLOAT" db:"ebit_mil,omitempty"`
	PretaxIncomeMil                           float32   `csv:"Pretax Income ($mil)" json:"pretax_income_mil" parquet:"name=pretax_income_mil, type=FLOAT" db:"pretax_income_mil,omitempty"`
	NetIncomeMil                              float32   `csv:"Net Income  ($mil)" json:"net_income_mil" parquet:"name=net_income_mil, type=FLOAT" db:"net_income_mil,omitempty"`
	CashFlowMil                               float32   `csv:"Cash Flow ($mil)" json:"cash_flow_mil" parquet:"name=cash_flow_mil, type=FLOAT" db:"cash_flow_mil,omitempty"`
	NetIncomeGrowthF0FNeg1                    float32   `csv:"Net Income Growth F(0)/F(-1)" json:"net_income_growth_f0_f_neg1" parquet:"name=net_income_growth_f0_f_neg1, type=FLOAT" db:"net_income_growth_f0_f_neg1,omitempty"`
	TwelveMoNetIncomeCurrentToLastPercent     float32   `csv:"12 Mo. Net Income Current/Last %" json:"twelve_mo_net_income_current_to_last_percent" parquet:"name=twelve_mo_net_income_current_to_last_percent, type=FLOAT" db:"twelve_mo_net_income_current_to_last_percent,omitempty"`
	TwelveMoNetIncomeCurrent1qToLast1qPercent float32   `csv:"12 Mo. Net Income Current-1Q/Last-1Q %" json:"twelve_mo_net_income_current_1q_to_last_1q_percent" parquet:"name=twelve_mo_net_income_current_1q_to_last_1q_percent, type=FLOAT" db:"twelve_mo_net_income_current_1q_to_last_1q_percent,omitempty"`
	DivYieldPercent                           float32   `csv:"Div. Yield %" json:"div_yield_percent" parquet:"name=div_yield_percent, type=FLOAT" db:"div_yield_percent,omitempty"`
	FiveYrDivYieldPercent                     float32   `csv:"5 Yr Div. Yield %" json:"five_yr_div_yield_percent" parquet:"name=five_yr_div_yield_percent, type=FLOAT" db:"five_yr_div_yield_percent,omitempty"`
	FiveYrHistDivGrowthPercent                float32   `csv:"5 Yr Hist. Div. Growth %" json:"five_yr_hist_div_growth_percent" parquet:"name=five_yr_hist_div_growth_percent, type=FLOAT" db:"five_yr_hist_div_growth_percent,omitempty"`
	Dividend                                  float32   `csv:"Dividend " json:"dividend" parquet:"name=dividend, type=FLOAT" db:"dividend,omitempty"`
	NetMarginPercent                          float32   `csv:"Net Margin %" json:"net_margin_percent" parquet:"name=net_margin_percent, type=FLOAT" db:"net_margin_percent,omitempty"`
	Turnover                                  float32   `csv:"Turnover" json:"turnover" parquet:"name=turnover, type=FLOAT" db:"turnover,omitempty"`
	OperatingMargin12MoPercent                float32   `csv:"Operating Margin 12 Mo %" json:"operating_margin_12_mo_percent" parquet:"name=operating_margin_12_mo_percent, type=FLOAT" db:"operating_margin_12_mo_percent,omitempty"`
	InventoryTurnover                         float32   `csv:"Inventory Turnover" json:"inventory_turnover" parquet:"name=inventory_turnover, type=FLOAT" db:"inventory_turnover,omitempty"`
	AssetUtilization                          float32   `csv:"Asset Utilization" json:"asset_utilization" parquet:"name=asset_utilization, type=FLOAT" db:"asset_utilization,omitempty"`
	ReceivablesMil                            float32   `csv:"Receivables ($mil)" json:"receivables_mil" parquet:"name=receivables_mil, type=FLOAT" db:"receivables_mil,omitempty"`
	IntangiblesMil                            float32   `csv:"Intangibles ($mil)" json:"intangibles_mil" parquet:"name=intangibles_mil, type=FLOAT" db:"intangibles_mil,omitempty"`
	InventoryMil                              float32   `csv:"Inventory ($mil)" json:"inventory_mil" parquet:"name=inventory_mil, type=FLOAT" db:"inventory_mil,omitempty"`
	CurrentAssetsMil                          float32   `csv:"Current Assets  ($mil)" json:"current_assets_mil" parquet:"name=current_assets_mil, type=FLOAT" db:"current_assets_mil,omitempty"`
	CurrentLiabilitiesMil                     float32   `csv:"Current Liabilities ($mil)" json:"current_liabilities_mil" parquet:"name=current_liabilities_mil, type=FLOAT" db:"current_liabilities_mil,omitempty"`
	LongTermDebtMil                           float32   `csv:"Long Term Debt ($mil)" json:"long_term_debt_mil" parquet:"name=long_term_debt_mil, type=FLOAT" db:"long_term_debt_mil,omitempty"`
	PreferredEquityMil                        float32   `csv:"Preferred Equity ($mil)" json:"preferred_equity_mil" parquet:"name=preferred_equity_mil, type=FLOAT" db:"preferred_equity_mil,omitempty"`
	CommonEquityMil                           float32   `csv:"Common Equity ($mil)" json:"common_equity_mil" parquet:"name=common_equity_mil, type=FLOAT" db:"common_equity_mil,omitempty"`
	BookValue                                 float32   `csv:"Book Value" json:"book_value" parquet:"name=book_value, type=FLOAT" db:"book_value,omitempty"`
	DebtToTotalCapital                        float32   `csv:"Debt/Total Capital" json:"debt_to_total_capital" parquet:"name=debt_to_total_capital, type=FLOAT" db:"debt_to_total_capital,omitempty"`
	DebtToEquityRatio                         float32   `csv:"Debt/Equity Ratio" json:"debt_to_equity_ratio" parquet:"name=debt_to_equity_ratio, type=FLOAT" db:"debt_to_equity_ratio,omitempty"`
	CurrentRatio                              float32   `csv:"Current Ratio" json:"current_ratio" parquet:"name=current_ratio, type=FLOAT" db:"current_ratio,omitempty"`
	QuickRatio                                float32   `csv:"Quick Ratio" json:"quick_ratio" parquet:"name=quick_ratio, type=FLOAT" db:"quick_ratio,omitempty"`
	CashRatio                                 float32   `csv:"Cash Ratio" json:"cash_ratio" parquet:"name=cash_ratio, type=FLOAT" db:"cash_ratio,omitempty"`
}

func LoadRatings(ratingsFn string, dateStr string, limit int) []*ZacksRecord {
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

	// cleanup records
	for _, r := range records {
		// set event date
		r.EventDateStr = dateStr
		r.EventDate, _ = time.Parse("2006-01-02", dateStr)

		// parse date fields
		dt, err := time.Parse("200601", r.LastReportedFiscalYrStr)
		if err == nil {
			r.LastReportedFiscalYrStr = dt.Format("2006-01-02")
			r.LastReportedFiscalYr = dt
		} else {
			if r.LastReportedFiscalYrStr != "" {
				log.Warn().Str("Ticker", r.Ticker).Str("InputString", r.LastReportedFiscalYrStr).Msg("could not parse last reported fiscal year")
			}
		}

		dt, err = time.Parse("200601", r.LastReportedQtrDateStr)
		if err == nil {
			r.LastReportedQtrDateStr = dt.Format("2006-01-02")
			r.LastReportedQtrDate = dt
		} else {
			if r.LastReportedQtrDateStr != "" {
				log.Warn().Str("Ticker", r.Ticker).Str("InputString", r.LastReportedQtrDateStr).Msg("could not parse last reported quarter date")
			}
		}

		dt, err = time.Parse("20060102", r.LastEpsReportDateStr)
		if err == nil {
			r.LastEpsReportDateStr = dt.Format("2006-01-02")
			r.LastEpsReportDate = dt
		} else {
			if r.LastEpsReportDateStr != "" {
				log.Warn().Str("Ticker", r.Ticker).Str("InputString", r.LastEpsReportDateStr).Msg("could not parse last eps report date")
			}
		}

		dt, err = time.Parse("20060102", r.NextEpsReportDateStr)
		if err == nil {
			r.NextEpsReportDateStr = dt.Format("2006-01-02")
			r.NextEpsReportDate = dt
		} else {
			if r.NextEpsReportDateStr != "" {
				log.Warn().Str("Ticker", r.Ticker).Str("InputString", r.NextEpsReportDateStr).Msg("could not parse next eps report date")
			}
		}

	}

	return records
}

func isValidExchange(record *ZacksRecord) bool {
	return (record.Exchange != "OTC" &&
		record.Exchange != "OTCBB")
}

func EnrichWithFigi(records []*ZacksRecord) []*ZacksRecord {
	conn, err := pgx.Connect(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Could not connect to database")
	}
	defer conn.Close(context.Background())

	// build a list of all active records that have composite figi's
	tickerMap := make(map[string]*Ticker)

	rows, err := conn.Query(context.Background(), "SELECT ticker, name, composite_figi FROM tickers_v1 WHERE active='t' AND composite_figi IS NOT NULL")
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Failed to retrieve tickers from database")
	}

	for rows.Next() {
		var ticker Ticker
		err := rows.Scan(&ticker.Ticker, &ticker.CompanyName, &ticker.CompositeFigi)
		if err != nil {
			log.Error().Str("OriginalError", err.Error()).Msg("Failed to retrieve ticker row from database")
		}
		tickerMap[ticker.Ticker] = &ticker
	}

	for _, r := range records {
		if ticker, ok := tickerMap[r.Ticker]; ok {
			r.CompositeFigi = ticker.CompositeFigi
		} else {
			if isValidExchange(r) {
				log.Warn().Str("Ticker", r.Ticker).Str("Exchange", r.Exchange).Msg("could not find composite figi for ticker")
			}
		}
	}

	return records
}

func SaveToDB(records []*ZacksRecord) {
	conn, err := pgx.Connect(context.Background(), viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Error().Str("OriginalError", err.Error()).Msg("Could not connect to database")
	}
	defer conn.Close(context.Background())

	cnt := 0
	for _, r := range records {
		if r.CompositeFigi != "" {
			_, err = conn.Exec(context.Background(),
				`INSERT INTO zacks_financials_v1 (
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
			) ON CONFLICT ON CONSTRAINT zacks_financials_v1_pkey
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
				log.Warn().Str("OriginalError", err.Error()).Msg("insert into db failed")
			} else {
				cnt++
			}
		}
	}

	log.Info().Int("NumRecords", cnt).Msg("records saved to DB")
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
				Str("EventDate", r.EventDateStr).Str("Ticker", r.Ticker).
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
