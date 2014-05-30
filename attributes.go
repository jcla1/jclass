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

func writeAttributes(w io.Writer, attrs Attributes) error {
	err := binary.Write(w, byteOrder, uint16(len(attrs)))
	if err != nil {
		return err
	}

	for _, attr := range attrs {
		err := attr.Dump(w)
		if err != nil {
			return err
		}
	}

	return nil
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
	name := constPool.GetUTF8(attrBase.NameIndex)

	switch name {
	case "ConstantValue":
		attr = &ConstantValue{baseAttribute: attrBase}
	case "Code":
		attr = &Code{baseAttribute: attrBase}
	// case "StackMapTable":
	//     attr = &StackMapTable{baseAttribute: attrBase}
	case "Exceptions":
		attr = &Exceptions{baseAttribute: attrBase}
	case "InnerClasses":
		attr = &InnerClasses{baseAttribute: attrBase}
	case "EnclosingMethod":
		attr = &EnclosingMethod{baseAttribute: attrBase}
	case "Synthetic":
		attr = &Synthetic{baseAttribute: attrBase}
	case "Signature":
		attr = &Signature{baseAttribute: attrBase}
	case "SourceFile":
		attr = &SourceFile{baseAttribute: attrBase}
	case "SourceDebugExtension":
		attr = &SourceDebugExtension{baseAttribute: attrBase}
	case "LineNumberTable":
		attr = &LineNumberTable{baseAttribute: attrBase}
	case "LocalVariableTable":
		attr = &LocalVariableTable{baseAttribute: attrBase}
	case "LocalVariableTypeTable":
		attr = &LocalVariableTypeTable{baseAttribute: attrBase}
	case "Deprecated":
		attr = &Deprecated{baseAttribute: attrBase}
	// case "RuntimeVisibleAnnotations":
	// 	attr = &RuntimeVisibleAnnotations{baseAttribute: attrBase}
	// case "RuntimeInvisibleAnnotations":
	// 	attr = &RuntimeInvisibleAnnotations{baseAttribute: attrBase}
	// case "RuntimeVisibleParameterAnnotations":
	// 	attr = &RuntimeVisibleParameterAnnotations{baseAttribute: attrBase}
	// case "RuntimeInvisibleParameterAnnotations":
	// 	attr = &RuntimeInvisibleParameterAnnotations{baseAttribute: attrBase}
	// case "AnnotationDefault":
	// 	attr = &AnnotationDefault{baseAttribute: attrBase}
	case "BootstrapMethods":
		attr = &BootstrapMethods{baseAttribute: attrBase}
	default:
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
	NameIndex ConstPoolIndex
	Length    uint32
}

func (_ baseAttribute) UnknownAttr() *UnknownAttr     { panic("jclass: value is not UnknownAttr") }
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

type UnknownAttr struct {
	baseAttribute
	Data []uint8
}

func (a *UnknownAttr) UnknownAttr() *UnknownAttr { return a }
func (_ *UnknownAttr) GetTag() AttributeType     { return UnknownTag }

func (a *UnknownAttr) Read(r io.Reader, _ ConstantPool) error {
	a.Data = make([]uint8, a.Length)
	return binary.Read(r, byteOrder, a.Data)
}

func (a *UnknownAttr) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, a.Data),
	})
}

// field_info, may single
// ACC_STATIC only
type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

func (a *ConstantValue) ConstantValue() *ConstantValue { return a }
func (_ *ConstantValue) GetTag() AttributeType         { return ConstantValueTag }

func (a *ConstantValue) Read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &a.Index)
}

func (a *ConstantValue) Dump(w io.Writer) error { return binary.Write(w, byteOrder, a) }

// method_info, single
// not if native or abstract
type Code struct {
	baseAttribute

	MaxStackSize   uint16
	MaxLocalsCount uint16

	ByteCode        []uint8
	ExceptionsTable []CodeException

	// only LineNumberTable, LocalVariableTable,
	// LocalVariableTypeTable, StackMapTable
	Attributes
}

