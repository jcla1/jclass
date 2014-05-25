package class

// ClassFile represents a single class file as specified in:
// http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html
type ClassFile struct {
	// Magic number found in all valid Java class files.
	// This will always equal 0xCAFEBABE
	Magic uint32

	// Major.Minor denotes the class file version, that
	// has to be supported by the executing JVM.
	MinorVersion uint16
	MajorVersion uint16

	// The constant pool is a table of structures
	// representing various string constants, class,
	// interface & field names and other constants that
	// are referred to in the class file structure.
	ConstPoolSize uint16
	ConstPool     []*ConstInfo

	// AccessFlags is a mask of flags used to denote
	// access permissions and properties of this class
	// or interface.
	AccessFlags ClassAccessFlag

	// Index into the constant pool, where you should
	// find a CONSTANT_Class_info struct that describes
	// this class.
	ThisClass ConstPoolIndex

	// Index into the constant pool or zero, where you
	// should find a CONSTANT_Class_info struct that
	// describes this class' super class.
	// If SuperClass is zero, then this class must
	// represent the Object class.
	// For an interface, the corresponding value in the
	// constant pool, must represent the Object class.
	SuperClass ConstPoolIndex

	// Interfaces contains indexes into the constant pool,
	// where every referenced entry describes a
	// CONSTANT_Class_info struct representing a direct
	// super-interface of this class or interface.
	InterfacesCount uint16
	Interfaces      []ConstPoolIndex

	// Fields contains indexes into the constant pool,
	// referencing field_info structs, giving a complete
	// description of a field in this class or interface.
	// The Fields table only contains fields declared in
	// this class or interface, not any inherited ones.
	FieldsCount uint16
	Fields      []*FieldInfo

	// Methods contains method_info structs describing
	// a method of this class or interface.
	// If neiter METHOD_ACC_NATIVE or METHOD_ACC_ABSTRACT
	// flags are set, the corresponding code for the method
	// will also be supplied.
	MethodsCount uint16
	Methods      []*MethodInfo

	// Attributes describes properties of this class or
	// interface through attribute_info structs.
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
