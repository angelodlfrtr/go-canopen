package canopen

import (
	"testing"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

const TestPort string = "/dev/tty.usbserial-14220"

func TestSearch(t *testing.T) {
	transport := &transports.USBCanAnalyzer{
		Port:     TestPort,
		BaudRate: 2000000,
	}

	bus := can.Bus{Transport: transport}

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	network := NewNetwork(bus)

	// Run search (in my case), node ids a returned after ~~500ms
	// So be secure with timeout
	timeout := time.Duration(3) * time.Second
	nodes, err := network.Search(256, timeout)

	if err != nil {
		t.Fatal(err)
	}

	// Expect a least one node in results
	if len(nodes) == 0 {
		t.Fatal("No nodes found")
	}

	t.Log(nodes)
}
