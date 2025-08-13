package pumps

import (
	"context"

	"github.com/pachirode/iam_study/pkg/log"
)

type DummyPump struct {
	CommonPumpConfig
}

func (p *DummyPump) New() Pump {
	newPump := DummyPump{}

	return &newPump
}

func (p *DummyPump) GetName() string {
	return "Dummy Pump"
}

func (p *DummyPump) Init(conf interface{}) error {
	log.Debug("Dummy Initialized")

	return nil
}

func (p *DummyPump) WriteData(ctx context.Context, data []interface{}) error {
	log.Infof("Writing %d records", len(data))

	return nil
}
