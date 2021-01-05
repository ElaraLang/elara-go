package interpreter

type Instance struct {
	Type   Type
	Values map[string]*Value
}
