package canopen

import (
	"time"
)

const (
	MapPDONotValid   int = 1 << 31
	MapRTRNotAllowed int = 1 << 30
)

type PDOMap struct {
	PDONode   *PDONode
	ComRecord DicObject
	MapArray  DicObject

	Enabled    bool
	CobID      int
	RTRAllowed bool
	TransType  byte
	EventTimer byte

	Map map[int]DicObject

	// Maybe use bytes.Buffer
	OldData []byte
	Data    []byte

	Timestamp *time.Time
	Period    *time.Duration

	IsReceived bool
}

func NewPDOMap(pdoNode *PDONode, comRecord, mapArray DicObject) *PDOMap {
	pdoMap := &PDOMap{
		PDONode:    pdoNode,
		ComRecord:  comRecord,
		MapArray:   mapArray,
		RTRAllowed: true,
	}

	return pdoMap
}

// Find a map by index
func (m *PDOMap) FindIndex(idx int) DicObject {
	if ma, ok := m.Map[idx]; ok {
		return ma
	}

	return nil
}

// Find a map by name
func (m *PDOMap) FindByName(name string) DicObject {
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

func (m *PDOMap) Listen() error {
	// @TODO
	return nil
}

// Read a map
func (m *PDOMap) Read() error {
	// Get COB ID
	if err := m.ComRecord.FindIndex(1).Read(); err != nil {
		return err
	}

	cobID := int(*m.ComRecord.FindIndex(1).GetIntVal())
	m.CobID = cobID

	// Is enabled
	m.Enabled = (cobID & MapPDONotValid) == 0

	// Is RTRAllowed
	m.RTRAllowed = (cobID & MapRTRNotAllowed) == 0

	// Get Trans type
	if err := m.ComRecord.FindIndex(2).Read(); err != nil {
		return err
	}

	transType := *m.ComRecord.FindIndex(2).GetByteVal()
	m.TransType = transType

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
	nofEntries := m.MapArray.FindIndex(0).GetData()

	for i := range nofEntries {
		ii := uint16(i + 1)
		if err := m.MapArray.FindIndex(ii).Read(); err != nil {
			return err
		}
		val := *m.MapArray.FindIndex(ii).GetIntVal()

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

func (m *PDOMap) Save() error {
	return nil
}
