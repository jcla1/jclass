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
}

type ConstInfoTag uint8
type AccessFlag uint16

type ConstInfo struct {
	Tag  ConstInfoTag
	Info []uint8
}
