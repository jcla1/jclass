package class

import (
	"encoding/binary"
	"io"
)

type Attribute interface {
	isAttr()

	// read must parse from r and populate the struct,
	// not including the baseAttribute embedded struct.
	// This method can expect, that baseAttribute will
	// be populated.
	read(io.Reader, ConstantPool) error

	// write must write the whole binary representation
	// of the attribute to w. Including the baseAttribute.
	write(io.Writer) error
}

func readAttribute(r io.Reader, constPool ConstantPool) (Attribute, error) {
	var err error

	attrBase := new(baseAttribute)
	err = attrBase.read(r)
	if err != nil {
		return nil, err
	}

	attr := initAttributeType(attrBase, constPool)
	attr.read(r, constPool)
	if err != nil {
		return nil, err
	}

	return attr, nil
}

func initAttributeType(base *baseAttribute, constPool ConstantPool) Attribute {
	name := string(constPool[base.NameIndex].Info)

	var attr Attribute

	switch name {
	case "ConstantValue":
		attr = &ConstantValue{}
	case "Code":
		attr = &Code{}

	case "Exceptions":
		attr = &Exceptions{}
	case "InnerClasses":
		attr = &InnerClasses{}
	case "EnclosingMethod":
		attr = &EnclosingMethod{}
	case "Synthetic":
		attr = &Synthetic{}
	case "Signature":
		attr = &Signature{}
	case "SourceDebugExtension":
		attr = &SourceDebugExtension{}
	case "LineNumberTable":
		attr = &LineNumberTable{}
	case "LocalVariableTable":
		attr = &LocalVariableTable{}
	case "Deprecated":
		attr = &Deprecated{}
	// case "RuntimeVisibleAnnotations":
	// 	attr = &RuntimeVisibleAnnotations{}
	// case "RuntimeInvisibleAnnotations":
	// 	attr = &RuntimeInvisibleAnnotations{}
	// case "RuntimeVisibleParameterAnnotations":
	// 	attr = &RuntimeVisibleParameterAnnotations{}
	// case "RuntimeInvisibleParameterAnnotations":
	// 	attr = &RuntimeInvisibleParameterAnnotations{}
	case "AnnotationDefault":
		attr = &AnnotationDefault{}
	case "BootstrapMethods":
		attr = &BootstrapMethods{}
	}

	return attr
}

type Attributes []Attribute

func (as Attributes) write(w io.Writer) error {
	err := binary.Write(w, byteOrder, len(as))
	if err != nil {
		return err
	}

	for _, attr := range as {
		err := attr.write(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (as Attributes) read(r io.Reader, constPool ConstantPool) (uint16, error) {
	var count uint16

	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return 0, err
	}

	as = make(Attributes, 0, count)

	for i := uint16(0); i < count; i++ {
		attr, err := readAttribute(r, constPool)
		if err != nil {
			return 0, err
		}

		as = append(as, attr)
	}

	return count, nil
}

type baseAttribute struct {
	NameIndex ConstPoolIndex
	Length    uint32
}

func (ba baseAttribute) read(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &ba.NameIndex),
		binary.Read(r, byteOrder, &ba.Length),
	})
}

func (ba baseAttribute) write(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, ba.NameIndex),
		binary.Write(w, byteOrder, ba.Length),
	})
}

type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

func (cv *ConstantValue) read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &cv.Index)
}

func (cv *ConstantValue) write(w io.Writer) error {
	return multiError([]error{
		cv.baseAttribute.write(w),
		binary.Write(w, byteOrder, cv.Index)},
	)
}

type Code struct {
	baseAttribute

	MaxStackDepth uint16
	// Warning: Here again, caution: long & double take 2 slots!
	// http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.7.3
	MaxLocalVars uint16

	CodeLength uint32
	Code       []uint8

	// This bit is important for try,catch,finally constructs
	ExceptionsLength uint16
	Exceptions       []struct {
		StartPC   uint16
		EndPC     uint16
		HandlerPC uint16
		CatchType ConstPoolIndex
	}

	AttributesCount uint16
	Attributes
}

