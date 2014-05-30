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

var dumpFuncs = []func(*ClassFile, io.Writer) error{
	(*ClassFile).writeMagic,
	(*ClassFile).writeVersion,
	(*ClassFile).writeConstPool,
	(*ClassFile).writeAccessFlags,
	(*ClassFile).writeThisClass,
	(*ClassFile).writeSuperClass,
	(*ClassFile).writeInterfaces,
	(*ClassFile).writeFields,
	(*ClassFile).writeMethods,
	(*ClassFile).writeAttributes,
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

// Dump writes the binary representation of the
// ClassFile struct to the provied io.Writer
// When a class file is parsed and then dumped
// (unmodified), both (files) should be exactly
// the same.
func (c *ClassFile) Dump(w io.Writer) error {
	var err error

	for _, f := range dumpFuncs {
		err = f(c, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClassFile) readMagic(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.Magic)
}

func (c *ClassFile) writeMagic(w io.Writer) error {
	return binary.Write(w, byteOrder, c.Magic)
}

func (c *ClassFile) readVersion(r io.Reader) error {
	return multiError([]error{
		binary.Read(r, byteOrder, &c.MinorVersion),
		binary.Read(r, byteOrder, &c.MajorVersion),
	})
}

func (c *ClassFile) writeVersion(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, c.MinorVersion),
		binary.Write(w, byteOrder, c.MajorVersion),
	})
}

func (c *ClassFile) readAccessFlags(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.AccessFlags)
}

func (c *ClassFile) writeAccessFlags(w io.Writer) error {
	return binary.Write(w, byteOrder, c.AccessFlags)
}

func (c *ClassFile) readThisClass(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.ThisClass)
}

func (c *ClassFile) writeThisClass(w io.Writer) error {
	return binary.Write(w, byteOrder, c.ThisClass)
}

func (c *ClassFile) readSuperClass(r io.Reader) error {
	return binary.Read(r, byteOrder, &c.SuperClass)
}

func (c *ClassFile) writeSuperClass(w io.Writer) error {
	return binary.Write(w, byteOrder, c.SuperClass)
}

func (c *ClassFile) readInterfaces(r io.Reader) error {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return err
	}

	c.Interfaces = make([]ConstPoolIndex, count)

	return binary.Read(r, byteOrder, c.Interfaces)
}

func (c *ClassFile) writeInterfaces(w io.Writer) error {
	err := binary.Write(w, byteOrder, uint16(len(c.Interfaces)))
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, c.Interfaces)
}

func (c *ClassFile) readFields(r io.Reader) error {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return err
	}

	c.Fields = make([]*Field, 0, count)

	for i := uint16(0); i < count; i++ {
		fieldMethod, err := readFieldMethod(r, c.ConstantPool)
		if err != nil {
			return err
		}

		field := &Field{*fieldMethod}
		c.Fields = append(c.Fields, field)
	}

	return nil
}

func (c *ClassFile) writeFields(w io.Writer) error {
	err := binary.Write(w, byteOrder, uint16(len(c.Fields)))
	if err != nil {
		return err
	}

	for _, field := range c.Fields {
		err := field.Dump(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClassFile) readMethods(r io.Reader) error {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return err
	}

	c.Methods = make([]*Method, 0, count)

	for i := uint16(0); i < count; i++ {
		fieldMethod, err := readFieldMethod(r, c.ConstantPool)
		if err != nil {
			return err
		}

		method := &Method{*fieldMethod}
		c.Methods = append(c.Methods, method)
	}

	return nil
}

func (c *ClassFile) writeMethods(w io.Writer) error {
	err := binary.Write(w, byteOrder, uint16(len(c.Methods)))
	if err != nil {
		return err
	}

	for _, method := range c.Methods {
		err := method.Dump(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClassFile) readAttributes(r io.Reader) error {
	var err error
	c.Attributes, err = readAttributes(r, c.ConstantPool)
	return err
}

func (c *ClassFile) writeAttributes(w io.Writer) error {
	return writeAttributes(w, c.Attributes)
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
