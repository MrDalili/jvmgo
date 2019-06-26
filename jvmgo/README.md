# 手写JVM虚拟机

## 我为什么想去自己实现一个JVM虚拟机

我今天在图书馆借阅书籍的时候在《揭秘Java虚拟机》中看到了这么几句话。  

> ​        被真相所蒙蔽，是一件痛苦的事。我们在一个被层层封装的世界里进行开发和设计，操作系统、各种中间件与框架，将底层世界隐藏的严严实实。我们一方面享受着高级编程语言所带来的高效、稳定、快速的开发体验，然而另一方面，却又如同行走于黑暗之中。我们不知道下面是否有坑，即使有坑，我们也不知道如何排除。  
>
> ​        我不知道Java能走多远，未来是否会被淘汰，但你不能因此就否定研究JVM的意义。    

很喜欢这俩句话，我觉得说一个程序员不去研究底层，不去翻阅底层源码只会调用API无疑是一种悲哀。

## 参考书籍

 > 自己动手写Java虚拟机》——张秀宏  
 >《揭秘Java虚拟机-JVM设计原理与实现》——封亚飞

## 环境

windows 10  

JDK 1.8  

Go 1.12.5  

## 为什么使用Go语言

- 你既然想写JVM，那就是你想离及计算机底层更近一些，所以排除JavaScript、python、Ruby等语言，而java又是本语言无挑战，c++又太难太复杂，c过于底层。
- 它是编译型语言，接近底层
- 语法很简洁，学习曲线平滑

## 学习Go的简单语法以及一些常用的包

### 语法

参考菜鸟教程学习即可

### 常用包

- **flag包**  ：是用来实现命令行参数的解析
- **fmt包**：用于在命令行输出和输入的
- **os包**：用于解析电脑系统的相关包

## 一、命令行工具

在命令行中，Java是可以用命令进行调用，所以我们先使用go语言来编写一个很简单的，可以查看vision的命令工具，如下图  

### 准备工作
从跟目录开始配置整个项目目录  
|-jvmgo  
  |-ch01  
    |-cmd.go  
    |main.go  
![Java在命令行中使用java -version][1]    

![自己需要实现的-version][2]  

在这之前我先说一句，所有代码的意思我不单独做说明，因为我已经几乎给每一句话打了注释  
+ 我们在workspace/src/ch01中新建一个cmd.go。用来解析如同java -version这种命令的
```go
package main

//引包
import "flag"
import "fmt"
import "os"

//类似于java中类
type Cmd struct {
	helpFlag    bool
	versionFlag bool
	cpOption    string
	class       string
	args        []string
}

//方法
func parseCmd() *Cmd {
	//new一个Cmd
	cmd := &Cmd{}

	//调用方法
	flag.Usage = printUsage
	//如果是-help，那么显示这个，并将结果赋值给cmd结构体中的helpFlag，下面的大致都如此
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")
	flag.Parse()


	args := flag.Args()
	//如果args的长度大于0
	if len(args) > 0 {
		//将args[0]取出赋值给class
		cmd.class = args[0]
		//将数组第二个元素到最后一个元素赋值给args
		cmd.args = args[1:]
	}

	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-options] class [args...]\n", os.Args[0])
	//flag.PrintDefaults()
}
```
+ 添加一个主类，用来运行整个程序  
```go
package main

import "fmt"

func main() {
	cmd := parseCmd()

	if cmd.versionFlag {
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		startJVM(cmd)
	}
}

func startJVM(cmd *Cmd) {
	fmt.Printf("classpath:%s class:%s args:%v\n",
		cmd.cpOption, cmd.class, cmd.args)
}
```
这样我们就已经完成了这个小demo的编写。  
这也只是我们迈出的第一步，接下来我们要做的就是去加载类文件了。

## 二、搜索class文件
在上一讲中，我们实现了通过传入命令的方式去启动java的虚拟机，或者查看版本等功能，那么这一讲主要就是通过程序告诉你虚拟机是如何加载一个class文件的。  
### 类路径
Java虚拟机规范虽然没有规定虚拟机从哪里寻找类，这里使用的是Oracle的Java虚拟机实现根据类路径来搜素，依次按照
1. 启动类路径
   + 默认对应的是jre/lib目录，Java标准中大部分的库也都在这里
2. 扩展类路径
   + 默认对应在jre/lib/ext目录下，Java的扩展类一般都在这里
3. 用户类路径  
   + 存放自己实现的类与第三方的类库，一般可以通过xbootclasspatch选项修改启动类路径
   + 一般默认是.路径
   + 一般使用-cp/-classpath来修改用户此次运行的时候的默认路径
   + 可指定目录，JAR,ZIP等文件
   + 也可以指定多个文件和目录
   + 也可以使用通配符*来指定某个目录下的所有JAR文件
### 准备工作
配置路径，在上一讲中原来的路径的基础上，创建ch02，并将ch01中的文件拷贝进ch02中  
|-jvmgo  
  |-ch01  
  |-ch02  
    |-classpath  
    |-cmd.go  
    |-main.go  