// TODO: check if &c.Exceptions works
func (c *Code) read(r io.Reader, constPool ConstantPool) error {
	err := multiError([]error{
		binary.Read(r, byteOrder, &c.MaxStackDepth),
		binary.Read(r, byteOrder, &c.MaxLocalVars),
		binary.Read(r, byteOrder, &c.CodeLength),
		binary.Read(r, byteOrder, &c.Code),
		binary.Read(r, byteOrder, &c.ExceptionsLength),
		binary.Read(r, byteOrder, &c.Exceptions),
	})

	if err != nil {
		return err
	}

	count, err := c.Attributes.read(r, constPool)
	if err != nil {
		return err
	}

	c.AttributesCount = count

	return nil
}

func (c *Code) write(w io.Writer) error {
	return multiError([]error{
		c.baseAttribute.write(w),
		binary.Write(w, byteOrder, c.MaxStackDepth),
		binary.Write(w, byteOrder, c.MaxLocalVars),
		binary.Write(w, byteOrder, c.CodeLength),
		binary.Write(w, byteOrder, c.Code),
		binary.Write(w, byteOrder, c.ExceptionsLength),
		binary.Write(w, byteOrder, c.Exceptions),
		c.Attributes.write(w),
	})
}

type Exceptions struct {
	baseAttribute
	ExceptionsCount uint16
	Exceptions      []ConstPoolIndex
}

func (e *Exceptions) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &e.ExceptionsCount),
		binary.Read(r, byteOrder, &e.Exceptions),
	})
}

func (e *Exceptions) write(w io.Writer) error {
	return multiError([]error{
		e.baseAttribute.write(w),
		binary.Write(w, byteOrder, e.ExceptionsCount),
		binary.Write(w, byteOrder, e.Exceptions),
	})
}

type InnerClasses struct {
	baseAttribute
	ClassesCount uint16
	Classes      []struct {
		InnerClassInfo        ConstPoolIndex
		OuterClassInfo        ConstPoolIndex
		InnerNameIndex        ConstPoolIndex
		InnerClassAccessFlags NestedClassAccessFlag
	}
}

// TODO: check if &i.Classes works
func (i *InnerClasses) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &i.ClassesCount),
		binary.Read(r, byteOrder, &i.Classes),
	})
}

func (i *InnerClasses) write(w io.Writer) error {
	return binary.Write(w, byteOrder, i)
}

type EnclosingMethod struct {
	baseAttribute
	ClassIndex  ConstPoolIndex
	MethodIndex ConstPoolIndex
}

func (e *EnclosingMethod) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &e.ClassIndex),
		binary.Read(r, byteOrder, &e.MethodIndex),
	})
}

func (e *EnclosingMethod) write(w io.Writer) error {
	return binary.Write(w, byteOrder, e)
}

type Synthetic struct {
	baseAttribute
}

func (s *Synthetic) read(_ io.Reader, _ ConstantPool) error { return nil }
func (s *Synthetic) write(w io.Writer) error {
	return s.baseAttribute.write(w)
}

type Signature struct {
	baseAttribute
	SignatureIndex ConstPoolIndex
}

func (s *Signature) read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &s.SignatureIndex)
}

func (s *Signature) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, s.SignatureIndex),
	})
}

type SourceFile struct {
	baseAttribute
	SourceFileIndex ConstPoolIndex
}

func (s *SourceFile) read(r io.Reader, _ ConstantPool) error {
	return binary.Read(r, byteOrder, &s.SourceFileIndex)
}

func (s *SourceFile) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, s.SourceFileIndex),
	})
}

type SourceDebugExtension struct {
	baseAttribute
	DebugExtension string
}

func (s *SourceDebugExtension) read(r io.Reader, _ ConstantPool) error {
	str := make([]byte, s.Length)
	err := binary.Read(r, byteOrder, str)
	if err != nil {
		return err
	}

	s.DebugExtension = string(str)

	return nil
}

func (s *SourceDebugExtension) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, []byte(s.DebugExtension)),
	})
}

