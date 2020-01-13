package canopen

import (
	"errors"
	"reflect"
	"time"

	"github.com/angelodlfrtr/go-can/frame"
	"github.com/google/uuid"
)

const (
	MapPDONotValid   int = 1 << 31
	MapRTRNotAllowed int = 1 << 30
)

type PDOMapChangeChan struct {
	ID string
	C  chan []byte
}

// PDOMap @TODO : mutex
type PDOMap struct {
	PDONode   *PDONode
	ComRecord DicObject
	MapArray  DicObject

	Listening  bool
	Enabled    bool
	CobID      int
	RTRAllowed bool
	TransType  byte
	EventTimer byte

	Map map[int]DicObject

	OldData []byte
	Data    []byte

	Timestamp *time.Time
	Period    *time.Duration

	IsReceived bool

	ChangeChans []*PDOMapChangeChan
}

// NewPDOMap return a PDOMap initialized
func NewPDOMap(pdoNode *PDONode, comRecord, mapArray DicObject) *PDOMap {
	return &PDOMap{
		PDONode:    pdoNode,
		ComRecord:  comRecord,
		MapArray:   mapArray,
		RTRAllowed: true,
	}
}

// FindIndex find a object by index
func (m *PDOMap) FindIndex(idx int) DicObject {
	if ma, ok := m.Map[idx]; ok {
		return ma
	}

	return nil
}

// FindName find a object by name
func (m *PDOMap) FindName(name string) DicObject {
	var r DicObject

	for _, rr := range m.Map {
		if rr.GetName() == name {
			r = rr
			break
		}
	}

	return r
}

// GetTotalSize of a map
func (m *PDOMap) GetTotalSize() int {
	size := 0

	for _, rr := range m.Map {
		size += rr.GetDataLen()
	}

	return size
}

func (m *PDOMap) UpdateDataSize() {
	tSize := m.GetTotalSize()
	m.Data = make([]byte, 0, tSize)
}

func (m *PDOMap) SetData(data []byte) {
	m.OldData = m.Data
	m.Data = data
}

// Listen for changes on map from network
func (m *PDOMap) Listen() error {
	if m.CobID == 0 {
		return errors.New("call Read() on this map before listening")
	}

	if m.Listening {
		return nil
	}

	m.Listening = true

	now := time.Now()
	m.Timestamp = &now

	filterFunc := func(frm *frame.Frame) bool {
		return frm.ArbitrationID == uint32(m.CobID)
	}

	framesChan := m.PDONode.Node.Network.AcquireFramesChan(&filterFunc)

	go func() {
		for {
			// Stop routine if listening == false
			if !m.Listening {
				return
			}

			select {
			case frm := <-framesChan.C:
				m.IsReceived = true
				m.SetData(frm.GetData())
				// @TODO m.Period = frm.Timestamp - m.Timestamp;
				now := time.Now()
				m.Timestamp = &now

				// If data changed
				if !reflect.DeepEqual(m.OldData, m.Data) {
					for _, changeChan := range m.ChangeChans {
						changeChan.C <- m.Data
					}
				}
			default:
				continue
			}
		}
	}()

	return nil
}

// Unlisten for changes on map from network
func (m *PDOMap) Unlisten() {
	m.Listening = false
}

// AcquireChangesChan create a new PDOMapChangeChan
func (m *PDOMap) AcquireChangesChan() *PDOMapChangeChan {
	// Create frame chan
	chanID := uuid.Must(uuid.NewRandom()).String()
	changesChan := &PDOMapChangeChan{
		ID: chanID,
		C:  make(chan []byte),
	}

	// Append m.ChangeChans
	m.ChangeChans = append(m.ChangeChans, changesChan)

	return changesChan
}

// ReleaseChangesChan release (close) a PDOMapChangeChan
func (m *PDOMap) ReleaseChangesChan(id string) error {
	var changesChan *PDOMapChangeChan
	var changesChanIndex *int

	for idx, fc := range m.ChangeChans {
		if fc.ID == id {
			changesChan = fc
			idxx := idx
			changesChanIndex = &idxx
			break
		}
	}

	if changesChanIndex == nil {
		return errors.New("no PDOMapChangeChan found with specified ID")
	}

	// Close chan
	close(changesChan.C)

	// Remove frameChan from network.FramesChans
	m.ChangeChans = append(
		m.ChangeChans[:*changesChanIndex],
		m.ChangeChans[*changesChanIndex+1:]...,
	)

	return nil
}

