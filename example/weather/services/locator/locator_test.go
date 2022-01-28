package locator

import (
	"context"
	"fmt"
	"testing"

	"github.com/goadesign/clue/example/weather/services/locator/clients/ipapi"
)

func TestGetIPLocation(t *testing.T) {
	var (
		testLat     = 37.4224764
		testLong    = -122.0842499
		testCity    = "Mountain View"
		testRegion  = "CA"
		testCountry = "United States"
	)
	// Create mock call sequence with first successful call returning an IP
	// location then failing.
	ipc := ipapi.NewMock(t)
	ipc.AddGetLocationFunc(func(ctx context.Context, ip string) (*ipapi.WorldLocation, error) {
		return &ipapi.WorldLocation{testLat, testLong, testCity, testRegion, testCountry}, nil
	})
	ipc.AddGetLocationFunc(func(ctx context.Context, ip string) (*ipapi.WorldLocation, error) {
		return nil, fmt.Errorf("test failure")
	})

	// Create locator service.
	s := New(ipc)

	// Call service, first call should succeed.
	l, err := s.GetLocation(context.Background(), "8.8.8.8")
	if err != nil {
		t.Errorf("got error %v, expected nil", err)
	}
	if l.Lat != testLat {
		t.Errorf("got lat %v, expected %f", l.Lat, testLat)
	}
	if l.Long != testLong {
		t.Errorf("got long %v, expected %f", l.Long, testLong)
	}
	if l.City != testCity {
		t.Errorf("got city %q, expected %q", l.City, testCity)
	}
	if l.Region != testRegion {
		t.Errorf("got region code %q, expected %q", l.Region, testRegion)
	}

	// Call service, second call should fail.
	_, err = s.GetLocation(context.Background(), "8.8.8.8")
	if err == nil {
		t.Errorf("got nil, expected error")
	}

	// Make sure all calls were made.
	if ipc.HasMore() {
		t.Errorf("expected all calls to be made")
	}
}
