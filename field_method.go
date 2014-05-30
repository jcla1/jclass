package class

import (
	"encoding/binary"
	"io"
)

type FieldInfo struct {
	fieldMethodInfo
}

type MethodInfo struct {
	fieldMethodInfo
}

type fieldMethodInfo struct {
	AccessFlags
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	Attributes
}

func readFieldMethod(r io.Reader, constPool ConstantPool) (*fieldMethodInfo, error) {
	fom := &fieldMethodInfo{}

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

func (fom fieldMethodInfo) Dump(w io.Writer) error {
	return multiError([]error{
		binary.Write(w, byteOrder, fom.AccessFlags),
		binary.Write(w, byteOrder, fom.NameIndex),
		binary.Write(w, byteOrder, fom.DescriptorIndex),
		writeAttributes(w, fom.Attributes),
	})
}