type CodeException struct {
	StartPC   uint16
	EndPC     uint16
	HandlerPC uint16
	// may be zero, then used for finally
	CatchType ConstPoolIndex
}

func (a *Code) Code() *Code           { return a }
func (_ *Code) GetTag() AttributeType { return CodeTag }

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

	a.ExceptionsTable = make([]CodeException, exceptionsCount)
	err = binary.Read(r, byteOrder, a.Exceptions)
	if err != nil {
		return err
	}

	a.Attributes, err = readAttributes(r, constPool)
	return err
}

func (a *Code) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, a.MaxStackSize),
		binary.Write(w, byteOrder, a.MaxLocalsCount),
		binary.Write(w, byteOrder, uint32(len(a.ByteCode))),
		binary.Write(w, byteOrder, a.ByteCode),
		binary.Write(w, byteOrder, uint16(len(a.ExceptionsTable))),
		binary.Write(w, byteOrder, a.ExceptionsTable),
		writeAttributes(w, a.Attributes),
	})
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
func (_ *Exceptions) GetTag() AttributeType   { return ExceptionsTag }

func (a *Exceptions) Read(r io.Reader, _ ConstantPool) error {
	var exceptionsCount uint16
	err := binary.Read(r, byteOrder, &exceptionsCount)
	if err != nil {
		return err
	}

	a.ExceptionsTable = make([]ConstPoolIndex, exceptionsCount)
	return binary.Read(r, byteOrder, a.ExceptionsTable)
}

func (a *Exceptions) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.ExceptionsTable))),
		binary.Write(w, byteOrder, a.ExceptionsTable),
	})
}

// ClassFile, may single
type InnerClasses struct {
	baseAttribute
	Classes []InnerClass
}

type InnerClass struct {
	InnerClassIndex  ConstPoolIndex
	OuterClassIndex  ConstPoolIndex
	InnerName        ConstPoolIndex
	InnerAccessFlags AccessFlags
}

func (a *InnerClasses) InnerClasses() *InnerClasses { return a }
func (_ *InnerClasses) GetTag() AttributeType       { return InnerClassesTag }

func (a *InnerClasses) Read(r io.Reader, _ ConstantPool) error {
	var classesCount uint16
	err := binary.Read(r, byteOrder, &classesCount)
	if err != nil {
		return err
	}

	a.Classes = make([]InnerClass, classesCount)
	return binary.Read(r, byteOrder, a.Classes)
}

func (a *InnerClasses) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.Classes))),
		binary.Write(w, byteOrder, a.Classes),
	})
}

// ClassFile, may single
// iff local class or anonymous class
type EnclosingMethod struct {
	baseAttribute
	ClassIndex  ConstPoolIndex
	MethodIndex ConstPoolIndex
}

func (a *EnclosingMethod) EnclosingMethod() *EnclosingMethod { return a }
func (_ *EnclosingMethod) GetTag() AttributeType             { return EnclosingMethodTag }

func (a *EnclosingMethod) Read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &a.ClassIndex),
		binary.Read(r, byteOrder, &a.MethodIndex),
	})
}

func (a *EnclosingMethod) Dump(w io.Writer) error { return binary.Write(w, byteOrder, a) }

// ClassFile, method_info or field_info, may single
// if compiler generated
// instead maybe: ACC_SYNTHETIC
type Synthetic struct{ baseAttribute }

func (a *Synthetic) Synthetic() *Synthetic                  { return a }
func (_ *Synthetic) GetTag() AttributeType                  { return SyntheticTag }
func (_ *Synthetic) Read(_ io.Reader, _ ConstantPool) error { return nil }
func (a *Synthetic) Dump(w io.Writer) error                 { return binary.Write(w, byteOrder, a) }

// ClassFile, field_info, or method_info, may single
type Signature struct {
	baseAttribute
	SignatureIndex ConstPoolIndex
}

