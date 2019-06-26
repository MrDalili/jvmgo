package classpath

import (
	"os"
	"path/filepath"
)

type Classpath struct {
	bootClasspath Entry //jre/lib启动路径
	extClasspath Entry //jre/ext/lib拓展路径
	userClasspath Entry //用户自己定义classPath路径
}
//分别解析系统自带路径和用户定义的路径
func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}//创建实体
	cp.parseBootAndExtClasspath(jreOption)//从系统的库中去寻找
	cp.parseUserClasspath(cpOption)//从用于定义的path中去寻找
	return cp //返回寻找到的路径
}

//从系统的默认路径去寻找
func (self *Classpath) parseBootAndExtClasspath(jreOption string)  {
	//把路径传进去获取到jre所在的目录
	jreDir := getJreDir(jreOption)
	//去jre中的lib这个路径中寻找
	jreLibpath := filepath.Join(jreDir, "lib", "*")
	//实例一个通配符路径的entry，将里面的jar包都解析了，返回一个数组类型，全部存放jar包路径
	self.bootClasspath = newWildcardEntry(jreLibpath)
	//去拓展库里面找
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	self.extClasspath = newWildcardEntry(jreExtPath)
}
//获得jre的路径
func getJreDir(jreOption string) string {
	//判断路径不为空且存在这个路劲
	if jreOption != "" && exists(jreOption) {
		//返回这个路劲
		return jreOption
	}
	//为空或者不存在那个路径，看是不是存在jre这个目录
	if exists("./jre") {
		//返回./jre，表示这个路径中的jre目录
		return "./jre"
	}
	//上面的都不对，那么去系统中的环境变量中寻找，且有这个环境变量
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		//将路径连接
		return filepath.Join(jh,"jre")
	}
	panic("Can not find jre folder!")
}
//看是否存在某个路径
func exists(path string) bool{
	//看是否存在
	if _, err := os.Stat(path); err != nil{
		//出错了以后，判断是不是不存在这个路径的错误
		if os.IsNotExist(err) {
			return false
		}
	}
	//存在这个路径
	return true
}

func (self *Classpath) ReadClass(className string) ([]byte, Entry, error){
	//给要找的class再路径上加上后缀
	className = className + ".class"
	//先从默认库中读，看有没有这个class
	if data, entry, err := self.bootClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	//从拓展库里面读看有没有
	if data, entry, err := self.extClasspath.readClass(className); err == nil {
		return data,entry,err
	}
	//从用户自己的库里面读看有没有
	return self.userClasspath.readClass(className)
	//这里貌似就跟java一样了，到时候调用就到那些特定的类里面去调用特定的方法
}

//toString函数
func (self *Classpath) String() string {
	return self.userClasspath.String()
}
//如果用户-cp/-classPath为空，那么就让他为默认的本地路径
func (self *Classpath) parseUserClasspath(cpOption string) {
	if  cpOption == ""{
		cpOption = "."
	}
	//一个普通的路径
	self.userClasspath = newEntry(cpOption)
}