// Read map values
func (m *PDOMap) Read() error {
	// Get COB ID
	if err := m.ComRecord.FindIndex(1).Read(); err != nil {
		return err
	}

	cobID := int(*m.ComRecord.FindIndex(1).GetUintVal())
	m.CobID = cobID

	// Is enabled
	m.Enabled = (cobID & MapPDONotValid) == 0

	// Is RTRAllowed
	m.RTRAllowed = (cobID & MapRTRNotAllowed) == 0

	// Get Trans type
	if err := m.ComRecord.FindIndex(2).Read(); err != nil {
		return err
	}

	transType := *m.ComRecord.FindIndex(2).GetUintVal()
	m.TransType = byte(transType)

	// Get EventTimer
	if transType > 254 {
		comr := m.ComRecord.FindIndex(5)

		if comr != nil {
			if err := comr.Read(); err != nil {
				return err
			}

			m.EventTimer = *comr.GetByteVal()
		}
	}

	// Init m.Map
	m.Map = make(map[int]DicObject)
	offset := 0

	// Nof entries
	if err := m.MapArray.FindIndex(0).Read(); err != nil {
		return err
	}

	nofEntries := int(m.MapArray.FindIndex(0).GetData()[0])

	for i := 1; i <= (nofEntries + 1); i++ {
		ii := uint16(i)
		if err := m.MapArray.FindIndex(ii).Read(); err != nil {
			return err
		}

		val := *m.MapArray.FindIndex(ii).GetUintVal()

		index := uint16(val >> 16)
		subindex := uint16((val >> 8) & 0xFF)
		size := val & 0xFF

		if size == 0 {
			continue
		}

		dicVar := m.PDONode.Node.ObjectDic.FindIndex(index)
		// Instead of dicVar.Size = size @TODO: use uint64 for size
		dicVar.SetSize(int(size))

		// Set sdo client
		dicVar.SetSDO(m.PDONode.Node.SDOClient)

		if !dicVar.IsDicVariable() {
			dicVar = dicVar.FindIndex(subindex)
		}

		dicVar.SetOffset(offset)
		// @TODO: check working
		m.Map[i] = dicVar

		// @TODO: use uint64
		offset += int(size)
	}

	m.UpdateDataSize()

	return m.Listen()
}

// Save pdo map
// @TODO: Not Working, DO NOT USE
func (m *PDOMap) Save() error {
	if true {
		return errors.New("PDOMap.Save() not implemented")
	}

	// Get COB ID
	if err := m.ComRecord.FindIndex(1).Read(); err != nil {
		return err
	}

	cobID := int(*m.ComRecord.FindIndex(1).GetUintVal())
	m.CobID = cobID

	// Is enabled
	m.Enabled = (cobID & MapPDONotValid) == 0

	// Is RTRAllowed
	m.RTRAllowed = (cobID & MapRTRNotAllowed) == 0

	// Setting COB-ID 0x%X and temporarily disabling PDO
	m.ComRecord.FindIndex(1).SetData([]byte{byte(cobID | MapPDONotValid)})
	if err := m.ComRecord.FindIndex(1).Save(); err != nil {
		return err
	}

	// Set transType
	if m.TransType != 0 {
		m.ComRecord.FindIndex(2).SetData([]byte{m.TransType})
		if err := m.ComRecord.FindIndex(2).Save(); err != nil {
			return err
		}
	}

	if m.EventTimer != 0 {
		m.ComRecord.FindIndex(5).SetData([]byte{m.EventTimer})
		if err := m.ComRecord.FindIndex(5).Save(); err != nil {
			return err
		}
	}

	if len(m.Map) > 0 {
		for i, dicVar := range m.Map {
			subIndex := i + 1
			// nolint
			val := ((((dicVar.GetIndex() << 16) | uint16(dicVar.GetSubIndex())) << 8) | uint16(dicVar.GetDataLen()))
			mm := m.MapArray.FindIndex(uint16(subIndex))
			mm.SetUintVal(uint64(val))

			if err := mm.Save(); err != nil {
				return err
			}
		}

		mapLen := len(m.Map)
		m.MapArray.FindIndex(0).SetIntVal(int64(mapLen))
		if err := m.MapArray.FindIndex(0).Save(); err != nil {
			return err
		}

		m.UpdateDataSize()
	}

	return nil
}

// RebuildData rebuild map data object from map variables
func (m *PDOMap) RebuildData() {
	data := make([]byte, m.GetTotalSize()/8)

	for _, dicVar := range m.Map {
		dicVarData := dicVar.GetData()
		if len(dicVarData) == 0 {
			continue
		}

		dicOffset := dicVar.GetOffset() / 8

		for i := 0; i < len(dicVarData); i++ {
			offset := i + dicOffset
			data[offset] = dicVarData[i]
		}
	}

	m.SetData(data)
}

// Transmit map data
func (m *PDOMap) Transmit(rebuild bool) error {
	if rebuild {
		m.RebuildData()
	}

	return m.PDONode.Node.Network.Send(uint32(m.CobID), m.Data)
}