type LineNumberTable struct {
	baseAttribute
	TableSize uint16
	Table     []struct {
		StartPC    uint16
		LineNumber uint16
	}
}

// TODO: check if &s.Table works
func (s *LineNumberTable) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &s.TableSize),
		binary.Read(r, byteOrder, &s.Table),
	})
}

func (s *LineNumberTable) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, s.TableSize),
		binary.Write(w, byteOrder, s.Table),
	})
}

type LocalVariableTable struct {
	baseAttribute
	TableSize uint16
	Table     []struct {
		StartPC         uint16
		Length          uint16
		NameIndex       ConstPoolIndex
		DescriptorIndex ConstPoolIndex
		// Again: long & double occupy 2 slots
		Index uint16
	}
}

// TODO: check if &s.Table works
func (s *LocalVariableTable) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &s.TableSize),
		binary.Read(r, byteOrder, &s.Table),
	})
}

func (s *LocalVariableTable) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, s.TableSize),
		binary.Write(w, byteOrder, s.Table),
	})
}

type LocalVariableTypeTable struct {
	baseAttribute
	TableSize uint16
	Table     []struct {
		StartPC        uint16
		Length         uint16
		NameIndex      ConstPoolIndex
		SignatureIndex ConstPoolIndex
		// Again: long & double occupy 2 slots
		Index uint16
	}
}

// TODO: check if &s.Table works
func (s *LocalVariableTypeTable) read(r io.Reader, _ ConstantPool) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &s.TableSize),
		binary.Read(r, byteOrder, &s.Table),
	})
}

func (s *LocalVariableTypeTable) write(w io.Writer) error {
	return multiError([]error{
		s.baseAttribute.write(w),
		binary.Write(w, byteOrder, s.TableSize),
		binary.Write(w, byteOrder, s.Table),
	})
}

type Deprecated struct {
	baseAttribute
}

func (d *Deprecated) read(_ io.Reader, _ ConstantPool) error { return nil }
func (d *Deprecated) write(w io.Writer) error {
	return d.baseAttribute.write(w)
}

type ElementValue struct {
	// TODO: implement!
}

// type Annotation struct {
// 	TypeIndex  ConstPoolIndex
// 	PairsCount uint16
// 	Pairs      []struct {
// 		NameIndex ConstPoolIndex
// 		Value     ElementValue
// 	}
// }

// type Annotations []Annotation

// func (a *Annotations) read(r io.Reader) error {

// }

// type annot struct {
// 	AnnotationsCount uint16
// 	Annotations
// }

// func (a *annot) read(r io.Reader) error {
// 	return multiError([]error{
// 		binary.Read(r, byteOrder, &a.AnnotationsCount),
// 		a.Annotations.read(w),
// 	})
// }

// func (a *annot) write(w io.Writer) error {
// 	return multiError([]error{
// 		binary.Write(w, byteOrder, a.AnnotationsCount),
// 		a.Annotations.write(w),
// 	})
// }

// type RuntimeAnnotations struct {
// 	baseAttribute
// 	annot
// }

// func (a *RuntimeAnnotations) read(r io.Reader, _ ConstantPool) error {
// 	return a.annot.read(r)
// }

// func (a *RuntimeAnnotations) write(w io.Writer) error {
// 	return multiError([]error{
// 		a.baseAttribute.write(w),
// 		a.annot.write(w),
// 	})
// }

// type RuntimeVisibleAnnotations RuntimeAnnotations

// func (a *RuntimeVisibleAnnotations) read(r io.Reader, _ ConstantPool) error {
// 	return a.RuntimeAnnotations.read(r)
// }

// func (a *RuntimeVisibleAnnotations) write(r io.Writer) error {
// 	return a.RuntimeAnnotations.write(r)
// }

// type RuntimeInvisibleAnnotations RuntimeAnnotations

// func (a *RuntimeInvisibleAnnotations) read(r io.Reader, _ ConstantPool) error {
// 	return a.RuntimeAnnotations.read(r)
// }

