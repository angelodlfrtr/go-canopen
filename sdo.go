package canopen

import (
	"errors"
	"time"

	"github.com/angelodlfrtr/go-can/frame"
)

const (
	SDORequestUpload    uint8 = 2 << 5
	SDOResponseUpload   uint8 = 2 << 5
	SDORequestDownload  uint8 = 1 << 5
	SDOResponseDownload uint8 = 3 << 5

	SDORequestSegmentUpload    uint8 = 3 << 5
	SDOResponseSegmentUpload   uint8 = 0 << 5
	SDORequestSegmentDownload  uint8 = 0 << 5
	SDOResponseSegmentDownload uint8 = 1 << 5

	SDOExpedited     uint8 = 0x2
	SDOSizeSpecified uint8 = 0x1
	SDOToggleBit     uint8 = 0x10
	SDONoMoreData    uint8 = 0x1
)

// Client represent an SDO client
type SDOClient struct {
	Node    *Node
	RXCobID uint32
	TXCobID uint32
}

func NewSDOClient(node *Node) *SDOClient {
	return &SDOClient{
		Node:    node,
		RXCobID: uint32(0x600 + node.ID),
		TXCobID: uint32(0x580 + node.ID),
	}
}

// SendRequest to network bus
func (sdoClient *SDOClient) SendRequest(req []byte) error {
	return sdoClient.Node.Network.Send(uint32(sdoClient.RXCobID), req)
}

// Send message and optionaly wait for response
func (sdoClient *SDOClient) Send(
	req []byte,
	expectFunc networkFramesChanFilterFunc,
	timeout *time.Duration,
) (*frame.Frame, error) {
	// Set default timeout
	if timeout == nil {
		dtm := time.Duration(1) * time.Second
		timeout = &dtm
	}

	var framesChan *NetworkFramesChan

	// If response wanted, require data chan to network
	if expectFunc != nil {
		framesChan = sdoClient.Node.Network.AcquireFramesChan(expectFunc)
	}

	if err := sdoClient.SendRequest(req); err != nil {
		return nil, err
	}

	// If no response wanted, just return nothing
	if expectFunc == nil {
		return nil, nil
	}

	// Wait for response frame
	var frm *frame.Frame
	start := time.Now()

	for {
		if time.Since(start) > *timeout {
			return nil, errors.New("timeout exceeded")
		}

		fr, ok := <-framesChan.Chan
		if ok {
			frm = fr
			break
		}
	}

	// Release data chan
	if err := sdoClient.Node.Network.ReleaseFramesChan(framesChan.ID); err != nil {
		return frm, err
	}

	return frm, nil
}

// Read sdo
func (sdoClient *SDOClient) Read(index uint16, subIndex uint8) ([]byte, error) {
	reader := NewSDOReader(sdoClient, index, subIndex)
	return reader.ReadAll()
}
