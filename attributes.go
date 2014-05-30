package class

import (
	"encoding/binary"
	"io"
)

func readAttributes(r io.Reader, constPool ConstantPool) (Attributes, error) {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return nil, err
	}

	attrs := make(Attributes, 0, count)

	for i := uint16(0); i < count; i++ {
		attr, err := readAttribute(r, constPool)
		if err != nil {
			return nil, err
		}

		attrs = append(attrs, attr)
	}

	return attrs, nil
}

func readAttribute(r io.Reader, constPool ConstantPool) (Attribute, error) {
	attrBase := baseAttribute{}

	err := multiError([]error{
		binary.Read(r, byteOrder, &attrBase.NameIndex),
		binary.Read(r, byteOrder, &attrBase.Length),
	})

	if err != nil {
		return nil, err
	}

	return fillAttribute(r, attrBase, constPool)
}

func fillAttribute(r io.Reader, attrBase baseAttribute, constPool ConstantPool) (Attribute, error) {
	var attr Attribute
	name := constPool.GetString(attrBase.NameIndex)

	switch name {
	case "ConstantValue":
		attrBase.attrType = ConstantValueTag
		attr = &ConstantValue{baseAttribute: attrBase}
	case "Code":
		attrBase.attrType = CodeTag
		attr = &Code{baseAttribute: attrBase}
	// case "StackMapTable":
	//     attrBase.attrType = StackMapTableTag
	//     attr = &StackMapTable{baseAttribute: attrBase}
	case "Exceptions":
		attrBase.attrType = ExceptionsTag
		attr = &Exceptions{baseAttribute: attrBase}
	case "InnerClasses":
		attrBase.attrType = InnerClassesTag
		attr = &InnerClasses{baseAttribute: attrBase}
	case "EnclosingMethod":
		attrBase.attrType = EnclosingMethodTag
		attr = &EnclosingMethod{baseAttribute: attrBase}
	case "Synthetic":
		attrBase.attrType = SyntheticTag
		attr = &Synthetic{baseAttribute: attrBase}
	case "Signature":
		attrBase.attrType = SignatureTag
		attr = &Signature{baseAttribute: attrBase}
	case "SourceFile":
		attrBase.attrType = SourceFileTag
		attr = &SourceFile{baseAttribute: attrBase}
	case "SourceDebugExtension":
		attrBase.attrType = SourceDebugExtensionTag
		attr = &SourceDebugExtension{baseAttribute: attrBase}
	case "LineNumberTable":
		attrBase.attrType = LineNumberTableTag
		attr = &LineNumberTable{baseAttribute: attrBase}
	case "LocalVariableTable":
		attrBase.attrType = LocalVariableTableTag
		attr = &LocalVariableTable{baseAttribute: attrBase}
	case "LocalVariableTypeTable":
		attrBase.attrType = LocalVariableTypeTableTag
		attr = &LocalVariableTypeTable{baseAttribute: attrBase}
	case "Deprecated":
		attrBase.attrType = DeprecatedTag
		attr = &Deprecated{baseAttribute: attrBase}
	// case "RuntimeVisibleAnnotations":
	// 	attrBase.attrType = RuntimeVisibleAnnotationsTag
	// 	attr = &RuntimeVisibleAnnotations{baseAttribute: attrBase}
	// case "RuntimeInvisibleAnnotations":
	// 	attrBase.attrType = RuntimeInvisibleAnnotationsTag
	// 	attr = &RuntimeInvisibleAnnotations{baseAttribute: attrBase}
	// case "RuntimeVisibleParameterAnnotations":
	// 	attrBase.attrType = RuntimeVisibleParameterAnnotationsTag
	// 	attr = &RuntimeVisibleParameterAnnotations{baseAttribute: attrBase}
	// case "RuntimeInvisibleParameterAnnotations":
	// 	attrBase.attrType = RuntimeInvisibleParameterAnnotationsTag
	// 	attr = &RuntimeInvisibleParameterAnnotations{baseAttribute: attrBase}
	// case "AnnotationDefault":
	// 	attrBase.attrType = AnnotationDefaultTag
	// 	attr = &AnnotationDefault{baseAttribute: attrBase}
	case "BootstrapMethods":
		attrBase.attrType = BootstrapMethodsTag
		attr = &BootstrapMethods{baseAttribute: attrBase}
	default:
		attrBase.attrType = UnknownTag
		attr = &UnknownAttr{baseAttribute: attrBase}
	}

	err := attr.Read(r, constPool)
	if err != nil {
		return nil, err
	}

	return attr, nil
}

