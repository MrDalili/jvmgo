package main

import "fmt"
import "strings"
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
