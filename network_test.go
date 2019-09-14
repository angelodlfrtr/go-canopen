package canopen

import (
	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
	"testing"
)

func TestSearch(t *testing.T) {
	transport := &transports.USBCanAnalyzer{
		BaudRate: 2000000,
		Port:     "/dev/ttyusblala",
	}

	canBus := can.Bus{Transport: transport}

	if err := canBus.Open(); err != nil {
		t.Fatal(err)
	}
}