func (a *Signature) Signature() *Signature { return a }
func (_ *Signature) GetTag() AttributeType { return SignatureTag }

func (a *Signature) Read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &a.SignatureIndex)
}

func (a *Signature) Dump(w io.Writer) error { return binary.Write(w, byteOrder, a) }

// ClassFile, may single
type SourceFile struct {
	baseAttribute
	SourceFileIndex ConstPoolIndex
}

func (a *SourceFile) SourceFile() *SourceFile { return a }
func (_ *SourceFile) GetTag() AttributeType   { return SourceFileTag }

func (a *SourceFile) Read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &a.SourceFileIndex)
}

func (a *SourceFile) Dump(w io.Writer) error { return binary.Write(w, byteOrder, a) }

// ClassFile, may single
type SourceDebugExtension struct {
	baseAttribute
	DebugExtension string
}

func (a *SourceDebugExtension) SourceDebugExtension() *SourceDebugExtension { return a }
func (_ *SourceDebugExtension) GetTag() AttributeType                       { return SourceDebugExtensionTag }

func (a *SourceDebugExtension) Read(r io.Reader, _ ConstantPool) error {
	var err error

	var length uint32
	err = binary.Read(r, byteOrder, &length)
	if err != nil {
		return err
	}

	str := make([]uint8, length)
	err = binary.Read(r, byteOrder, str)
	if err != nil {
		return err
	}

	a.DebugExtension = string(str)

	return nil
}

func (a *SourceDebugExtension) Dump(w io.Writer) error {
	err := binary.Write(w, byteOrder, uint16(len(a.DebugExtension)))
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, []byte(a.DebugExtension))
}

// Code, may multiple
type LineNumberTable struct {
	baseAttribute
	Table []LineNumber
}

type LineNumber struct {
	StartPC    uint16
	LineNumber uint16
}

func (a *LineNumberTable) LineNumberTable() *LineNumberTable { return a }
func (_ *LineNumberTable) GetTag() AttributeType             { return LineNumberTableTag }

func (a *LineNumberTable) Read(r io.Reader, _ ConstantPool) error {
	var linesCount uint16
	err := binary.Read(r, byteOrder, &linesCount)
	if err != nil {
		return err
	}

	a.Table = make([]LineNumber, linesCount)
	return binary.Read(r, byteOrder, a.Table)
}

func (a *LineNumberTable) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.Table))),
		binary.Write(w, byteOrder, a.Table),
	})
}

// Code, may multiple
type LocalVariableTable struct {
	baseAttribute
	Table []LocalVariable
}

type LocalVariable struct {
	StartPC         uint16
	Length          uint16
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	// index into local variable array of current frame
	Index uint16
}

func (a *LocalVariableTable) LocalVariableTable() *LocalVariableTable { return a }
func (_ *LocalVariableTable) GetTag() AttributeType                   { return LocalVariableTableTag }

func (a *LocalVariableTable) Read(r io.Reader, _ ConstantPool) error {
	var varsCount uint16
	err := binary.Read(r, byteOrder, &varsCount)
	if err != nil {
		return err
	}

	a.Table = make([]LocalVariable, varsCount)
	return binary.Read(r, byteOrder, a.Table)
}

func (a *LocalVariableTable) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.Table))),
		binary.Write(w, byteOrder, a.Table),
	})
}

// Code, may multiple
type LocalVariableTypeTable struct {
	baseAttribute
	Table []LocalVariableType
}

type LocalVariableType struct {
	StartPC        uint16
	Length         uint16
	NameIndex      ConstPoolIndex
	SignatureIndex ConstPoolIndex
	// index into local variable array of current frame
	Index uint16
}

func (a *LocalVariableTypeTable) LocalVariableTypeTable() *LocalVariableTypeTable { return a }
func (_ *LocalVariableTypeTable) GetTag() AttributeType                           { return LocalVariableTypeTableTag }

