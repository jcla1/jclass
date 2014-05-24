package class

type ClassFile struct {
	Magic uint32

	MinorVersion uint16
	MajorVersion uint16

	ConstPoolSize uint16
	ConstPool     []*ConstInfo

	AccessFlags ClassAccessFlag
	ThisClass   ConstPoolIndex
	SuperClass  ConstPoolIndex

	InterfacesCount uint16
	Interfaces      []ConstPoolIndex

	FieldsCount uint16
	Fields      []*FieldInfo

	MethodsCount uint16
	Methods      []*MethodInfo

	AttributesCount uint16
	Attributes      []*AttributeInfo
}

type ClassAccessFlag uint16
type NestedClassAccessFlag uint16
type FieldAccessFlag uint16
type MethodAccessFlag uint16

type FieldInfo struct {
	AccessFlags FieldAccessFlag
	fieldOrMethodInfo
}
type MethodInfo struct {
	AccessFlags MethodAccessFlag
	fieldOrMethodInfo
}

type fieldOrMethodInfo struct {
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

type ConstInfoTag uint8
type ConstPoolIndex uint16

type ConstInfo struct {
	Tag  ConstInfoTag
	Info []uint8
}

type AttributeInfo struct {
	NameIndex ConstPoolIndex
	Length    uint32
	Info      []uint8
}
