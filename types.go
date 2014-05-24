package class

type ClassFile struct {
	Magic uint32

	MinorVersion uint16
	MajorVersion uint16

	ConstPoolSize uint16
	ConstPool     []*ConstInfo

	AccessFlags AccessFlag
	ThisClass   ConstPoolIndex
	SuperClass  ConstPoolIndex

	InterfacesCount uint16
	Interfaces      []ConstPoolIndex

	FieldsCount uint16
	Fields      []*FieldInfo

	MethodsCount uint16
	Methods      []*MethodInfo
}

type ConstInfoTag uint8
type ConstPoolIndex uint16
type AccessFlag uint16

type FieldInfo fieldOrMethodInfo
type MethodInfo fieldOrMethodInfo

type ConstInfo struct {
	Tag  ConstInfoTag
	Info []uint8
}

type fieldOrMethodInfo struct {
	AccessFlags     AccessFlag
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

type AttributeInfo struct {
	NameIndex ConstPoolIndex
	Length    uint32
	Info      []uint8
}
