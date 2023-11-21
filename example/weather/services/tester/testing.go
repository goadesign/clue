package tester

import (
	"context"

	"goa.design/clue/example/weather/services/tester/clients/forecaster"
	"goa.design/clue/example/weather/services/tester/clients/locator"
	"goa.design/clue/log"
)

type (
	Service struct {
		lc locator.Client
		fc forecaster.Client
	}
)

func New(lc locator.Client, fc forecaster.Client) *Service {
	return &Service{
		lc: lc,
		fc: fc,
	}
}

func logError(ctx context.Context, err error) error {
	log.Error(ctx, err)
	return err
}
