package class

type Attribute interface {
	isAttr()
}

type baseAttribute struct {
	NameIndex ConstPoolIndex
	Length    uint32
}

type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

type Code struct {
	baseAttribute

	MaxStackDepth uint16
	// Warning: Here again, caution: long & double take 2 slots!
	// http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.7.3
	MaxLocalVars uint16

	CodeLength uint32
	Code       []uint8

	// This bit is important for try,catch,finally constructs
	ExceptionsLength uint16
	Exceptions       []struct {
		StartPC   uint16
		EndPC     uint16
		HandlerPC uint16
		CatchType ConstPoolIndex
	}

	AttributesCount uint16
	Attributes      []Attribute
}
