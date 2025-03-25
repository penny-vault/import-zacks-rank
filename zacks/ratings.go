package zacks

import (
	"context"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"

	"github.com/rs/zerolog/log"
)

func LoadRatings(ratingsData []byte, dateStr string, limit int) []*ZacksRecord {
	records := []*ZacksRecord{}

	stringData := string(ratingsData[:])
	stringData = strings.ReplaceAll(stringData, `"NA"`, `"0"`)

	if err := gocsv.UnmarshalString(stringData, &records); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal byte data")
		return make([]*ZacksRecord, 0)
	}

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Error().Err(err).Str("DateStr", dateStr).Msg("cannot parse dateStr")
	}

	// cleanup records
	for _, r := range records {
		r.Ticker = strings.ReplaceAll(r.Ticker, ".", "/")

		// set event date
		r.EventDateStr = dateStr
		r.EventDate = date
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
	conn, err := pgx.Connect(context.Background(), viper.GetString("database.url"))
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to database")
	}
	defer conn.Close(context.Background())

	// build a list of all active records that have composite figi's
	tickerMap := make(map[string]*Ticker)

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
