package pdo

import (
	"time"

	"github.com/angelodlfrtr/go-canopen/dic"
)

type Map struct {
	Node      *Node
	ComRecord dic.Object
	MapArray  dic.Object

	Enabled    bool
	CobID      int
	RTRAllowed bool
	TransType  byte
	EventTimer byte

	Map map[int]dic.Object

	OldData []byte
	Data    []byte

	Timestamp *time.Time
	Period    *time.Duration

	IsReceived bool
}

func NewMap(pdoNode *Node, comRecord, mapArray dic.Object) *Map {
	pdoMap := &Map{
		Node:       pdoNode,
		ComRecord:  comRecord,
		MapArray:   mapArray,
		RTRAllowed: true,
	}

	return pdoMap
}
