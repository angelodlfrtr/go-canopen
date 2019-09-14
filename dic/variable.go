package dic

type Variable struct {
	Unit        string
	Factor      int
	Min         int
	Max         int
	Default     int64
	DataType    byte
	AccessType  string
	Description string

	Data   []byte
	Offset int

	Index    int
	SubIndex int
	Name     string

	ValueDescriptions map[string]string
	BitDefinitions    map[string][]byte
}

func (variable *Variable) GetIndex() int {
	return variable.Index
}

func (variable *Variable) GetName() string {
	return variable.Name
}

func (variable *Variable) AddMember(object Object) {
	// Do nothing
}

func (variable *Variable) FindIndex(index int) (Object, error) {
	// Do nothing
	return nil, nil
}

func (variable *Variable) FindName(name string) (Object, error) {
	// Do nothing
	return nil, nil
}

func (variable *Variable) AddValueDescription(name string, des string) {
	variable.ValueDescriptions[name] = des
}

func (variable *Variable) AddBitDefinition(name string, bits []byte) {
	variable.BitDefinitions[name] = bits
}

func (variable *Variable) Read() ([]byte, error) {
	// @TODO
	return nil, nil
}

func (variable *Variable) Write([]byte) error {
	// @TODO
	return nil
}

func (variable *Variable) SetRaw(val []byte) {
	variable.Data = val
}
