package sdo

import (
	"github.com/angelodlfrtr/go-canopen"
)

// Client represent an SDO client
type Client struct {
	Node    *canopen.Node
	Network *canopen.Network
	RXCobID uint32
	TXCobID uint32
}

// SendRequest to network bus
func (client *Client) SendRequest(req []byte) error {
	return client.Network.Send(client.RXCobID, req)
}