func (a *LocalVariableTypeTable) Read(r io.Reader, _ ConstantPool) error {
	var varsCount uint16
	err := binary.Read(r, byteOrder, &varsCount)
	if err != nil {
		return err
	}

	a.Table = make([]LocalVariableType, varsCount)
	return binary.Read(r, byteOrder, a.Table)
}

func (a *LocalVariableTypeTable) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.Table))),
		binary.Write(w, byteOrder, a.Table),
	})
}

// ClassFile, field_info, or method_info, may single
type Deprecated struct{ baseAttribute }

func (a *Deprecated) Deprecated() *Deprecated                { return a }
func (_ *Deprecated) GetTag() AttributeType                  { return DeprecatedTag }
func (_ *Deprecated) Read(r io.Reader, _ ConstantPool) error { return nil }
func (a *Deprecated) Dump(w io.Writer) error                 { return binary.Write(w, byteOrder, a) }

type RuntimeVisibleAnnotations struct {
	baseAttribute
}

func (a *RuntimeVisibleAnnotations) RuntimeVisibleAnnotations() *RuntimeVisibleAnnotations { return a }

type RuntimeInvisibleAnnotations struct {
	baseAttribute
}

func (a *RuntimeInvisibleAnnotations) RuntimeInvisibleAnnotations() *RuntimeInvisibleAnnotations {
	return a
}

type RuntimeVisibleParameterAnnotations struct {
	baseAttribute
}

func (a *RuntimeVisibleParameterAnnotations) RuntimeVisibleParameterAnnotations() *RuntimeVisibleParameterAnnotations {
	return a
}

type RuntimeInvisibleParameterAnnotations struct {
	baseAttribute
}

func (a *RuntimeInvisibleParameterAnnotations) RuntimeInvisibleParameterAnnotations() *RuntimeInvisibleParameterAnnotations {
	return a
}

type AnnotationDefault struct {
	baseAttribute
}

func (a *AnnotationDefault) AnnotationDefault() *AnnotationDefault { return a }

// ClassFile, may single
// iff constpool conatains CONSTANT_InvokeDynamic_info
type BootstrapMethods struct {
	baseAttribute
	Methods []BootstrapMethod
}

type BootstrapMethod struct {
	MethodRef ConstPoolIndex
	Args      []ConstPoolIndex
}

func (a *BootstrapMethods) BootstrapMethods() *BootstrapMethods { return a }
func (_ *BootstrapMethods) GetTag() AttributeType               { return BootstrapMethodsTag }

func (a *BootstrapMethods) Read(r io.Reader, _ ConstantPool) error {
	var methodsCount uint16
	err := binary.Read(r, byteOrder, &methodsCount)
	if err != nil {
		return err
	}

	a.Methods = make([]BootstrapMethod, 0, methodsCount)

	for i := uint16(0); i < methodsCount; i++ {
		method := BootstrapMethod{}
		err := method.read(r)
		if err != nil {
			return err
		}

		a.Methods = append(a.Methods, method)
	}

	return nil
}

func (a *BootstrapMethods) Dump(w io.Writer) error {
	err := multiError([]error{
		binary.Write(w, byteOrder, a.baseAttribute),
		binary.Write(w, byteOrder, uint16(len(a.Methods))),
	})
	if err != nil {
		return err
	}

	for _, method := range a.Methods {
		err := method.dump(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *BootstrapMethod) read(r io.Reader) error {
	var err error

	err = binary.Read(r, byteOrder, &a.MethodRef)
	if err != nil {
		return err
	}

	var argsCount uint16
	err = binary.Read(r, byteOrder, &a.MethodRef)
	if err != nil {
		return err
	}

	a.Args = make([]ConstPoolIndex, argsCount)
	return binary.Read(r, byteOrder, a.Args)
}

func (a *BootstrapMethod) dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, a.MethodRef),
		binary.Write(w, byteOrder, uint16(len(a.Args))),
		binary.Write(w, byteOrder, a.Args),
	})
}
