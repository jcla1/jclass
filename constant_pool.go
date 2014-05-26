package class

import (
	"encoding/binary"
	"io"
)

type ConstInfoTag uint8
type ConstPoolIndex uint16

type Constant struct {
	Tag ConstInfoTag
	ConstantData
}

type ConstantData interface {
	isConst()

	read(io.Reader) error
	write(io.Writer) error
}

type ConstantPool []*Constant

func (c ConstantPool) read(r io.Reader) (uint16, error) {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return 0, err
	}

	c = make(ConstantPool, 0, count)
	for i := uint16(1); i < count; i++ {
		constant := &Constant{}
		err := constant.read(r)
		if err != nil {
			return 0, err
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
		if constant.Tag == CONSTANT_Long || constant.Tag == CONSTANT_Double {
			i++
		}

		c = append(c, constant)
	}

	return count, nil
}

func (c ConstantPool) write(w io.Writer) error {
	err := binary.Write(w, byteOrder, len(c))
	if err != nil {
		return err
	}

	for _, constant := range c {
		err := constant.write(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Constant) read(r io.Reader) error {
	var err error
	err = binary.Read(r, byteOrder, &c.Tag)
	if err != nil {
		return err
	}

	switch c.Tag {
	case CONSTANT_Class:
		c.ConstantData = &ConstClassInfo{}
	case CONSTANT_FieldRef:
		c.ConstantData = &ConstFieldRefInfo{}
	case CONSTANT_MethodRef:
		c.ConstantData = &ConstMethodRefInfo{}
	case CONSTANT_InterfaceMethodRef:
		c.ConstantData = &ConstInterfaceMethodRefInfo{}
	case CONSTANT_String:
		c.ConstantData = &ConstStringInfo{}
	case CONSTANT_Integer:
		c.ConstantData = &ConstIntegerInfo{}
	case CONSTANT_Float:
		c.ConstantData = &ConstFloatInfo{}
	case CONSTANT_Long:
		c.ConstantData = &ConstLongInfo{}
	case CONSTANT_Double:
		c.ConstantData = &ConstDoubleInfo{}
	case CONSTANT_NameAndType:
		c.ConstantData = &ConstNameAndTypeInfo{}
	case CONSTANT_UTF8:
		c.ConstantData = &ConstUTF8Info{}
	}

	err = c.ConstantData.read(r)
	if err != nil {
		return err
	}

	return nil
}

type ConstClassInfo struct {
	NameIndex ConstPoolIndex
}

func (c *ConstClassInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstClassInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstClassInfo) isConst() {}

type refInfo struct {
	ClassIndex       ConstPoolIndex
	NameAndTypeIndex ConstPoolIndex
}

type ConstFieldRefInfo refInfo

func (c *ConstFieldRefInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstFieldRefInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstFieldRefInfo) isConst() {}

type ConstMethodRefInfo refInfo

func (c *ConstMethodRefInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstMethodRefInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstMethodRefInfo) isConst() {}

type ConstInterfaceMethodRefInfo refInfo

func (c *ConstInterfaceMethodRefInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstInterfaceMethodRefInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstInterfaceMethodRefInfo) isConst() {}

type ConstStringInfo struct {
	StringIndex ConstPoolIndex
}

func (c *ConstStringInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstStringInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstStringInfo) isConst() {}

type ConstIntegerInfo struct {
	Value int32
}

func (c *ConstIntegerInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstIntegerInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstIntegerInfo) isConst() {}

type ConstFloatInfo struct {
	Value float32
}

func (c *ConstFloatInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstFloatInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstFloatInfo) isConst() {}

type ConstLongInfo struct {
	Value int64
}

func (c *ConstLongInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstLongInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstLongInfo) isConst() {}

type ConstDoubleInfo struct {
	Value float64
}

func (c *ConstDoubleInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstDoubleInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstDoubleInfo) isConst() {}

type ConstNameAndTypeInfo struct {
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
}

func (c *ConstNameAndTypeInfo) read(r io.Reader) error {
	return binary.Read(r, byteOrder, c)
}

func (c *ConstNameAndTypeInfo) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstNameAndTypeInfo) isConst() {}

type ConstUTF8Info struct {
	Length uint16
	Value  string
}

func (c *ConstUTF8Info) read(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.Length)
	if err != nil {
		return err
	}

	str := make([]byte, c.Length)
	err = binary.Read(r, byteOrder, str)
	if err != nil {
		return err
	}

	c.Value = string(str)

	return nil
}

func (c *ConstUTF8Info) write(w io.Writer) error {
	return binary.Write(w, byteOrder, c)
}

func (_ *ConstUTF8Info) isConst() {}
