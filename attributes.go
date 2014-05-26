package class

import (
	"encoding/binary"
	"io"
)

type Attribute interface {
	isAttr()

	read(io.Reader) error
	write(io.Writer) error
}

type baseAttribute struct {
	NameIndex ConstPoolIndex
	Length    uint32
}

func (ba baseAttribute) read(r io.Reader) error {
	err := binary.Read(r, byteOrder, &ba.NameIndex)
	if err != nil {
		return err
	}

	return binary.Read(r, byteOrder, &ba.Length)
}

func (ba baseAttribute) write(w io.Writer) error {
	err := binary.Write(w, byteOrder, ba.NameIndex)
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, ba.Length)
}

type ConstantValue struct {
	baseAttribute
	Index ConstPoolIndex
}

func (cv *ConstantValue) read(r io.Reader) error {
	err := cv.baseAttribute.read(r)
	if err != nil {
		return err
	}

	return binary.Read(r, byteOrder, &cv.Index)
}

func (cv *ConstantValue) write(w io.Writer) error {
	err := cv.baseAttribute.write(w)
	if err != nil {
		return err
	}

	return binary.Write(w, byteOrder, cv.Index)
}

type Code struct {
	baseAttribute

	MaxStackDepth uint16
	// Warning: Here again, caution: long & double take 2 slots!
	// http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.7.3
	MaxLocalVars uint16

	CodeLength uint32
	Code       []uint8

	// This bit is important for try,catch,finally constructs
	ExceptionsLength uint16
	Exceptions       []struct {
		StartPC   uint16
		EndPC     uint16
		HandlerPC uint16
		CatchType ConstPoolIndex
	}

	AttributesCount uint16
	Attributes      []Attribute
}

func (_ *ConstantValue) isAttr() {}
func (_ *Code) isAttr()          {}
