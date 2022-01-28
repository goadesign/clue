package locator

import (
	"context"

	"github.com/goadesign/clue/example/weather/services/locator/clients/ipapi"
	genlocator "github.com/goadesign/clue/example/weather/services/locator/gen/locator"
)

type (
	// Service is the locator service implementation.
	Service struct {
		ipc ipapi.Client
	}
)

// New instantiates a new locator service.
func New(ipc ipapi.Client) *Service {
	return &Service{ipc: ipc}
}

// GetLocation returns the location for the given IP address.
func (s *Service) GetLocation(ctx context.Context, ip string) (*genlocator.WorldLocation, error) {
	l, err := s.ipc.GetLocation(ctx, ip)
	if err != nil {
		return nil, err
	}
	lval := genlocator.WorldLocation(*l)
	return &lval, nil
}
