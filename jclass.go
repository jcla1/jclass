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
	ConstInvokeDynamic                   = 18
)

const (
	AccPublic     AccessFlag = 0x0001
	AccFinal                 = 0x0010
	AccSuper                 = 0x0020
	AccInterface             = 0x0200
	AccAbstract              = 0x0400
	AccSynthetic             = 0x1000
	AccAnnotation            = 0x2000
	AccEnum                  = 0x4000
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
	return binary.Read(r, byteOrder, c.AccessFlags)
}

func (c *ClassFile) readThisClass(r io.Reader) error {
	return binary.Read(r, byteOrder, c.ThisClass)
}

func (c *ClassFile) readSuperClass(r io.Reader) error {
	return binary.Read(r, byteOrder, c.SuperClass)
}
