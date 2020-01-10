package canopen

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/frame"
	"github.com/angelodlfrtr/go-canopen/utils"
	"github.com/google/uuid"
)

type networkFramesChanFilterFunc *(func(*frame.Frame) bool)

// FrameChan contain a Chan, and ID and a Filter function
// Each FrameChan can have a filter function which return a boolean,
// and for each frame, the filter func is called. If func return true, the frame is returned,
// else dont send frame.
type NetworkFramesChan struct {
	ID     string
	Chan   chan *frame.Frame
	Filter networkFramesChanFilterFunc
}

// Network represent the global nodes network
type Network struct {
	// Bus is the go-can bus
	Bus can.Bus

	// Nodes contain the network nodes
	Nodes map[int]*Node

	// FramesChans contains a list of chan when is sent each frames from network bus.
	FramesChans []*NetworkFramesChan

	// NMTMaster contain nmt control struct
	NMTMaster *NMTMaster

	// mutex for FramesChans access
	mutex sync.Mutex

	// listening is network reading datas on can bus
	listening bool

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
	// @TODO: check bus is opened

	// Start network nmt master hearbeat listener
	if err := network.NMTMaster.ListenForHeartbeat(); err != nil {
		return err
	}

	// Set as listening
	network.listening = true

	go func() {
		for {
			if !network.listening {
				// Stop loop and goroutine
				break
			}

			// Read frame
			frm := &frame.Frame{}
			ok, err := network.Bus.Read(frm)

			if err != nil {
				network.BusReadErrChan <- err
				continue
			}

			// If not data continue
			if !ok {
				continue
			}

			// Send frame to frames chans
			for _, ch := range network.FramesChans {
				if ch.Filter != nil {
					if (*ch.Filter)(frm) {
						ch.Chan <- frm
					}
				} else {
					ch.Chan <- frm
				}
			}
		}
	}()

	return nil
}

// Stop handlers for frames on bus
func (network *Network) Stop() error {
	if !network.listening {
		return nil
	}

	// Start network nmt master hearbeat listener
	if err := network.NMTMaster.UnlistenForHeartbeat(); err != nil {
		return err
	}

	network.listening = false

	// @TODO: stop all nmtmasters, and all chan listeners

	return nil
}

func (network *Network) Send(arbID uint32, data []byte) error {
	frm := &frame.Frame{
		ArbitrationID: arbID,
		DLC:           uint8(len(data)),
	}

	if len(data) > 8 {
		frm.DLC = uint8(8)
	}

	// Copy data to 8 byte array
	var arr [8]byte
	copy(arr[:], data[:int(frm.DLC)])

	frm.Data = arr

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
	}

	// Set node network
	node.SetNetwork(network)

	// Set ObjectDic
	node.SetObjectDic(objectDic)

	// Init node
	node.Init()

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
	if node, ok := network.Nodes[nodeID]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("no node with id %d", nodeID)
}

// AcquireFramesChan create a new FrameChan
func (network *Network) AcquireFramesChan(filterFunc networkFramesChanFilterFunc) *NetworkFramesChan {
	network.mutex.Lock()
	defer network.mutex.Unlock()

	// Create frame chan
	chanID := uuid.Must(uuid.NewRandom()).String()
	frameChan := &NetworkFramesChan{
		ID:     chanID,
		Filter: filterFunc,
		Chan:   make(chan *frame.Frame),
	}

	// Append network.FramesChans
	network.FramesChans = append(network.FramesChans, frameChan)

	return frameChan
}

// ReleaseFramesChan release (close) a FrameChan
func (network *Network) ReleaseFramesChan(id string) error {
	network.mutex.Lock()
	defer network.mutex.Unlock()

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
		return errors.New("no FrameChan found with specified ID")
	}

	// Close chan
	close(framesChan.Chan)

	// Remove frameChan from network.FramesChans
	network.FramesChans = append(
		network.FramesChans[:*framesChanIndex],
		network.FramesChans[*framesChanIndex+1:]...,
	)

	return nil
}

// Search send data to network and wait for nodes response
func (network *Network) Search(limit int, timeout time.Duration) ([]*Node, error) {
	// Nodes found
	nodes := make([]*Node, 0, limit)

	reqData := []byte{0x40, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Send ping for `limit` nodes
	for i := 0; i <= limit+1; i++ {
		if err := network.Send(uint32(0x600+i), reqData); err != nil {
			return nodes, err
		}
	}

	// Canopen service
	services := []uint32{0x700, 0x580, 0x180, 0x280, 0x380, 0x480, 0x80}

	// Handle pongs
	start := time.Now()

	for {
		if time.Since(start) > timeout {
			break
		}

		frm := &frame.Frame{}
		ok, err := network.Bus.Read(frm)

		if err != nil {
			return nil, err
		}

		if ok {
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

	return nodes, nil
}
