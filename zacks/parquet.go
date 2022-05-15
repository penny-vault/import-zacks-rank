package zacks

import (
	"github.com/rs/zerolog/log"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

func SaveToParquet(records []*ZacksRecord, fn string) error {
	var err error

	fh, err := local.NewLocalFileWriter(fn)
	if err != nil {
		log.Error().Err(err).Str("FileName", fn).Msg("cannot create local file")
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
		log.Error().Err(err).Msg("Parquet write failed")
		return err
	}

	log.Info().Int("NumRecords", len(records)).Msg("Parquet write finished")
	return nil
}
