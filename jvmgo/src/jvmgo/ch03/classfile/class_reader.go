package classfile

import "encoding/binary"

//用来给将class文件当数据流来处理，但是操作字节不方便，所以先定义一个结构体来帮助读取数据
type ClassReader struct {
	data []byte
}

func (self *ClassReader) readUint8() uint8 {
	val := self.data[0]
	self.data = self.data[1:]
	return val
}

func (self *ClassReader) readUit16() uint16{
	//从[]byte中解码多字节数据
	val := binary.BigEndian.Uint16(self.data)
	self.data = self.data[2:]
	return val
}

func (self *ClassReader) readUit32() uint32{
	val := binary.BigEndian.Uint32(self.data)
	self.data = self.data[4:]
	return val
}

func (self *ClassReader) readUint64() uint64{
	val := binary.BigEndian.Uint64(self.data)
	self.data = self.data[8:]
	return val
}
//读取unit16表，表的大小由开头的uint16数据指出
func (self *ClassReader) readUint16s() []uint16{
	n := self.readUit16()
	s := make([]uint16,n)
	for i := range s {
		s[i] = self.readUit16()
	}
	return s
}
//用于读取指定数量的字节
func (self *ClassReader) readBytes(length uint32) []byte{
	bytes := self.data[:length]
	self.data = self.data[length:]
	return bytes
}