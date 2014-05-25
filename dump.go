package class

import (
	"encoding/binary"
	"io"
)

var dumpFuncs = []func(*ClassFile, io.Writer) error{
	(*ClassFile).writeMagic,
	(*ClassFile).writeVersion,
	(*ClassFile).writeConstPool,
	(*ClassFile).writeAccessFlags,
	(*ClassFile).writeThisClass,
	(*ClassFile).writeSuperClass,
	(*ClassFile).writeInterfaces,
	(*ClassFile).writeFields,
	// (*ClassFile).writeMethods,
	// (*ClassFile).writeAttributes,
}

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

func (c *ClassFile) writeMagic(w io.Writer) error {
	return binary.Write(w, byteOrder, c.Magic)
}

func (c *ClassFile) writeVersion(w io.Writer) error {
	err := binary.Write(w, byteOrder, c.MinorVersion)
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, c.MajorVersion)
}

func (c *ClassFile) writeConstPool(w io.Writer) error {
	err := binary.Write(w, byteOrder, c.ConstPoolSize)
	if err != nil {
		return err
	}

	for _, constant := range c.ConstPool {
		err := c.writeConstInfo(w, constant)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ClassFile) writeConstInfo(w io.Writer, info *ConstInfo) error {
	var err error

	err = binary.Write(w, byteOrder, info.Tag)
	if err != nil {
		return err
	}

	if info.Tag == CONSTANT_UTF8 {
		err = binary.Write(w, byteOrder, uint16(len(info.Info)))
		if err != nil {
			return err
		}
	}

	err = binary.Write(w, byteOrder, info.Info)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClassFile) writeAccessFlags(w io.Writer) error {
	return binary.Write(w, byteOrder, c.AccessFlags)
}

func (c *ClassFile) writeThisClass(w io.Writer) error {
	return binary.Write(w, byteOrder, c.ThisClass)
}

func (c *ClassFile) writeSuperClass(w io.Writer) error {
	return binary.Write(w, byteOrder, c.SuperClass)
}

func (c *ClassFile) writeInterfaces(w io.Writer) error {
	err := binary.Write(w, byteOrder, c.InterfacesCount)
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, c.Interfaces)
}

func (c *ClassFile) writeFields(w io.Writer) error {
	var err error

	err = binary.Write(w, byteOrder, c.FieldsCount)
	if err != nil {
		return err
	}

	for _, field := range c.Fields {
		err = binary.Write(w, byteOrder, field.AccessFlags)
		if err != nil {
			return err
		}

		err = writeFieldOrMethodInfo(w, field.fieldOrMethodInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeFieldOrMethodInfo(w io.Writer, fom fieldOrMethodInfo) error {
	errs := []error{
		binary.Write(w, byteOrder, fom.NameIndex),
		binary.Write(w, byteOrder, fom.DescriptorIndex),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	for _, attr := range fom.Attributes {
		err := writeAttribute(w, attr)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeAttribute(w io.Writer, attr *AttributeInfo) error {
	errs := []error{
		binary.Write(w, byteOrder, attr.NameIndex),
		binary.Write(w, byteOrder, attr.Length),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	err := binary.Write(w, byteOrder, attr.Info)
	if err != nil {
		return err
	}

	return nil
}
