package class

import (
	"encoding/binary"
	"io"
)

func (constPool ConstantPool) GetUTF8(index ConstPoolIndex) string {
	return constPool[index-1].UTF8().Value
}

func (c *ClassFile) writeConstPool(w io.Writer) error {
	err := binary.Write(w, byteOrder, c.ConstPoolSize)
	if err != nil {
		return err
	}

	for _, constant := range c.ConstantPool {
		err := constant.Dump(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClassFile) readConstPool(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.ConstPoolSize)
	if err != nil {
		return err
	}

	c.ConstantPool = make(ConstantPool, 0, c.ConstPoolSize)

	for i := uint16(1); i < c.ConstPoolSize; i++ {
		constant, err := readConstant(r)
		if err != nil {
			return err
		}

		// This is one of the WORST! bugs ever!
		// They even admit it in the JVM spec:
		//
		// "In retrospect, making 8-byte constants
		// take two constant pool entries was a poor choice."
		//
		// The problem is, that both longs & doubles
		// take up TWO!! slots in the const pool (an it
		// looks like they take 2 in the local variable
		// pool too). That's why we need to advance an
		// extra slot, iff we encounter one of them
		if constant.GetTag() == CONSTANT_Long || constant.GetTag() == CONSTANT_Double {
			i++
		}

		c.ConstantPool = append(c.ConstantPool, constant)
	}

	return nil
}

func readConstant(r io.Reader) (Constant, error) {
	constBase := baseConstant{}

	err := binary.Read(r, byteOrder, &constBase.tag)
	if err != nil {
		return nil, err
	}

	return fillConstant(r, constBase)
}

func fillConstant(r io.Reader, constBase baseConstant) (Constant, error) {
	var constant Constant

	switch constBase.GetTag() {
	case CONSTANT_Class:
		constant = &ClassRef{baseConstant: constBase}
	case CONSTANT_FieldRef:
		constant = &FieldRef{fieldMethodInterfaceRef{baseConstant: constBase}}
	case CONSTANT_MethodRef:
		constant = &MethodRef{fieldMethodInterfaceRef{baseConstant: constBase}}
	case CONSTANT_InterfaceMethodRef:
		constant = &InterfaceMethodRef{fieldMethodInterfaceRef{baseConstant: constBase}}
	case CONSTANT_String:
		constant = &StringRef{baseConstant: constBase}
	case CONSTANT_Integer:
		constant = &IntegerRef{baseConstant: constBase}
	case CONSTANT_Float:
		constant = &FloatRef{baseConstant: constBase}
	case CONSTANT_Long:
		constant = &LongRef{baseConstant: constBase}
	case CONSTANT_Double:
		constant = &DoubleRef{baseConstant: constBase}
	case CONSTANT_NameAndType:
		constant = &NameAndTypeRef{baseConstant: constBase}
	case CONSTANT_UTF8:
		constant = &UTF8Ref{baseConstant: constBase}
	case CONSTANT_MethodHandle:
		constant = &MethodHandleRef{baseConstant: constBase}
	case CONSTANT_MethodType:
		constant = &MethodTypeRef{baseConstant: constBase}
	case CONSTANT_InvokeDynamic:
		constant = &InvokeDynamicRef{baseConstant: constBase}
	default:
		panic("jclass: unknown constant pool tag")
	}

	err := constant.Read(r)
	if err != nil {
		return nil, err
	}

	return constant, nil
}

type ConstantType uint8

type baseConstant struct {
	tag ConstantType
}

func (b baseConstant) GetTag() ConstantType {
	return b.tag
}

func (_ baseConstant) Class() *ClassRef   { panic("jclass: constant is not Class") }
func (_ baseConstant) Field() *FieldRef   { panic("jclass: constant is not Field") }
func (_ baseConstant) Method() *MethodRef { panic("jclass: constant is not Method") }
func (_ baseConstant) InterfaceMethod() *InterfaceMethodRef {
	panic("jclass: constant is not InterfaceMethod")
}
func (_ baseConstant) StringRef() *StringRef          { panic("jclass: constant is not StringRef") }
func (_ baseConstant) Integer() *IntegerRef           { panic("jclass: constant is not Integer") }
func (_ baseConstant) Float() *FloatRef               { panic("jclass: constant is not Float") }
func (_ baseConstant) Long() *LongRef                 { panic("jclass: constant is not Long") }
func (_ baseConstant) Double() *DoubleRef             { panic("jclass: constant is not Double") }
func (_ baseConstant) NameAndType() *NameAndTypeRef   { panic("jclass: constant is not NameAndType") }
func (_ baseConstant) UTF8() *UTF8Ref                 { panic("jclass: constant is not UTF8") }
func (_ baseConstant) MethodHandle() *MethodHandleRef { panic("jclass: constant is not MethodHandle") }
func (_ baseConstant) MethodType() *MethodTypeRef     { panic("jclass: constant is not MethodType") }
func (_ baseConstant) InvokeDynamic() *InvokeDynamicRef {
	panic("jclass: constant is not InvokeDynamic")
}

type ClassRef struct {
	baseConstant
	NameIndex ConstPoolIndex
}

func (c *ClassRef) Class() *ClassRef { return c }

func (c *ClassRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.NameIndex)
}

