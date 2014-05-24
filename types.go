package class

const (
	ConstUTF8               ConstInfoTag = 1
	ConstInteger                         = 3
	ConstFloat                           = 4
	ConstLong                            = 5
	ConstDouble                          = 6
	ConstClass                           = 7
	ConstString                          = 8
	ConstFieldRef                        = 9
	ConstMethodRef                       = 10
	ConstInterfaceMehtodRef              = 11
	ConstNameAndType                     = 12
	ConstMethodHandle                    = 15
	ConstMethodType                      = 16
	ConstInvokeDynamic                   = 18
)

const (
	AccPublic     AccessFlag = 0x0001
	AccFinal                 = 0x0010
	AccSuper                 = 0x0020
	AccInterface             = 0x0200
	AccAbstract              = 0x0400
	AccSynthetic             = 0x1000
	AccAnnotation            = 0x2000
	AccEnum                  = 0x4000
)

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

	AttributesCount uint16
	Attributes      []*AttributeInfo
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
