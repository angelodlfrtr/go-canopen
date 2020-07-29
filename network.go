package canopen

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-canopen/utils"
	"github.com/google/uuid"
)

// Network represent the global nodes network
type Network struct {
	// mutex for FramesChans access
	sync.Mutex

	// Bus is the go-can bus
	Bus can.Bus

	// Nodes contain the network nodes
	Nodes map[int]*Node

	// FramesChans contains a list of chan when is sent each frames from network bus.
	FramesChans []*NetworkFramesChan

	// NMTMaster contain nmt control struct
	NMTMaster *NMTMaster

	// stopChan permit to stop network
	stopChan chan bool

	// running is network running
	running bool

	// BusReadErrChan @TODO on go-can
	BusReadErrChan chan error
}

// NewNetwork a new Network with given bus
func NewNetwork(bus can.Bus) (*Network, error) {
	// Create network
	netw := &Network{Bus: bus}

	// Set nmt and listen for nmt hearbeat messages
	netw.NMTMaster = NewNMTMaster(0, netw)

	return netw, nil
}

// Run listen handlers for frames on bus
func (network *Network) Run() error {
	if network.running {
		return nil
	}

	network.running = true
	network.stopChan = make(chan bool, 1)

	// Start network nmt master hearbeat listener
	if err := network.NMTMaster.ListenForHeartbeat(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-network.stopChan:
				// Stop goroutine
				return
			case frm := <-network.Bus.ReadChan():
				network.Lock()

				// Send frame to frames chans
				for _, ch := range network.FramesChans {
					ch.Publish(frm)
				}

				network.Unlock()
			}
		}
	}()

	return nil
}

// Stop handlers for frames on bus
func (network *Network) Stop() error {
	if !network.running {
		return nil
	}

	// Start network nmt master hearbeat listener
	if err := network.NMTMaster.UnlistenForHeartbeat(); err != nil {
		return err
	}

	// Stop each nodes
	for _, node := range network.Nodes {
		node.Stop()
	}

	network.stopChan <- true

	return nil
}

// Send a frame on network
func (network *Network) Send(arbID uint32, data []byte) error {
	frm := &can.Frame{
		ArbitrationID: arbID,
		DLC:           uint8(len(data)),
	}

	if len(data) > 8 {
		frm.DLC = uint8(8)
	}

	// Copy data to 8 byte array
	var arr [8]byte
	copy(arr[0:], data[:int(frm.DLC)])

	// Set data in frame
	frm.Data = arr

	// Write frame to serial port
	return network.Bus.Write(frm)
}

// AddNode add a node to the network
func (network *Network) AddNode(node *Node, objectDic *DicObjectDic, uploadEDS bool) *Node {
	if uploadEDS {
		// @TODO: download definition from node if true
		log.Fatal("uploading EDS not supported for now")
	}

	if node == nil {
		log.Fatal("Cannot use nil Node")
		return nil
	}

	// Set node network
	node.SetNetwork(network)

	// Set ObjectDic
	node.SetObjectDic(objectDic)

	// Set nmt and listen for nmt hearbeat messages
	node.NMTMaster = NewNMTMaster(node.ID, network)

	// Init node
	node.Init()

	// Start nmt master hearbeat listener
	if err := node.NMTMaster.ListenForHeartbeat(); err != nil {
		log.Fatalf("Failed to start nmt master on node %d with err %v", node.ID, err)
	}

	network.Lock()
	defer network.Unlock()
	// Initialize Nodes
	if network.Nodes == nil {
		network.Nodes = map[int]*Node{}
	}

	// Append node to network
	network.Nodes[node.ID] = node

	return node
}

// GetNode by node id. Return error if node dont exist in network.Nodes
func (network *Network) GetNode(nodeID int) (*Node, error) {
	network.Lock()
	defer network.Unlock()

	if node, ok := network.Nodes[nodeID]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("no node with id %d", nodeID)
}

// AcquireFramesChan create a new FrameChan
func (network *Network) AcquireFramesChan(filterFunc networkFramesChanFilterFunc) *NetworkFramesChan {
	network.Lock()
	defer network.Unlock()

	// Create frame chan
	chanID := uuid.Must(uuid.NewRandom()).String()
	frameChan := &NetworkFramesChan{
		ID:     chanID,
		Filter: filterFunc,
		C:      make(chan *can.Frame),
	}

	// Append network.FramesChans
	network.FramesChans = append(network.FramesChans, frameChan)

	return frameChan
}

// ReleaseFramesChan release (close) a FrameChan
func (network *Network) ReleaseFramesChan(id string) {
	network.Lock()
	defer network.Unlock()

	var framesChan *NetworkFramesChan
	var framesChanIndex *int

	for idx, fc := range network.FramesChans {
		if fc.ID == id {
			framesChan = fc
			idxx := idx
			framesChanIndex = &idxx
			break
		}
	}

	if framesChanIndex == nil {
		return
	}

	// Close chan
	close(framesChan.C)

	// Remove frameChan from network.FramesChans
	network.FramesChans = append(
		network.FramesChans[:*framesChanIndex],
		network.FramesChans[*framesChanIndex+1:]...,
	)
}

// Search send data to network and wait for nodes response
func (network *Network) Search(limit int, timeout time.Duration) ([]*Node, error) {
	// Nodes found
	nodes := make([]*Node, 0, limit)

	// Send ping for `limit` nodes
	reqData := []byte{0x40, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}
	for i := 1; i <= limit; i++ {
		if err := network.Send(uint32(0x600+i), reqData); err != nil {
			return nil, err
		}
	}

	// Canopen service
	services := []uint32{0x700, 0x580, 0x180, 0x280, 0x380, 0x480, 0x80}

	framesChan := network.AcquireFramesChan(nil)
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-timer.C:
			network.ReleaseFramesChan(framesChan.ID)
			return nodes, nil
		case frm := <-framesChan.C:
			service := frm.ArbitrationID & 0x780
			nodeID := int(frm.ArbitrationID & 0x7F)

			if nodeID != 0 {
				if !utils.ContainsUint32(services, uint32(nodeID)) {
					if utils.ContainsUint32(services, service) {
						// Append only if not already exist in nodes slice
						nodeExist := false
						for _, n := range nodes {
							if n.ID == nodeID {
								nodeExist = true
								break
							}
						}

						if nodeExist {
							continue
						}

						nNode := NewNode(nodeID, nil, nil)
						nodes = append(nodes, nNode)
					}
				}
			}
		}
	}
}
