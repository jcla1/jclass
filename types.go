package class

type ClassFile struct {
	Magic uint32

	MinorVersion uint16
	MajorVersion uint16

	ConstPoolSize uint16
	ConstPool     []*ConstInfo
}

type ConstInfoTag uint8

type ConstInfo struct {
	Tag  ConstInfoTag
	Info []uint8
}
