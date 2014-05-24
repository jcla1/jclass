package class

import (
	"encoding/binary"
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
	case ConstClass, ConstString, ConstMethodType:
		bytesToRead = 2

	case ConstMethodHandle:
		bytesToRead = 3

	case ConstFieldRef, ConstMethodRef,
		ConstInterfaceMehtodRef, ConstInteger,
		ConstFloat, ConstNameAndType, ConstInvokeDynamic:
		bytesToRead = 4

	case ConstLong, ConstDouble:
		bytesToRead = 8

	case ConstUTF8:
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

	var field FieldInfo
	var fieldOrMethod *fieldOrMethodInfo
	for i := uint16(0); i < c.FieldsCount; i++ {
		fieldOrMethod, err = c.readFieldOrMethod(r)
		if err != nil {
			return err
		}

		field = FieldInfo(*fieldOrMethod)
		c.Fields = append(c.Fields, &field)
	}

	return nil
}

func (c *ClassFile) readMethods(r io.Reader) error {
	err := binary.Read(r, byteOrder, &c.MethodsCount)
	if err != nil {
		return err
	}

	c.Methods = make([]*MethodInfo, 0, c.MethodsCount)

	var method MethodInfo
	var fieldOrMethod *fieldOrMethodInfo
	for i := uint16(0); i < c.MethodsCount; i++ {
		fieldOrMethod, err = c.readFieldOrMethod(r)
		if err != nil {
			return err
		}

		method = MethodInfo(*fieldOrMethod)
		c.Methods = append(c.Methods, &method)
	}

	return nil
}

func (c *ClassFile) readFieldOrMethod(r io.Reader) (*fieldOrMethodInfo, error) {
	fom := &fieldOrMethodInfo{}

	errs := []error{
		binary.Read(r, byteOrder, &fom.AccessFlags),
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