// func (a *RuntimeInvisibleAnnotations) write(r io.Writer) error {
// 	return a.RuntimeAnnotations.write(r)
// }

// type RuntimeParameterAnnotations struct {
// 	ParametersCount      uint8
// 	ParameterAnnotations []annot
// }

// // TODO: check if this works with ParameterAnnotations
// func (a *RuntimeParameterAnnotations) read(r io.Reader, _ ConstantPool) error {
// 	return binary.Read(r, byteOrder, a)
// }

// // TODO: check if this works with ParameterAnnotations
// func (a *RuntimeParameterAnnotations) write(w io.Writer) error {
// 	return binary.Write(w, byteOrder, a)
// }

// type RuntimeVisibleParameterAnnotations RuntimeParameterAnnotations
// type RuntimeInvisibleParameterAnnotations RuntimeParameterAnnotations

type AnnotationDefault struct {
	baseAttribute
	DefaultValue ElementValue
}

func (a *AnnotationDefault) read(r io.Reader, _ ConstantPool) error {
	return a.DefaultValue.read(r)
}

func (a *AnnotationDefault) write(w io.Writer) error {
	return multiError([]error{
		a.baseAttribute.write(w),
		a.DefaultValue.write(w),
	})
}

type BootstrapMethod struct {
	MethodRef      ConstPoolIndex
	ArgumentsCount uint16
	Arguments      []ConstPoolIndex
}

func (b *BootstrapMethod) read(r io.Reader) error {
	err := multiError([]error{
		binary.Read(r, byteOrder, &b.MethodRef),
		binary.Read(r, byteOrder, &b.ArgumentsCount),
	})

	if err != nil {
		return err
	}

	b.Arguments = make([]ConstPoolIndex, b.ArgumentsCount)
	return binary.Read(r, byteOrder, b.Arguments)
}

func (b *BootstrapMethod) write(w io.Writer) error {
	err := multiError([]error{
		binary.Write(w, byteOrder, b.MethodRef),
		binary.Write(w, byteOrder, b.ArgumentsCount),
	})

	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, b.Arguments)
}

type BootstrapMethods struct {
	baseAttribute
	MethodsCount uint16
	Methods      []BootstrapMethod
}

func (b *BootstrapMethods) read(r io.Reader, _ ConstantPool) error {
	err := binary.Read(r, byteOrder, &b.MethodsCount)
	if err != nil {
		return err
	}

	b.Methods = make([]BootstrapMethod, 0, b.MethodsCount)

	for i := uint16(0); i < b.MethodsCount; i++ {
		method := BootstrapMethod{}
		err := method.read(r)
		if err != nil {
			return err
		}

		b.Methods = append(b.Methods, method)
	}

	return nil
}

func (b *BootstrapMethods) write(w io.Writer) error {
	err := binary.Write(w, byteOrder, b.MethodsCount)
	if err != nil {
		return err
	}

	for _, method := range b.Methods {
		err := method.write(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (_ *ConstantValue) isAttr()          {}
func (_ *Code) isAttr()                   {}
func (_ *Exceptions) isAttr()             {}
func (_ *InnerClasses) isAttr()           {}
func (_ *EnclosingMethod) isAttr()        {}
func (_ *Synthetic) isAttr()              {}
func (_ *Signature) isAttr()              {}
func (_ *SourceFile) isAttr()             {}
func (_ *SourceDebugExtension) isAttr()   {}
func (_ *LineNumberTable) isAttr()        {}
func (_ *LocalVariableTable) isAttr()     {}
func (_ *LocalVariableTypeTable) isAttr() {}
func (_ *Deprecated) isAttr()             {}

// func (_ *RuntimeVisibleAnnotations) isAttr()            {}
// func (_ *RuntimeInvisibleAnnotations) isAttr()          {}
// func (_ *RuntimeVisibleParameterAnnotations) isAttr()   {}
// func (_ *RuntimeInvisibleParameterAnnotations) isAttr() {}
func (_ *AnnotationDefault) isAttr() {}
func (_ *BootstrapMethods) isAttr()  {}

func multiError(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
