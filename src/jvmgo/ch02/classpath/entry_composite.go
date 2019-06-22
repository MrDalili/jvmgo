package classpath

import "errors"
import "strings"

type CompositeEntry []Entry

func newCompositeEntry(pathList string) CompositeEntry {
	//创建一个silce切片就像java里面的动态数组一样
	compositeEntry := []Entry{}
	//这里使用的是for range循环，for _表示遍历数组的下标
	for _, path := range strings.Split(pathList, pathListSeparartor){
		//调用接口的方法
		entry := newEntry(path)
		//append是将每次循环的路径进行判断，然后添加到这个slice中
		compositeEntry = append(compositeEntry,entry)
	}
	return compositeEntry
}

func (self CompositeEntry) readClass(className string) ([]byte,Entry,error)  {

}

func (self CompositeEntry) String() string {

}
