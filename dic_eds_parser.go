package canopen

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// Parse EDS File
// If in is string, it must be a path to a file
// else if in must be eds data as []byte
func DicEDSParse(in interface{}) (*DicObjectDic, error) {
	// Load ini file
	iniData, err := ini.Load(in)
	if err != nil {
		return nil, nil
	}

	// Create object dictionary
	ddic := NewDicObjectDic()

	// Get NodeID & Baudrate
	if sec, err := iniData.GetSection("DeviceComissioning"); err == nil {
		// Get NodeID
		if key, err := sec.GetKey("NodeId"); err == nil {
			ab, _ := strconv.ParseInt(key.String(), 0, 0)
			ddic.NodeID = int(ab)
		}

		// Get Baudrate
		if key, err := sec.GetKey("Baudrate"); err == nil {
			ab, _ := strconv.ParseInt(key.String(), 0, 0)
			ddic.Baudrate = int(ab)
		}
	}

	matchIdxRegexp := regexp.MustCompile(`^[0-9A-Fa-f]{4}$`)
	matchSubIdxRegexp := regexp.MustCompile(`^([0-9A-Fa-f]{4})sub([0-9A-Fa-f]+)$`)

	// Iterate over sections
	for _, sec := range iniData.Sections() {
		sectionName := sec.Name()

		// Match index
		if matchIdxRegexp.MatchString(sectionName) {
			idx, err := strconv.ParseUint(sectionName, 16, 16)
			if err != nil {
				return nil, err
			}

			index := uint16(idx)

			name := sec.Key("ParameterName").String()
			objectType, _ := strconv.ParseUint(sec.Key("ObjectType").String(), 0, 8)

			// Object type == VARIABLE
			if byte(objectType) == DicVar {
				variable, err := buildVariable(index, 0, name, sec, iniData)
				if err != nil {
					return nil, err
				}
				ddic.AddObject(variable)
			}

			// Object type == ARRAY
			if byte(objectType) == DicArr {
				array := &DicArray{Index: index, Name: name}
				ddic.AddObject(array)
			}

			// Object type == RECORD
			if byte(objectType) == DicRec {
				record := &DicRecord{Index: index, Name: name}
				ddic.AddObject(record)
			}

			continue
		}

		// Match sub-indexs
		if matchSubIdxRegexp.MatchString(sectionName) {
			idx, err := strconv.ParseUint(sectionName[0:4], 16, 16)
			if err != nil {
				return nil, err
			}

			index := uint16(idx)
			sidx, err := strconv.ParseUint(sectionName[7:], 16, 8)
			if err != nil {
				return nil, err
			}

			subIndex := uint8(sidx)
			name := sec.Key("ParameterName").String()

			object := ddic.FindIndex(index)
			if object == nil {
				return nil, fmt.Errorf("index with id %d not found", index)
			}

			variable, err := buildVariable(index, subIndex, name, sec, iniData)
			if err != nil {
				return nil, err
			}
			object.AddMember(variable)

			continue
		}
	}

	return ddic, nil
}

// @TODO: check working
func buildVariable(
	index uint16,
	subIndex uint8,
	name string,
	sec *ini.Section,
	iniData *ini.File,
) (*DicVariable, error) {
	variable := &DicVariable{
		Index:      index,
		SubIndex:   subIndex,
		Name:       name,
		AccessType: strings.ToLower(sec.Key("AccessType").String()),
	}

	// Get & set DataType
	i, err := sec.Key("DataType").Int()
	if err != nil {
		return nil, err
	}
	variable.DataType = byte(i)

	if variable.DataType > 0x1B {
		dTypeStr := fmt.Sprintf("%dsub1", variable.DataType)

		i, err := iniData.Section(dTypeStr).Key("DefaultValue").Uint()
		if err != nil {
			return nil, err
		}

		variable.DataType = byte(i)
	}

	if lowl, err := sec.GetKey("LowLimit"); err == nil {
		i, err := lowl.Int()
		if err != nil {
			return nil, err
		}
		variable.Min = i
	}

	if howl, err := sec.GetKey("HighLimit"); err == nil {
		i, err := howl.Int()
		if err != nil {
			return nil, err
		}
		variable.Max = i
	}

	if def, err := sec.GetKey("DefaultValue"); err == nil {
		variable.Default = []byte(def.Value())
	}

	return variable, nil
}
