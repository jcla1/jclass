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
	err := binary.Read(r, byteOrder, &c.FieldsCount)
	if err != nil {
		return err
	}

	c.Fields = make([]*FieldInfo, 0, c.FieldsCount)

	for i := uint16(0); i < c.FieldsCount; i++ {
		var access AccessFlags
		err := binary.Read(r, byteOrder, &access)
		if err != nil {
			return err
		}

		fieldOrMethod, err := readFieldOrMethod(r)
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
		var access AccessFlags
		err := binary.Read(r, byteOrder, &access)
		if err != nil {
			return err
		}

		fieldOrMethod, err := readFieldOrMethod(r)
		if err != nil {
			return err
		}

		method := &MethodInfo{access, *fieldOrMethod}
		c.Methods = append(c.Methods, method)
	}

	return nil
}

func readFieldOrMethod(r io.Reader) (*fieldOrMethodInfo, error) {
	fom := &fieldOrMethodInfo{}

	err := multiError([]error{
		binary.Read(r, byteOrder, &fom.NameIndex),
		binary.Read(r, byteOrder, &fom.DescriptorIndex),
		binary.Read(r, byteOrder, &fom.AttributesCount),
	})

	if err != nil {
		return nil, err
	}

	fom.Attributes = make([]*AttributeInfo, 0, fom.AttributesCount)

	for i := uint16(0); i < fom.AttributesCount; i++ {
		attr, err := readAttribute(r)
		if err != nil {
			return nil, err
		}

		fom.Attributes = append(fom.Attributes, attr)
	}

	return fom, nil
}

func (c *ClassFile) readAttributes(r io.Reader) error {
	attrs, err := readAttributes(r, c.ConstantPool)
	if err != nil {
		return err
	}

	c.Attributes = attrs

	return nil
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
