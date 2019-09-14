package canopen

import (
	"github.com/angelodlfrtr/go-canopen/dic"
)

// Node is a canopen node
type Node struct {
	// Each node has an id, which is ArbitrationID & 0x7F
	ID int

	// Network the node network
	Network *Network

	// ObjectDic is the node object dictionnary
	ObjectDic *dic.ObjectDic
}

// SetNetwork set node.Network to the desired network
func (node *Node) SetNetwork(network *Network) {
	node.Network = network
}

// SetObjectDic set node.ObjectDic to the desired ObjectDic
func (node *Node) SetObjectDic(objectDic *dic.ObjectDic) {
	node.ObjectDic = objectDic
}