type fieldMethodInterfaceRef struct {
	baseConstant
	ClassIndex       ConstPoolIndex
	NameAndTypeIndex ConstPoolIndex
}

func (c *fieldMethodInterfaceRef) Read(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &c.ClassIndex),
		binary.Read(r, byteOrder, &c.NameAndTypeIndex),
	})
}

type FieldRef struct {
	fieldMethodInterfaceRef
}

func (c *FieldRef) Field() *FieldRef { return c }

type MethodRef struct {
	fieldMethodInterfaceRef
}

func (c *MethodRef) Method() *MethodRef { return c }

type InterfaceMethodRef struct {
	fieldMethodInterfaceRef
}

func (c *InterfaceMethodRef) InterfaceMethod() *InterfaceMethodRef { return c }

type StringRef struct {
	baseConstant
	Index ConstPoolIndex
}

func (c *StringRef) StringRef() *StringRef { return c }

func (c *StringRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Index)
}

type IntegerRef struct {
	baseConstant
	Value int32
}

func (c *IntegerRef) Integer() *IntegerRef { return c }

func (c *IntegerRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

type FloatRef struct {
	baseConstant
	Value float32
}

func (c *FloatRef) Float() *FloatRef { return c }

func (c *FloatRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

type LongRef struct {
	baseConstant
	Value int64
}

func (c *LongRef) Long() *LongRef { return c }

func (c *LongRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

type DoubleRef struct {
	baseConstant
	Value float64
}

func (c *DoubleRef) Double() *DoubleRef { return c }

func (c *DoubleRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

type NameAndTypeRef struct {
	baseConstant
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
}

func (c *NameAndTypeRef) NameAndType() *NameAndTypeRef { return c }

func (c *NameAndTypeRef) Read(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &c.NameIndex),
		binary.Read(r, byteOrder, &c.DescriptorIndex),
	})
}

type UTF8Ref struct {
	baseConstant
	Value string
}

func (c *UTF8Ref) UTF8() *UTF8Ref { return c }

// TODO: check if the string is handled correctly
func (c *UTF8Ref) Read(r io.Reader) error {
	var err error

	var length uint16
	err = binary.Read(r, byteOrder, &length)
	if err != nil {
		return err
	}

	str := make([]uint8, length)
	err = binary.Read(r, byteOrder, str)
	if err != nil {
		return err
	}

	c.Value = string(str)

	return nil
}

type MethodHandleRef struct {
	baseConstant
	ReferenceKind  uint8
	ReferenceIndex ConstPoolIndex
}

func (c *MethodHandleRef) MethodHandle() *MethodHandleRef { return c }

func (c *MethodHandleRef) Read(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &c.ReferenceKind),
		binary.Read(r, byteOrder, &c.ReferenceIndex),
	})
}

type MethodTypeRef struct {
	baseConstant
	DescriptorIndex ConstPoolIndex
}

func (c *MethodTypeRef) MethodType() *MethodTypeRef { return c }

func (c *MethodTypeRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.DescriptorIndex)
}

type InvokeDynamicRef struct {
	baseConstant
	BootstrapMethodAttrIndex ConstPoolIndex
	NameAndTypeIndex         ConstPoolIndex
}

func (c *InvokeDynamicRef) InvokeDynamic() *InvokeDynamicRef { return c }

func (c *InvokeDynamicRef) Read(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &c.BootstrapMethodAttrIndex),
		binary.Read(r, byteOrder, &c.NameAndTypeIndex),
	})
}
