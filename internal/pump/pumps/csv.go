package pumps

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/pachirode/iam_study/internal/pump/analytics"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

type CSVPump struct {
	csvConf *CSVConf
	CommonPumpConfig
}

func (c *CSVPump) New() Pump {
	newPump := CSVPump{}

	return &newPump
}

func (c *CSVPump) GetName() string {
	return "CSV Pump"
}

func (c *CSVPump) Init(conf interface{}) error {
	c.csvConf = &CSVConf{}
	err := mapstructure.Decode(conf, &c.csvConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	err = os.MkdirAll(c.csvConf.CSVDir, 0o777)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("CSV Initialized")

	return nil
}

func (c *CSVPump) WriteData(ctx context.Context, data []interface{}) error {
	curTime := time.Now()
	fName := fmt.Sprintf("%d-%s-%d-%d.csv", curTime.Year(), curTime.Month().String(), curTime.Day(), curTime.Hour())
	fName = path.Join(c.csvConf.CSVDir, fName)

	var outFile *os.File
	var appendHeader bool

	if _, err := os.Stat(fName); os.IsNotExist(err) {
		var createErr error
		outFile, createErr = os.Create(fName)
		if createErr != nil {
			log.Errorf("Failed to create new CSV file: %s", createErr.Error())
		}
		appendHeader = true
	} else {
		var appendErr error
		outFile, appendErr = os.OpenFile(fName, os.O_APPEND|os.O_WRONLY, 0o600)
		if appendErr != nil {
			log.Errorf("Failed to open CSV file: %s", appendErr.Error())
		}
	}

	defer outFile.Close()
	writer := csv.NewWriter(outFile)

	if appendHeader {
		startRecord := analytics.AnalyticsRecord{}
		headers := startRecord.GetFieldNames()

		err := writer.Write(headers)
		if err != nil {
			log.Errorf("Failed to write file headers: %s", err.Error())

			return errors.Wrap(err, "Failed to write file headers")
		}
	}

	for _, v := range data {
		decoded, _ := v.(analytics.AnalyticsRecord)

		toWrite := decoded.GetLineValues()
		err := writer.Write(toWrite)
		if err != nil {
			log.Error("File write failed!")
			log.Error(err.Error())
		}
	}

	writer.Flush()

	return nil
}
