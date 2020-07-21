package canopen

import (
	"github.com/angelodlfrtr/go-can"
)

type networkFramesChanFilterFunc *(func(*can.Frame) bool)

// NetworkFramesChan contain a Chan, and ID and a Filter function
// Each FrameChan can have a filter function which return a boolean,
// and for each frame, the filter func is called. If func return true, the frame is returned,
// else dont send frame.
type NetworkFramesChan struct {
	ID     string
	C      chan *can.Frame
	Filter networkFramesChanFilterFunc
}

func (networkFramesCh *NetworkFramesChan) Publish(frm *can.Frame) {
	if networkFramesCh.Filter != nil {
		if (*networkFramesCh.Filter)(frm) {
			select {
			case networkFramesCh.C <- frm:
			default:
			}
		}

		return
	}

	select {
	case networkFramesCh.C <- frm:
	default:
	}
}
