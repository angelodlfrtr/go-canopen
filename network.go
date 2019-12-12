package canopen

import (
	"fmt"
	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/frame"
	"github.com/angelodlfrtr/go-canopen/dic"
	"github.com/thoas/go-funk"
	"log"
	"time"
)

// Network represent the global nodes network
type Network struct {
	// Bus is the go-can bus
	Bus can.Bus

	// Nodes contain the network nodes
	Nodes map[int]*Node
}

// NewNetwork a new Network with given bus
func NewNetwork(bus can.Bus) *Network {
	return &Network{Bus: bus}
}

// Listen for canopen message on bus
func (network *Network) Listen() {
	// @TODO
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
	copy(arr[:], data[:8])

	frm.Data = arr

	return network.Bus.Write(frm)
}

// AddNode add a node to the network
func (network *Network) AddNode(n Node, objectDic *dic.ObjectDic, uploadEDS bool) *Node {
	if uploadEDS {
		// @TODO: download definition from node if true
		log.Fatal("Uploading EDS not supported for now")
	}

	// Set node network
	n.SetNetwork(network)

	// Set ObjectDic
	n.SetObjectDic(objectDic)

	// Append node to network
	network.Nodes[n.ID] = &n

	return &n
}

// GetNode by node id. Return error if node dont exist in network.Nodes
func (network *Network) GetNode(nodeId int) (*Node, error) {
	if node, ok := network.Nodes[nodeId]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("No node with id %d", nodeId)
}

// Search send data to network and wait for nodes response
func (network *Network) Search(limit int, timeout time.Duration) ([]*Node, error) {
	// Nodes found
	var nodes []*Node

	reqData := []byte{0x40, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00}

	// Send ping for `limit` nodes
	for i := 0; i <= limit+1; i++ {
		if err := network.Send(uint32(0x600+1), reqData); err != nil {
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

			if !funk.Contains(services, uint32(nodeID)) && nodeID != 0 && funk.Contains(services, service) {
				nodes = append(nodes, &Node{ID: nodeID})
			}
		}
	}

	return nodes, nil
}
