package class

import (
	"encoding/binary"
	"io"
)

type Field struct {
	fieldMethod
}

type Method struct {
	fieldMethod
}

type fieldMethod struct {
	AccessFlags
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	Attributes
}

func readFieldMethod(r io.Reader, constPool ConstantPool) (*fieldMethod, error) {
	fom := &fieldMethod{}

	err := multiError([]error{
		binary.Read(r, byteOrder, &fom.AccessFlags),
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

func (fom fieldMethod) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, fom.AccessFlags),
		binary.Write(w, byteOrder, fom.NameIndex),
		binary.Write(w, byteOrder, fom.DescriptorIndex),
		writeAttributes(w, fom.Attributes),
	})
}
