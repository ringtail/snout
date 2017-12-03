package types

type Manager interface {
	Empty() bool
	Start(name ...string)
}
