package classpath

import (
	"os"
	"path/filepath"
	"strings"
)

//这个go是用来处理那些带有通配符的
func newWildcardEntry(path string) CompositeEntry {
	//创建了一个动态数组，从0到path的倒数第二个，去掉*号
	baseDir := path[:len(path)-1]
	//创建了一个数组
	compositeEntry := []Entry{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		//看看有没有错误，如果有错返回错误
		if err != nil{
			return err
		}
		//判断路径是不是存在而且路径不等于自己
		//等于自己就说明，没有*号这个东西
		if info.IsDir() && path != baseDir {
			//符合要求则报错，跳出这个目录 skip this directory
			return filepath.SkipDir
		}
		//如果后缀有.jar或者.JAR
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR"){
			//创建一个实例
			jarEntry := newZipEntry(path)
			//把这个实例放进compositeEntry这个数组里
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	//递归遍历目录
	filepath.Walk(baseDir,walkFn)

	return compositeEntry
}

