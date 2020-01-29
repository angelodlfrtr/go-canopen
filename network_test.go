package canopen

import (
	"fmt"
	"os"
	"sync"
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

	// Run search node ids a returned after ~500ms in my case
	// So be secure with timeout
	searchTimeout := time.Duration(1) * time.Second
	nodes, err := network.Search(256, searchTimeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) == 0 {
		t.Fatal("No nodes found")
	}

	fmt.Println("Nodes found", len(nodes))

	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, n := range nodes {
		wg.Add(1)

		go func(node *Node) {

			// Parse eds file
			dic := DicMustParse(DicEDSParse(objectDicFilePath))

			network.AddNode(node, dic, false)

			fmt.Println("Reading PDO")

			if err := node.PDONode.Read(); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}

			// node := nodes[0]
			// node, _ := network.GetNode(4)
			t.Log("PDO NODE readed ID", node.ID)

			wg.Done()
		}(n)
	}

	wg.Wait()

	select {
	case e := <-errChan:
		t.Fatal(e)
	default:
		close(errChan)
	}

	fmt.Println("Done")

	// Stop network
	network.Stop()
	fmt.Println("Network stopped")

	// Stop bus
	bus.Close()
	fmt.Println("Bus closed")
}

func TestReboot(t *testing.T) {
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

	// Run search node ids a returned after ~500ms in my case
	// So be secure with timeout
	searchTimeout := time.Duration(4) * time.Second
	nodes, err := network.Search(256, searchTimeout)
	if err != nil {
		t.Fatal(err)
	}

	if len(nodes) == 0 {
		t.Fatal("No nodes found")
	}

	fmt.Println("Nodes found", len(nodes))

	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, n := range nodes {
		wg.Add(1)

		go func(node *Node) {

			// Parse eds file
			dic := DicMustParse(DicEDSParse(objectDicFilePath))

			network.AddNode(node, dic, false)

			fmt.Println("Reading PDO")

			if err := node.PDONode.Read(); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}

			// node := nodes[0]
			// node, _ := network.GetNode(4)
			t.Log("PDO NODE readed ID", node.ID)

			wg.Done()
		}(n)
	}

	wg.Wait()

	select {
	case e := <-errChan:
		t.Fatal(e)
	default:
		close(errChan)
	}

	// Reboot first node
	node := nodes[0]

	fmt.Println("Rebooting")
	node.NMTMaster.SetState("RESET")

	done := make(chan bool)

	go func() {
		timeout := 20 * time.Second
		node.NMTMaster.WaitForBootup(&timeout)
		fmt.Println("Rebooted")
		done <- true
	}()

	<-done
}
