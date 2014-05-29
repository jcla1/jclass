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

	// only LineNumberTable, LocalVariableTable, LocalVariableTypeTable
	AttributesCount uint16
	Attributes
}

type StackMapTable struct {
	baseAttribute
}

type Exceptions struct {
	baseAttribute
}

type InnerClasses struct {
	baseAttribute
}

type EnclosingMethod struct {
	baseAttribute
}

type Synthetic struct {
	baseAttribute
}

type Signature struct {
	baseAttribute
}

type SourceFile struct {
	baseAttribute
}

type SourceDebugExtension struct {
	baseAttribute
}

type LineNumberTable struct {
	baseAttribute
}

type LocalVariableTable struct {
	baseAttribute
}

type LocalVariableTypeTable struct {
	baseAttribute
}

type Deprecated struct {
	baseAttribute
}

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

type BootstrapMethods struct {
	baseAttribute
}
