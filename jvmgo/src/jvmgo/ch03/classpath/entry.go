package classpath

import "os"
import "strings"

//用于路径分隔符
const pathListSeparartor = string(os.PathListSeparator)

type Entry interface {
	//负责寻找和加载class文件
	//参数是class文件的相对路径，路径之间用斜线（/）分隔，文件名有.class后缀
	//返回值是读到的字节数据、最终定位到class文件的Entry，以及错误信息
	readClass(className string)([]byte,Entry,error)
	//相当于Java中的tostring()
	String() string
}

func newEntry(path string) Entry {
	//看输入的path中是否有该系统的分隔符
	if strings.Contains(path,pathListSeparartor){
		return newCompositeEntry(path)
	}
	//查看路径是否后缀有*
	if strings.HasSuffix(path,"*"){
		return newWildcardEntry(path)
	}
	//看是否为jar包或者zip包
	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path,".ZIP"){
		return newZipEntry(path)
	}
	return newDirEntry(path)
}
