package classpath

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"path/filepath"
)

type ZipEntry struct {
	//存放ZIP或JAR文件的绝对路径
	absPath string
}

func newZipEntry(path string) *ZipEntry{
	absPath,err := filepath.Abs(path)
	if err != nil{
		panic(err)
	}
	return &ZipEntry{absPath}
}

func (self *ZipEntry) String() string{
	return self.absPath
}

//如何从ZIP中提取class文件
func (self *ZipEntry) readClass(className string) ([]byte, Entry, error)  {
	//首先打开zip文件，如果出错就退出
	r, err := zip.OpenReader(self.absPath)
	if err != nil {
		return nil, nil, err
	}
	//defer关键字的意思？

	//保证打开的文件可以关闭
	defer r.Close()
	//迭代压缩文件中的文件
	for _, f := range r.File {
		if f.Name == className {
			//再打开文件
			rc, err := f.Open()
			if err != nil {
				return nil, nil, err
			}
			//关闭文件
			defer rc.Close()
			//把里面的内容全部都读了
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, nil, err
			}
			//将数据，实体，还有错误返回
			return data, self, nil
		}
	}
	//没有找到就包class not found
	return nil, nil, errors.New("class not fount: " + className)
}

