package class

type ClassFile struct {
	Magic uint32

	MinorVersion uint16
	MajorVersion uint16

	ConstPoolSize uint16
	ConstPool     []*ConstInfo

	AccessFlags uint16
	ThisClass   uint16
	SuperClass  uint16

	InterfacesCount uint16
	Interfaces      []uint16

	FieldsCount uint16
	Fields      []*FieldInfo
}

type ConstInfoTag uint8
type AccessFlag uint16

type ConstInfo struct {
	Tag  ConstInfoTag
	Info []uint8
}

type FieldInfo struct {
	AccessFlags     AccessFlag
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

type AttributeInfo struct {
	NameIndex uint16
	Length    uint32
	Info      []uint8
}