打开cmd.go,为结构体添加一个字段变成如下：
```go
//类似于java中类
type Cmd struct {
	helpFlag    bool
	versionFlag bool
	cpOption    string
	XjreOption string //添加的字段，是用来接受指定的路径的
	class       string
	args        []string
}
```
添加一个解析该Xjre的语句  
```go
flag.Usage = printUsage
	//如果是-help，那么显示这个，并将结果赋值给cmd结构体中的helpFlag，下面的大致都如此
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")
	//这一条
	flag.StringVar(&cmd.XjreOption,"Xjre","","path to jre")
	flag.Parse()  
```
现在的cmd.go长这样
```go
package main

//引包
import "flag"
import "fmt"
import "os"

//类似于java中类
type Cmd struct {
	helpFlag    bool
	versionFlag bool
	cpOption    string
	XjreOption string
	class       string
	args        []string
}

//方法
func parseCmd() *Cmd {
	//new一个Cmd
	cmd := &Cmd{}

	//调用方法
	flag.Usage = printUsage
	//如果是-help，那么显示这个，并将结果赋值给cmd结构体中的helpFlag，下面的大致都如此
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")
	flag.StringVar(&cmd.XjreOption,"Xjre","","path to jre")
	flag.Parse()


	args := flag.Args()
	//如果args的长度大于0
	if len(args) > 0 {
		//将args[0]取出
		cmd.class = args[0]
		//
		cmd.args = args[1:]
	}

	return cmd
}

func printUsage() {
	fmt.Printf("Usage: %s [-options] class [args...]\n", os.Args[0])
	//flag.PrintDefaults()
}
```
### 实现具体方法
我们现在要想，如何让它来寻找传进来的class，那就是先从jre/lib，再到jre/ext/lib，再到用户所提供的路径，所以我们也这么来，但是我们还得去解析JAR，ZIP，遍历目录等操作。  
#### 定义一个Entry接口
将这个文件定义在classpath下，它的作用就如同Java中的接口一样，实现特定的一个类，然后重写方法，也用貌似Java中”多态“的感觉  
代码如下：  
**classpath/entry.go**  
```go
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

```
#### DirEntry，对应解虚普通路径的Entry
在4种形式中，这个实现对应的就是很普通的路径，如：
> c:/programmer/java1.8/jre   

在此顺带一提的就是，在Go语言中，类对接口的实现不用写什么，直接去实现对应的方法就ok了，该go中的逻辑与实现过程都在代码中有注释，就不再做过多的解释
**classpath/entry_dir.go**
```go
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
```

#### ZipEntry
在4中形式中，这个实现对应的是后缀为.JAR或.ZIP的这种路径，如:  
> c:/programmer/system64/test.jar

**classpath/entry_zip.go**
```go
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
```
#### CompositeEntry
在4种形式中，这个实现对应的就是很一次传递多个路径的参数时所要进行使用的，如：
> c:/programmer/system64/test.jar,c:/programmer/system64/lib

**classpath/entry_composite.go**
```go
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
```
#### WildcardEntry
在4中形式中，这个实现对应的是存在通配符的路径形式，如:
> c:/programmer/system64/*

**classpath/entry_wildcard.go**
```go
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


```
#### Classpath
这个结构体是用来将前面的各种实现与接口进行组合和调用  
**classpath/classpath.go**
```go
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

```
#### 修改主类
因为我们添加了很多的实现，所以我们要将其在主函数中可以调用，内容如下:   
**ch02/main.go**
```go
package main

import "fmt"
import "strings"的
import "jvmgo/ch02/classpath"

func main() {
	//实例化对象
	cmd := parseCmd()
	if cmd.versionFlag {
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		startJVM(cmd)
	}
}
//启动jvm
func startJVM(cmd *Cmd) {
	//将传进来的路径进行解虚
	//一个是修改后的系统启动路径，一个是用户指定启动路径
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	//输出找的路径
	fmt.Printf("classpath:%v class:%v args:%v\n",cp, cmd.class, cmd.args)
	//把原来的.都换成/，最后一个不转
	className := strings.Replace(cmd.class, ".", "/", -1)
	//从路径中读这个class的内容
	classData, _, err := cp.ReadClass(className)
	//看看报错不报错
	if err != nil {
		fmt.Print("Could not find or load main class %s\n",cmd.class)
		return
	}
	//输出class内容
	fmt.Printf("class data:%v\n", classData)
}
```
#### 大功告成，只差测试
1. 我们先对整个ch02进行编译，在bin目录下调用生成的exe，我们接下来使用它
2. 打开终端，进入bin目录
3. 执行命令
> $ go_build_jvmgo_ch02 -Xjre "C:\programmer\Java_JDK1.8\jre" java.lang.Object
> 这里go_build_jvmgo_ch02是我生成的可以执行的exe文件，这里换成你生成的文件

4. 查看结果，如下图：  
![显示结果][3]  

虽然说现在的文件读出来的时二进制，但我觉得能把这个做出来真的真的很厉害。我一直在翻译这个go语言，我以前没有学过，有查资料，有学基础语法，中间很累，但我觉得当这个二进制文件显示在控制台的时候，这时才是对自己最大的奖励哈哈哈哈，so happy~  
![开心][4]  

## 三、解析class文件

### Go和Java语言基本数据类型对照关系

| Go语言类型    | Java语言类型 | 说明                 |
| ------------- | ------------ | -------------------- |
| int8          | byte         | 8比特有符号整数      |
| uint8         | N/A          | 8比特无符号整数      |
| int16         | short        | 16比特有符号整数     |
| uint16        | char         | 16比特无符号整数     |
| int32（rnme） | int          | 32比特有符号整数     |
| uint32        | N/A          | 32比特无符号整数     |
| int64         | long         | 64比特无符号整数     |
| uint64        | N/A          | 64比特无符号整数     |
| float32       | float        | 32比特IEEE-754浮点数 |
| float64       | double       | 63比特IEEE-754浮点数 |



  

  

[1]: https:qiniuyun.ningdali.com/blog/19620jvm.png	"java自带的"
[2]: https:qiniuyun.ningdali.com/blog/19620jvm1.png	"我们需要实现的"
[3]: https:qiniuyun.ningdali.com/blog/19622jvm.png	"显示结果"
[4]: https:qiniuyun.ningdali.com/blog/19622jvm1.jpeg	"开心"

