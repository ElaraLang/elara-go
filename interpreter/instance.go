package interpreter

type Instance struct {
	Type   *StructType
	Values map[string]*Value
}

func (i *Instance) String() string {
	base := i.Type.Name() + " {"
	for _, v := range i.Type.Properties {
		value := i.Values[v.Name]
		base += "\n    " + v.Name + ": "
		base += value.String()
	}
	base += "\n}"
	return base
}
