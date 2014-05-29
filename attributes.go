package class

type AttributeType uint8

type baseAttribute struct {
	attrType  AttributeType
	NameIndex ConstPoolIndex
	Length    uint16
}

// field_info, may single
// ACC_STATIC only
type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

// method_info, single
// not if native or abstract
type Code struct {
	baseAttribute

	MaxStackSize   uint16
	MaxLocalsCount uint16

	CodeLength uint32
	Code       []uint8

	ExceptionsCount uint16
	Exceptions      []struct {
		StartPC   uint16
		EndPC     uint16
		HandlerPC uint16
		// may be zero, then used for finally
		CatchType ConstPoolIndex
	}

	// only LineNumberTable, LocalVariableTable,
	// LocalVariableTypeTable, StackMapTable
	AttributesCount uint16
	Attributes
}

type StackMapTable struct {
	baseAttribute
}

// method_info, may single
type Exceptions struct {
	baseAttribute
	ExceptionsCount uint16
	Exceptions      []ConstPoolIndex
}

// ClassFile, may single
type InnerClasses struct {
	baseAttribute

	ClassesCount uint16
	Classes      []struct {
		InnerClassIndex  ConstPoolIndex
		OuterClassIndex  ConstPoolIndex
		InnerName        ConstPoolIndex
		InnerAccessFlags NestedClassAccessFlag
	}
}

// ClassFile, may single
// iff local class or anonymous class
type EnclosingMethod struct {
	baseAttribute
	ClassIndex  ConstPoolIndex
	MethodIndex ConstPoolIndex
}

// ClassFile, method_info or field_info, may single
// if compiler generated
// instead maybe: ACC_SYNTHETIC
type Synthetic baseAttribute

// ClassFile, field_info, or method_info, may single
type Signature struct {
	baseAttribute
	SignatureIndex ConstPoolIndex
}

// ClassFile, may single
type SourceFile struct {
	baseAttribute
	SourceFileIndex ConstPoolIndex
}

// ClassFile, may single
type SourceDebugExtension struct {
	baseAttribute
	DebugExtension string
}

// Code, may multiple
type LineNumberTable struct {
	baseAttribute
	TableLength uint16
	Table       []struct {
		StartPC    uint16
		LineNumber uint16
	}
}

// Code, may multiple
type LocalVariableTable struct {
	baseAttribute
	TableLength uint16
	Table       []struct {
		StartPC         uint16
		Length          uint16
		NameIndex       ConstPoolIndex
		DescriptorIndex ConstPoolIndex
		// index into local variable array of current frame
		Index uint16
	}
}

// Code, may multiple
type LocalVariableTypeTable struct {
	baseAttribute
	TableLength uint16
	Table       []struct {
		StartPC        uint16
		Length         uint16
		NameIndex      ConstPoolIndex
		SignatureIndex ConstPoolIndex
		// index into local variable array of current frame
		Index uint16
	}
}

// ClassFile, field_info, or method_info, may single
type Deprecated baseAttribute

type RuntimeVisibleAnnotations struct {
	baseAttribute
}

type RuntimeInvisibleAnnotations struct {
	baseAttribute
}

type RuntimeVisibleParameterAnnotations struct {
	baseAttribute
}

type RuntimeInvisibleParameterAnnotations struct {
	baseAttribute
}

type AnnotationDefault struct {
	baseAttribute
}

// ClassFile, may single
// iff constpool conatains CONSTANT_InvokeDynamic_info
type BootstrapMethods struct {
	baseAttribute
	MethodsCount uint16
	Methods      []struct {
		MethodRef ConstPoolIndex
		ArgsCount uint16
		Args      []ConstPoolIndex
	}
}
