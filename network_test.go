package canopen

import (
	"os"
	"testing"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

func getTestPort() string {
	if a := os.Getenv("CAN_TEST_PORT"); len(a) > 0 {
		return a
	}

	return "/dev/tty.some-usbserial"
}

func getNetwork() (*Network, error) {
	testPort := getTestPort()
	transport := &transports.USBCanAnalyzer{
		Port:     testPort,
		BaudRate: 2000000,
	}

	bus := can.Bus{Transport: transport}

	if err := bus.Open(); err != nil {
		return nil, err
	}

	netw, err := NewNetwork(bus)
	if err != nil {
		return nil, err
	}

	if err := netw.Run(); err != nil {
		return nil, err
	}

	return netw, nil
}

func searchNodes() ([]*Node, error) {
	network, err := getNetwork()
	if err != nil {
		return nil, err
	}

	// Run search (in my case), node ids a returned after ~~500ms
	// So be secure with timeout
	timeout := time.Duration(2) * time.Second
	nodes, err := network.Search(256, timeout)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func TestSend(t *testing.T) {
	network, err := getNetwork()
	if err != nil {
		t.Fatal(err)
	}

	err = network.Send(uint32(0x01), []byte{0x0, 0x0, 0x0})

	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	nodes, err := searchNodes()
	if err != nil {
		t.Fatal(err)
	}

	// Expect a least one node in results
	if len(nodes) == 0 {
		t.Fatal("No nodes found")
	}

	t.Log(nodes)
}

func TestAddNode(t *testing.T) {
	network := &Network{}
	node := &Node{ID: 1}
	network.AddNode(node, nil, false)

	if len(network.Nodes) != 1 {
		t.Fatal("Invalid network.Nodes len")
	}
}

func TestGetNode(t *testing.T) {
	network := &Network{}
	node := &Node{ID: 1}
	network.AddNode(node, nil, false)

	if len(network.Nodes) != 1 {
		t.Fatal("Invalid network.Nodes len")
	}

	nodeGot, err := network.GetNode(node.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nodeGot == nil {
		t.Fatal("Node not found")
	}
}

func TestAll(t *testing.T) {
	testPort := getTestPort()
	transport := &transports.USBCanAnalyzer{
		Port:     testPort,
		BaudRate: 2000000,
	}

	bus := can.Bus{Transport: transport}

	if err := bus.Open(); err != nil {
		t.Fatal(err)
	}

	network, err := NewNetwork(bus)
	if err != nil {
		t.Fatal(err)
	}

	if err := network.Run(); err != nil {
		t.Fatal(err)
	}

	// Load object dic
	objectDicFilePath := os.Getenv("CAN_TEST_EDS")
	if len(objectDicFilePath) == 0 {
		t.Fatal("Invalid object dic file path")
	}

	// Parse eds file
	dic, err := DicEDSParse(objectDicFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Run search node ids a returned after ~500ms in my case
	// So be secure with timeout
	searchTimeout := time.Duration(2) * time.Second
	nodes, err := network.Search(256, searchTimeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) == 0 {
		t.Fatal("No nodes found")
	}

	for _, n := range nodes {
		network.AddNode(n, dic, false)
	}

	// node := nodes[0]
	node, _ := network.GetNode(41)
	t.Log("Handle node ID", node.ID)

	// Read node PDO
	if err := node.PDONode.Read(); err != nil {
		t.Fatal(err)
	}

	// Listen PDO
	eleMap := node.PDONode.FindName("SomeName")
	changesChan := eleMap.AcquireChangesChan()

	// Wait for any change
	timer := time.NewTicker(10 * time.Second)

	for {
		select {
		case data := <-changesChan.C:
			t.Log(data)
		case <-timer.C:
			t.Log("Done")
			return
		}
	}
}
