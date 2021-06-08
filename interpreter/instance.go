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

func (i *Instance) Equals(ctx *Context, other *Value) bool {
	otherAsInstance, otherIsInstance := other.Value.(*Instance)
	if !otherIsInstance {
		return false
	}
	if !i.Type.Accepts(otherAsInstance.Type, ctx) {
		return false
	}
	for key, value := range i.Values {
		otherVal, present := otherAsInstance.Values[key]
		if !present {
			return false
		}
		if !value.Equals(ctx, otherVal) {
			return false
		}
	}
	return true
}
