package class

import (
	"encoding/binary"
	// "fmt"
	"io"
)

func (constPool ConstantPool) GetUTF8(index ConstPoolIndex) string {
	return constPool[index-1].UTF8().Value
}

func (constPool ConstantPool) GetInteger(index ConstPoolIndex) int32 {
	return constPool[index-1].Integer().Value
}

func (constPool ConstantPool) GetFloat(index ConstPoolIndex) float32 {
	return constPool[index-1].Float().Value
}

func (constPool ConstantPool) GetLong(index ConstPoolIndex) int64 {
	return constPool[index-1].Long().Value
}

func (constPool ConstantPool) GetDouble(index ConstPoolIndex) float64 {
	return constPool[index-1].Double().Value
}

func (constPool ConstantPool) GetClass(index ConstPoolIndex) *ClassRef {
	return constPool[index-1].Class()
}

func (constPool ConstantPool) GetString(index ConstPoolIndex) *StringRef {
	return constPool[index-1].StringRef()
}

func (constPool ConstantPool) GetField(index ConstPoolIndex) *FieldRef {
	return constPool[index-1].Field()
}

func (constPool ConstantPool) GetMethod(index ConstPoolIndex) *MethodRef {
	return constPool[index-1].Method()
}

func (constPool ConstantPool) GetInterfaceMethod(index ConstPoolIndex) *InterfaceMethodRef {
	return constPool[index-1].InterfaceMethod()
}

func (constPool ConstantPool) GetNameAndType(index ConstPoolIndex) *NameAndTypeRef {
	return constPool[index-1].NameAndType()
}

func (constPool ConstantPool) GetMethodHandle(index ConstPoolIndex) *MethodHandleRef {
	return constPool[index-1].MethodHandle()
}

func (constPool ConstantPool) GetMethodType(index ConstPoolIndex) *MethodTypeRef {
	return constPool[index-1].MethodType()
}

func (constPool ConstantPool) GetInvokeDynamic(index ConstPoolIndex) *InvokeDynamicRef {
	return constPool[index-1].InvokeDynamic()
}

func (c *ClassFile) writeConstPool(w io.Writer) error {
	err := binary.Write(w, byteOrder, c.ConstPoolSize)
	if err != nil {
		return err
	}

	for i := uint16(0); i < c.ConstPoolSize-1; i++ {
		constant := c.ConstantPool[i]

		// In place because of the most annoying spec ever!
		// For more info, see: readConstPool
		if constant == nil {
			continue
		}

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

	c.ConstantPool = make(ConstantPool, c.ConstPoolSize)

	for i := uint16(1); i < c.ConstPoolSize; i++ {
		constant, err := readConstant(r)
		if err != nil {
			return err
		}

		c.ConstantPool[i-1] = constant

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
	}

	return nil
}

func readConstant(r io.Reader) (Constant, error) {
	constBase := baseConstant{}

	err := binary.Read(r, byteOrder, &constBase.Tag)
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
	Tag ConstantType
}

func (b baseConstant) GetTag() ConstantType {
	return b.Tag
}

func (b baseConstant) Class() *ClassRef   { panic("jclass: constant is not Class") }
func (b baseConstant) Field() *FieldRef   { panic("jclass: constant is not Field") }
func (b baseConstant) Method() *MethodRef { panic("jclass: constant is not Method") }
func (b baseConstant) InterfaceMethod() *InterfaceMethodRef {
	panic("jclass: constant is not InterfaceMethod")
}
func (b baseConstant) StringRef() *StringRef          { panic("jclass: constant is not StringRef") }
func (b baseConstant) Integer() *IntegerRef           { panic("jclass: constant is not Integer") }
func (b baseConstant) Float() *FloatRef               { panic("jclass: constant is not Float") }
func (b baseConstant) Long() *LongRef                 { panic("jclass: constant is not Long") }
func (b baseConstant) Double() *DoubleRef             { panic("jclass: constant is not Double") }
func (b baseConstant) NameAndType() *NameAndTypeRef   { panic("jclass: constant is not NameAndType") }
func (b baseConstant) UTF8() *UTF8Ref                 { panic("jclass: constant is not UTF8") }
func (b baseConstant) MethodHandle() *MethodHandleRef { panic("jclass: constant is not MethodHandle") }
func (b baseConstant) MethodType() *MethodTypeRef     { panic("jclass: constant is not MethodType") }
func (b baseConstant) InvokeDynamic() *InvokeDynamicRef {
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

func (c *ClassRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
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

func (c *fieldMethodInterfaceRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
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

func (c *StringRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

type IntegerRef struct {
	baseConstant
	Value int32
}

func (c *IntegerRef) Integer() *IntegerRef { return c }

func (c *IntegerRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

func (c *IntegerRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

type FloatRef struct {
	baseConstant
	Value float32
}

func (c *FloatRef) Float() *FloatRef { return c }

func (c *FloatRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

func (c *FloatRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

type LongRef struct {
	baseConstant
	Value int64
}

func (c *LongRef) Long() *LongRef { return c }

func (c *LongRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

func (c *LongRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

type DoubleRef struct {
	baseConstant
	Value float64
}

func (c *DoubleRef) Double() *DoubleRef { return c }

func (c *DoubleRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Value)
}

func (c *DoubleRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
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

func (c *NameAndTypeRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
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

func (c *UTF8Ref) Dump(w io.Writer) error {
	err := multiError([]error{
		binary.Write(w, byteOrder, c.baseConstant),
		binary.Write(w, byteOrder, uint16(len(c.Value))),
	})
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, []byte(c.Value))
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

func (c *MethodHandleRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

type MethodTypeRef struct {
	baseConstant
	DescriptorIndex ConstPoolIndex
}

func (c *MethodTypeRef) MethodType() *MethodTypeRef { return c }

func (c *MethodTypeRef) Read(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.DescriptorIndex)
}

func (c *MethodTypeRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
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

func (c *InvokeDynamicRef) Dump(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}
