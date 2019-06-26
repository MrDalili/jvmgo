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

//实现接口的方法
func (self CompositeEntry) readClass(className string) ([]byte,Entry,error)  {
	//遍历数组，去读每一个class，如果出错，则返回，没有出错，就
	for _, entry := range self{
		data, from, err := entry.readClass(className)
		if err == nil {
			return data , from , nil
		}
	}
	return nil, nil, errors.New("class not fount:" + className)
}

func (self CompositeEntry) String() string {
	//创建了一个slice
	strs := make([]string ,len(self))
	//对数组中的每个值进行遍历，调用其的string方法
	for i,entry := range self{
		strs[i] = entry.String()
	}
	return strings.Join(strs,pathListSeparartor)
}
