package dic

type Object interface {
	GetIndex() int
	GetName() string
	AddMember(Object)
	FindIndex(int) (Object, error)
	FindName(string) (Object, error)
}
