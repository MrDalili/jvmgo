package classpath

import "io/ioutil"
import "path/filepath"

type DirEntry struct {
	absDir string
}

//在文中将其当作构造函数使用
func newDirEntry(path string)  *DirEntry{
	//将路径转为绝对路径
	absDir, err := filepath.Abs(path)
	if err != nil {
		//如果出现错误，则调用panic种植程序
		panic(err)
	}
	//没有错误创建实例并返回
	return &DirEntry{absDir}
}

//这个是go语言里面的函数
func (self *DirEntry) readClass(className string) ([]byte , Entry, error) {
	//将目录和class文件名拼成一个完整的路径
	fileName := filepath.Join(self.absDir,className)
	//读class文件的内容
	data, err := ioutil.ReadFile(fileName)
	//返回内容，实体，以及错误
	return data, self, err
}

func (self *DirEntry) String() string  {
	//返回路径
	return self.absDir
}