type AttributeType uint8

type baseAttribute struct {
	attrType  AttributeType
	NameIndex ConstPoolIndex
	Length    uint16
}

func (b baseAttribute) GetTag() AttributeType {
	return b.attrType
}

func (_ baseAttribute) ConstantValue() *ConstantValue { panic("jclass: value is not ConstantValue") }
func (_ baseAttribute) Code() *Code                   { panic("jclass: value is not Code") }
func (_ baseAttribute) StackMapTable() *StackMapTable { panic("jclass: value is not StackMapTable") }
func (_ baseAttribute) Exceptions() *Exceptions       { panic("jclass: value is not Exceptions") }
func (_ baseAttribute) InnerClasses() *InnerClasses   { panic("jclass: value is not InnerClasses") }
func (_ baseAttribute) EnclosingMethod() *EnclosingMethod {
	panic("jclass: value is not EnclosingMethod")
}
func (_ baseAttribute) Synthetic() *Synthetic   { panic("jclass: value is not Synthetic") }
func (_ baseAttribute) Signature() *Signature   { panic("jclass: value is not Signature") }
func (_ baseAttribute) SourceFile() *SourceFile { panic("jclass: value is not SourceFile") }
func (_ baseAttribute) SourceDebugExtension() *SourceDebugExtension {
	panic("jclass: value is not SourceDebugExtension")
}
func (_ baseAttribute) LineNumberTable() *LineNumberTable {
	panic("jclass: value is not LineNumberTable")
}
func (_ baseAttribute) LocalVariableTable() *LocalVariableTable {
	panic("jclass: value is not LocalVariableTable")
}
func (_ baseAttribute) LocalVariableTypeTable() *LocalVariableTypeTable {
	panic("jclass: value is not LocalVariableTypeTable")
}
func (_ baseAttribute) Deprecated() *Deprecated { panic("jclass: value is not Deprecated") }
func (_ baseAttribute) RuntimeVisibleAnnotations() *RuntimeVisibleAnnotations {
	panic("jclass: value is not RuntimeVisibleAnnotations")
}
func (_ baseAttribute) RuntimeInvisibleAnnotations() *RuntimeInvisibleAnnotations {
	panic("jclass: value is not RuntimeInvisibleAnnotations")
}
func (_ baseAttribute) RuntimeVisibleParameterAnnotations() *RuntimeVisibleParameterAnnotations {
	panic("jclass: value is not RuntimeVisibleParameterAnnotations")
}
func (_ baseAttribute) RuntimeInvisibleParameterAnnotations() *RuntimeInvisibleParameterAnnotations {
	panic("jclass: value is not RuntimeInvisibleParameterAnnotations")
}
func (_ baseAttribute) AnnotationDefault() *AnnotationDefault {
	panic("jclass: value is not AnnotationDefault")
}
func (_ baseAttribute) BootstrapMethods() *BootstrapMethods {
	panic("jclass: value is not BootstrapMethods")
}

// field_info, may single
// ACC_STATIC only
type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

func (a *ConstantValue) ConstantValue() *ConstantValue { return a }

func (a *ConstantValue) Read(r io.Reader, constPool ConstantPool) error {
	return binary.Read(r, byteOrder, &a.Index)
}

func (a *ConstantValue) Dump(w io.Writer) error { return nil }

// method_info, single
// not if native or abstract
type Code struct {
	baseAttribute

	MaxStackSize   uint16
	MaxLocalsCount uint16

	ByteCode        []uint8
	ExceptionsTable []CodeExceptions

	// only LineNumberTable, LocalVariableTable,
	// LocalVariableTypeTable, StackMapTable
	Attributes
}

type CodeExceptions struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	// may be zero, then used for finally
	CatchType ConstPoolIndex
}

func (a *Code) Code() *Code { return a }

func (a *Code) Read(r io.Reader, constPool ConstantPool) error {
	var err error

	var codeLen uint32
	err = multiError([]error{
		binary.Read(r, byteOrder, &a.MaxStackSize),
		binary.Read(r, byteOrder, &a.MaxLocalsCount),
		binary.Read(r, byteOrder, &codeLen),
	})
	if err != nil {
		return err
	}

	a.ByteCode = make([]uint8, codeLen)
	err = binary.Read(r, byteOrder, a.ByteCode)
	if err != nil {
		return err
	}

	var exceptionsCount uint16
	err = binary.Read(r, byteOrder, &exceptionsCount)
	if err != nil {
		return err
	}

	a.Exceptions = make([]CodeException, exceptionsCount)
	err = binary.Read(r, byteOrder, a.Exceptions)
	if err != nil {
		return err
	}

	a.Attributes, err = readAttributes(r, constPool)
	return err
}

