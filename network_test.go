package canopen

import (
	"testing"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

const TestPort string = "/dev/tty.usbserial-14220"

func TestSearch(t *testing.T) {
	transport := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	canBus := can.Bus{Transport: transport}

	if err := canBus.Open(); err != nil {
		t.Fatal(err)
	}
}
