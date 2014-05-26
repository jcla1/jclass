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

func (c *ClassFile) readConstPool(r io.Reader) error {
	count, err := c.ConstantPool.read(r)
	if err != nil {
		return err
	}

	c.ConstPoolSize = count

	return nil
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

	c.Interfaces = make([]ConstPoolIndex, c.InterfacesCount)

	return binary.Read(r, byteOrder, c.Interfaces)
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

	err := multiError([]error{
		binary.Read(r, byteOrder, &fom.NameIndex),
		binary.Read(r, byteOrder, &fom.DescriptorIndex),
	})

	if err != nil {
		return nil, err
	}

	count, err := fom.Attributes.read(r, c.ConstantPool)
	if err != nil {
		return nil, err
	}

	fom.AttributesCount = count

	return fom, nil
}

func (c *ClassFile) readAttributes(r io.Reader) error {
	count, err := c.Attributes.read(r, c.ConstantPool)
	if err != nil {
		return err
	}

	c.AttributesCount = count

	return nil
}
