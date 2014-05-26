package class

type Attribute interface {
	isAttr()
}

type baseAttribute struct {
	NameIndex ConstPoolIndex
	Length    uint32
}

type ConstantValueAttr struct {
	baseAttribute
	Index ConstPoolIndex
}
