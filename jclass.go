package class

import (
	"encoding/binary"
	"fmt"
	"io"
)

var byteOrder = binary.BigEndian

var initFuncs = []func(*ClassFile, io.Reader) error{
	(*ClassFile).readMagic,
	(*ClassFile).readVersion,
	(*ClassFile).readConstPool,
	(*ClassFile).readAccessFlags,
	(*ClassFile).readThisClass,
	(*ClassFile).readSuperClass,
	(*ClassFile).readInterfaces,
	(*ClassFile).readFields,
	(*ClassFile).readMethods,
	(*ClassFile).readAttributes,
}

// Parse reads a Java class file from r and, on success,
// returns the parsed struct. Otherwise nil and the error.
func Parse(r io.Reader) (*ClassFile, error) {
	c := &ClassFile{}

	var err error

	for _, f := range initFuncs {
		err = f(c, r)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *ClassFile) readMagic(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Magic)
}

func (c *ClassFile) readVersion(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.MinorVersion)
	if err != nil {
		return err
	}

	return binary.Read(r, byteOrder, &c.MajorVersion)
}

func (c *ClassFile) readConstPool(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.ConstPoolSize)
	if err != nil {
		return err
	}

	c.ConstPool = make([]*ConstInfo, 0, c.ConstPoolSize)

	for i := uint16(1); i < c.ConstPoolSize; i++ {
		info, err := c.readConstInfo(r)
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
		if info.Tag == CONSTANT_Long || info.Tag == CONSTANT_Double {
			i += 1
		}

		c.ConstPool = append(c.ConstPool, info)
	}

	return nil
}

func (c *ClassFile) readConstInfo(r io.Reader) (*ConstInfo, error) {
	info := &ConstInfo{}

	err := binary.Read(r, byteOrder, &info.Tag)
	if err != nil {
		return nil, err
	}

	bytesToRead := uint16(0)
	switch info.Tag {
	case CONSTANT_Class, CONSTANT_String, CONSTANT_MethodType:
		bytesToRead = 2

	case CONSTANT_MethodHandle:
		bytesToRead = 3

	case CONSTANT_FieldRef, CONSTANT_MethodRef,
		CONSTANT_InterfaceMethodRef, CONSTANT_Integer,
		CONSTANT_Float, CONSTANT_NameAndType, CONSTANT_InvokeDynamic:
		bytesToRead = 4

	case CONSTANT_Long, CONSTANT_Double:
		bytesToRead = 8

	case CONSTANT_UTF8:
		err = binary.Read(r, byteOrder, &bytesToRead)
		if err != nil {
			return nil, err
		}

	default:
		panic("unknown tag in class file!")
	}

	info.Info = make([]uint8, bytesToRead)

	err = binary.Read(r, byteOrder, &info.Info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClassFile) readAccessFlags(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.AccessFlags)
}

func (c *ClassFile) readThisClass(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.ThisClass)
}

func (c *ClassFile) readSuperClass(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.SuperClass)
}

func (c *ClassFile) readInterfaces(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.InterfacesCount)
	if err != nil {
		return err
	}

	c.Interfaces = make([]ConstPoolIndex, 0, c.InterfacesCount)

	var idx ConstPoolIndex
	for i := uint16(0); i < c.InterfacesCount; i++ {
		err = binary.Read(r, byteOrder, &idx)
		if err != nil {
			return err
		}

		c.Interfaces = append(c.Interfaces, idx)
	}

	return nil
}

func (c *ClassFile) readFields(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.FieldsCount)
	if err != nil {
		return err
	}

	c.Fields = make([]*FieldInfo, 0, c.FieldsCount)

	for i := uint16(0); i < c.FieldsCount; i++ {
		var access FieldAccessFlag
		err := binary.Read(r, byteOrder, &access)
		if err != nil {
			return err
		}

		fieldOrMethod, err := c.readFieldOrMethod(r)
		if err != nil {
			return err
		}

		field := &FieldInfo{access, *fieldOrMethod}
		c.Fields = append(c.Fields, field)
	}

	return nil
}

func (c *ClassFile) readMethods(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.MethodsCount)
	if err != nil {
		return err
	}

	c.Methods = make([]*MethodInfo, 0, c.MethodsCount)

	for i := uint16(0); i < c.MethodsCount; i++ {
		var access MethodAccessFlag
		err := binary.Read(r, byteOrder, &access)
		if err != nil {
			return err
		}

		fieldOrMethod, err := c.readFieldOrMethod(r)
		if err != nil {
			return err
		}

		method := &MethodInfo{access, *fieldOrMethod}
		c.Methods = append(c.Methods, method)
	}

	return nil
}

func (c *ClassFile) readFieldOrMethod(r io.Reader) (*fieldOrMethodInfo, error) {
	fom := &fieldOrMethodInfo{}

	errs := []error{
		binary.Read(r, byteOrder, &fom.NameIndex),
		binary.Read(r, byteOrder, &fom.DescriptorIndex),
		binary.Read(r, byteOrder, &fom.AttributesCount),
	}

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	fom.Attributes = make([]*AttributeInfo, 0, fom.AttributesCount)

	for i := uint16(0); i < fom.AttributesCount; i++ {
		attr, err := c.readAttribute(r)
		if err != nil {
			return nil, err
		}

		fom.Attributes = append(fom.Attributes, attr)
	}

	return fom, nil
}

func (c *ClassFile) readAttributes(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.AttributesCount)
	if err != nil {
		return err
	}

	c.Attributes = make([]*AttributeInfo, 0, c.AttributesCount)

	var attr *AttributeInfo
	for i := uint16(0); i < c.AttributesCount; i++ {
		attr, err = c.readAttribute(r)
		if err != nil {
			return err
		}

		c.Attributes = append(c.Attributes, attr)
	}

	return nil
}

func (c *ClassFile) readAttribute(r io.Reader) (*AttributeInfo, error) {
	attr := &AttributeInfo{}

	errs := []error{
		binary.Read(r, byteOrder, &attr.NameIndex),
		binary.Read(r, byteOrder, &attr.Length),
	}

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	attr.Info = make([]uint8, attr.Length)

	err := binary.Read(r, byteOrder, &attr.Info)
	if err != nil {
		return nil, err
	}

	return attr, nil
}
