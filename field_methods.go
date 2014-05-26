package class

import (
	"encoding/binary"
	"io"
)

type Fields []*FieldInfo

func (f Fields) read(r io.Reader) (uint16, error) {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return 0, err
	}

	f = make(Fields, count)

	for i := uint16(0); i < count; i++ {
		field := &FieldInfo{}
		access, err := field.read(r)
		if err != nil {
			return 0, err
		}

		field.AccessFlags = FieldAccessFlag(access)
		f[i] = field
	}

	return count, nil
}

type Methods []*MethodInfo

func (f Methods) read(r io.Reader) (uint16, error) {
	var count uint16
	err := binary.Read(r, byteOrder, &count)
	if err != nil {
		return 0, err
	}

	f = make(Methods, count)

	for i := uint16(0); i < count; i++ {
		method := &MethodInfo{}
		access, err := method.read(r)
		if err != nil {
			return 0, err
		}

		method.AccessFlags = MethodAccessFlag(access)
		f[i] = method
	}

	return count, nil
}

type FieldInfo struct {
	AccessFlags FieldAccessFlag
	fieldOrMethodInfo
}

type MethodInfo struct {
	AccessFlags MethodAccessFlag
	fieldOrMethodInfo
}

type fieldOrMethodInfo struct {
	NameIndex       ConstPoolIndex
	DescriptorIndex ConstPoolIndex
	AttributesCount uint16
	Attributes
}

func (f *fieldOrMethodInfo) read(r io.Reader) (uint16, error) {
	var access uint16
	err := binary.Read(r, byteOrder, &access)
	if err != nil {
		return 0, err
	}

	err = multiError([]error{
		binary.Read(r, byteOrder, &f.NameIndex),
		binary.Read(r, byteOrder, &f.DescriptorIndex),
		binary.Read(r, byteOrder, &f.AttributesCount),
	})

	if err != nil {
		return 0, err
	}

	return access, nil
}
