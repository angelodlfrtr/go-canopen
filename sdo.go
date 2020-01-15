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
	return sdoClient.Node.Network.Send(sdoClient.RXCobID, req)
}

// FindName find an sdo object from object dictionary by name
func (sdoClient *SDOClient) FindName(name string) DicObject {
	if ob := sdoClient.Node.ObjectDic.FindName(name); ob != nil {
		ob.SetSDO(sdoClient)
		return ob
	}

	return nil
}

// Send message and optionaly wait for response
func (sdoClient *SDOClient) Send(
	req []byte,
	expectFunc networkFramesChanFilterFunc,
	timeout *time.Duration,
	retryCount *int,
) (*frame.Frame, error) {
	// If no response wanted, just send and return
	if expectFunc == nil {
		if err := sdoClient.SendRequest(req); err != nil {
			return nil, err
		}

		return nil, nil
	}

	// Set default timeout
	if timeout == nil {
		dtm := time.Duration(300) * time.Millisecond
		timeout = &dtm
	}

	if retryCount == nil {
		rtc := 3
		retryCount = &rtc
	}

	framesChan := sdoClient.Node.Network.AcquireFramesChan(expectFunc)

	// Retry loop
	remainingCount := *retryCount
	var frm *frame.Frame

	for {
		if remainingCount == 0 {
			break
		}

		if err := sdoClient.SendRequest(req); err != nil {
			return nil, err
		}

		timer := time.NewTicker(*timeout)

		select {
		case <-timer.C:
			// Double timeout for each retry
			newTimeout := *timeout * 2
			timeout = &newTimeout
		case fr := <-framesChan.C:
			frm = fr
		}

		timer.Stop()
		remainingCount--

		if frm != nil {
			break
		}
	}

	// Release data chan
	sdoClient.Node.Network.ReleaseFramesChan(framesChan.ID)

	// If no frm, timeout execeded
	if frm == nil {
		return nil, errors.New("timeout execeded")
	}

	return frm, nil
}

// Read sdo
func (sdoClient *SDOClient) Read(index uint16, subIndex uint8) ([]byte, error) {
	reader := NewSDOReader(sdoClient, index, subIndex)
	return reader.ReadAll()
}

// Write sdo
func (sdoClient *SDOClient) Write(index uint16, subIndex uint8, forceSegment bool, data []byte) error {
	writer := NewSDOWriter(sdoClient, index, subIndex, forceSegment)
	return writer.Write(data)
}
