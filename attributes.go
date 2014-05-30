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
		InnerAccessFlags AccessFlags
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
