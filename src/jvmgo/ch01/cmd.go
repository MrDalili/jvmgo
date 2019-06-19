package main
import "flag"
import "fmt"
import "os"
type Cmd struct{
	hrlpFlag bool
	versionFlag bool
	cpOption string
	class string 
	args []string
}