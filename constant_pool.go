package class

import (
	"encoding/binary"
	"io"
)

type ConstantType uint8

type baseConstant struct {
	tag uint8
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

type FieldRef fieldMethodInterfaceRef

func (c *FieldRef) Field() *FieldRef { return c }

type MethodRef fieldMethodInterfaceRef

func (c *MethodRef) InterfaceMethod() *MethodRef { return c }

type InterfaceMethodRef fieldMethodInterfaceRef

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
