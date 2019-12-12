package eds

import (
	"fmt"
	"github.com/angelodlfrtr/go-canopen/dic"
	"gopkg.in/ini.v1"
	"regexp"
	"strconv"
	"strings"
)

func Parse(filePath string) (*dic.ObjectDic, error) {
	// Load ini file
	iniData, err := ini.Load(filePath)
	if err != nil {
		return nil, nil
	}

	// Create object dictionary
	ddic := dic.NewObjectDic()

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

	// Iterate over sections
	for _, sec := range iniData.Sections() {
		sectionName := sec.Name()

		// Match index
		if ok, _ := regexp.MatchString(`^[0-9A-Fa-f]{4}$`, sectionName); ok {
			var index int
			if i, err := strconv.ParseInt(sectionName, 0, 0); err != nil {
				index = int(i)
			}

			name := sec.Key("ParameterName").String()
			objectType, _ := strconv.ParseInt(sec.Key("ObjectType").String(), 0, 0)

			// Object type == VARIABLE
			if byte(objectType) == dic.VAR {
				variable := buildVariable(index, 0, name, sec, iniData)
				ddic.AddObject(variable)
			}

			// Object type == ARRAY
			if byte(objectType) == dic.ARR {
				array := &dic.Array{Index: index, Name: name}
				ddic.AddObject(array)
			}

			// Object type == RECORD
			if byte(objectType) == dic.RECORD {
				record := &dic.Record{Index: index, Name: name}
				ddic.AddObject(record)
			}

			continue
		}

		// Match sub-indexs
		if ok, _ := regexp.MatchString(`^([0-9A-Fa-f]{4})sub([0-9A-Fa-f]+)$`, sectionName); ok {
			var index, subIndex int
			if i, err := strconv.ParseInt(sectionName[0:4], 16, 0); err == nil {
				index = int(i)
			}

			if i, err := strconv.ParseInt(sectionName[7:], 16, 0); err == nil {
				subIndex = int(i)
			}

			name := sec.Key("ParameterName").String()

			if object, err := ddic.FindIndex(index); err == nil {
				variable := buildVariable(index, subIndex, name, sec, iniData)
				object.AddMember(variable)
			}

			continue
		}
	}

	return ddic, nil
}

func buildVariable(index, subIndex int, name string, sec *ini.Section, iniData *ini.File) *dic.Variable {
	variable := &dic.Variable{
		Index:      index,
		SubIndex:   subIndex,
		Name:       name,
		AccessType: strings.ToLower(sec.Key("AccessType").String()),
	}

	// Set DataType
	if i, err := strconv.ParseInt(sec.Key("DataType").String(), 0, 0); err != nil {
		variable.DataType = byte(i)
	}

	if variable.DataType > 0x1B {
		dTypeStr := fmt.Sprintf("%dsub1", variable.DataType)

		if i, err := strconv.ParseInt(iniData.Section(dTypeStr).Key("DefaultValue").String(), 0, 0); err != nil {
			variable.DataType = byte(i)
		}
	}

	if lowl, err := sec.GetKey("LowLimit"); err == nil {
		if i, err := strconv.ParseInt(lowl.String(), 0, 0); err != nil {
			variable.Min = int(i)
		}
	}

	if howl, err := sec.GetKey("HighLimit"); err == nil {
		if i, err := strconv.ParseInt(howl.String(), 0, 0); err != nil {
			variable.Max = int(i)
		}
	}

	if def, err := sec.GetKey("DefaultValue"); err == nil {
		if i, err := strconv.ParseInt(def.String(), 0, 64); err != nil {
			variable.Default = i
		}
	}

	return variable
}
