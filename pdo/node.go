package pdo

import (
	"github.com/angelodlfrtr/go-canopen"
)

type Node struct {
	Node    *canopen.Node
	Network *canopen.Network
	RX      *Maps
	TX      *Maps
}
