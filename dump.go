package class

import (
	"encoding/binary"
	"io"
)

var dumpFuncs = []func(*ClassFile, io.Writer) error{
	(*ClassFile).writeMagic,
	(*ClassFile).writeVersion,
	// (*ClassFile).writeConstPool,
	// (*ClassFile).writeAccessFlags,
	// (*ClassFile).writeThisClass,
	// (*ClassFile).writeSuperClass,
	// (*ClassFile).writeInterfaces,
	// (*ClassFile).writeFields,
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
		err := binary.Write(w, byteOrder, constant)
		if err != nil {
			return err
		}
	}

	return nil
}
