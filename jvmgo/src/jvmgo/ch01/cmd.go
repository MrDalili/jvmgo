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
