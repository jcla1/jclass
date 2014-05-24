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
}

const (
	ConstUTF8               ConstInfoTag = 1
	ConstInteger                         = 3
	ConstFloat                           = 4
	ConstLong                            = 5
	ConstDouble                          = 6
	ConstClass                           = 7
	ConstString                          = 8
	ConstFieldRef                        = 9
	ConstMethodRef                       = 10
	ConstInterfaceMehtodRef              = 11
	ConstNameAndType                     = 12
	ConstMethodHandle                    = 15
	ConstMethodType                      = 16
	ConstInvokeDynamik                   = 18
)

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

	switch info.Tag {
	case ConstClass:
		panic("hello!")
	}
	return nil, nil
}
