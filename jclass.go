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
	return multiError([]error{
		binary.Read(r, byteOrder, &c.MinorVersion),
		binary.Read(r, byteOrder, &c.MajorVersion),
	})
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
	var count uint16
	err := binary.Read(r, byteOrder, count)
	if err != nil {
		return err
	}

	c.Interfaces = make([]ConstPoolIndex, count)

	return binary.Read(r, byteOrder, c.Interfaces)
}

func (c *ClassFile) readFields(r io.Reader) error {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return err
	}

	c.Fields = make([]*FieldInfo, 0, count)

	for i := uint16(0); i < count; i++ {
		fieldMethod, err := readFieldMethod(r, c.ConstantPool)
		if err != nil {
			return err
		}

		field := &FieldInfo{*fieldMethod}
		c.Fields = append(c.Fields, field)
	}

	return nil
}

func (c *ClassFile) readMethods(r io.Reader) error {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return err
	}

	c.Methods = make([]*MethodInfo, 0, count)

	for i := uint16(0); i < count; i++ {
		fieldMethod, err := readFieldMethod(r, c.ConstantPool)
		if err != nil {
			return err
		}

		method := &MethodInfo{*fieldMethod}
		c.Methods = append(c.Methods, method)
	}

	return nil
}

func readFieldMethod(r io.Reader, constPool ConstantPool) (*fieldMethodInfo, error) {
	fom := &fieldMethodInfo{}

	err := multiError([]error{
		binary.Read(r, byteOrder, &fom.NameIndex),
		binary.Read(r, byteOrder, &fom.DescriptorIndex),
	})

	if err != nil {
		return nil, err
	}

	fom.Attributes, err = readAttributes(r, constPool)
	if err != nil {
		return nil, err
	}

	return fom, nil
}

func (c *ClassFile) readAttributes(r io.Reader) error {
	var err error
	c.Attributes, err = readAttributes(r, c.ConstantPool)
	return err
}

// Useful when reading from data stream multiple times
func multiError(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