type StackMapTable struct {
	baseAttribute
}

func (a *StackMapTable) StackMapTable() *StackMapTable { return a }

// method_info, may single
type Exceptions struct {
	baseAttribute
	ExceptionsTable []ConstPoolIndex
}

func (a *Exceptions) Exceptions() *Exceptions { return a }

func (a *Exceptions) Read(r io.Reader, _ ConstantPool) error {
	var exceptionsCount uint16
	err := binary.Read(r, byteOrder, &exceptionsCount)
	if err != nil {
		return err
	}

	a.ExceptionsTable = make([]ConstPoolIndex, exceptionsCount)
	return binary.Read(r, byteOrder, a.ExceptionsTable)
}

// ClassFile, may single
type InnerClasses struct {
	baseAttribute

	ClassesCount uint16
	Classes      []struct {
		InnerClassIndex  ConstPoolIndex
		OuterClassIndex  ConstPoolIndex
		InnerName        ConstPoolIndex
		InnerAccessFlags AccessFlags
	}
}

func (_ *InnerClasses) InnerClasses() *InnerClasses { return a }

// ClassFile, may single
// iff local class or anonymous class
type EnclosingMethod struct {
	baseAttribute
	ClassIndex  ConstPoolIndex
	MethodIndex ConstPoolIndex
}

func (_ *EnclosingMethod) EnclosingMethod() *EnclosingMethod { return a }

// ClassFile, method_info or field_info, may single
// if compiler generated
// instead maybe: ACC_SYNTHETIC
type Synthetic baseAttribute

func (_ *Synthetic) Synthetic() *Synthetic { return a }

// ClassFile, field_info, or method_info, may single
type Signature struct {
	baseAttribute
	SignatureIndex ConstPoolIndex
}

func (_ *Signature) Signature() *Signature { return a }

// ClassFile, may single
type SourceFile struct {
	baseAttribute
	SourceFileIndex ConstPoolIndex
}

func (_ *SourceFile) SourceFile() *SourceFile { return a }

// ClassFile, may single
type SourceDebugExtension struct {
	baseAttribute
	DebugExtension string
}

func (_ *SourceDebugExtension) SourceDebugExtension() *SourceDebugExtension { return a }

// Code, may multiple
type LineNumberTable struct {
	baseAttribute
	TableLength uint16
	Table       []struct {
		StartPC    uint16
		LineNumber uint16
	}
}

func (_ *LineNumberTable) LineNumberTable() *LineNumberTable { return a }

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

func (_ *LocalVariableTable) LocalVariableTable() *LocalVariableTable { return a }

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

func (_ *LocalVariableTypeTable) LocalVariableTypeTable() *LocalVariableTypeTable { return a }

// ClassFile, field_info, or method_info, may single
type Deprecated baseAttribute

func (_ *Deprecated) Deprecated() *Deprecated { return a }

type RuntimeVisibleAnnotations struct {
	baseAttribute
}

func (_ *RuntimeVisibleAnnotations) RuntimeVisibleAnnotations() *RuntimeVisibleAnnotations { return a }

type RuntimeInvisibleAnnotations struct {
	baseAttribute
}

func (_ *RuntimeInvisibleAnnotations) RuntimeInvisibleAnnotations() *RuntimeInvisibleAnnotations {
	return a
}

type RuntimeVisibleParameterAnnotations struct {
	baseAttribute
}

func (_ *RuntimeVisibleParameterAnnotations) RuntimeVisibleParameterAnnotations() *RuntimeVisibleParameterAnnotations {
	return a
}

type RuntimeInvisibleParameterAnnotations struct {
	baseAttribute
}

func (_ *RuntimeInvisibleParameterAnnotations) RuntimeInvisibleParameterAnnotations() *RuntimeInvisibleParameterAnnotations {
	return a
}

type AnnotationDefault struct {
	baseAttribute
}

func (_ *AnnotationDefault) AnnotationDefault() *AnnotationDefault { return a }

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

func (_ *BootstrapMethods) BootstrapMethods() *BootstrapMethods { return a }
