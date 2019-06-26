package classfile

import "fmt"

//解析class文件
type ClassFIle struct {
	//magic uint32
	minorVersion uint16
	majorVersion uint16
	constantPool ConstantPool
	accessFlags uint16
	thisClass uint16
	superClass uint16
	interfaces []uint16
	fiedls []*MemberInfo
	methods []*MemberInfo
	attributes []AttributeInfo
}

func Parse(classData []byte) (cf *ClassFile, err error)  {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok{
				err = fmt.Errorf("%v",r)
			}
		}
	}()
	cr := &ClassReader(classData)
	cf = &ClassFIle{}
	cf.read(cr)
	return
}

func (self *ClassFIle) read(reader *ClassFIle)  {

}

func (self *ClassFIle) readAndCheckMagic(reader *ClassReader)  {

}

func (self *ClassFIle) readAndCheckVersion(reader *ClassReader)  {

}

func (self *ClassFIle) MinorVersion() uint16 {

}

func (self *ClassFIle) MajorVersion() uint16 {

}

func (self *ClassFIle) ConstantPool() ConstantPool {

}

func (self *ClassFIle) AccessFlags() uint16{

}

func (self *ClassFIle) Fields() []*MemberInfo  {

}

func (self *ClassFIle) Methods() []*MemberInfo  {

}

func (self *ClassFIle) ClassName() string {

}

func (self *ClassFIle) SuperClassName() string {

}

func (self *ClassFIle) InterfaceName() []string {

}