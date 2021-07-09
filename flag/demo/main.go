package main

import (
	"flag"
	"fmt"
)


// 使用flag.Type来定义命令行参数,并返回一个指针
var n = flag.Bool("n",false,"use a newline")
var sep = flag.String("s"," ","separator of words")

func main() {
	var p = new(int)
	var f = new(float64)
	// 可以使用flag.TypeVar来将变量的声明和与命令行的绑定分开
	flag.IntVar(p,"i",1,"a integer number")
	flag.Float64Var(f,"f",0,"a float number")

	// 解析命令行参数，一定要到所有的命令行参数全部定义完成之后才能调用
	flag.Parse()
	fmt.Println(*p,*f)
	fmt.Printf("第一个剩余的命令行参数是: %v\n",flag.Arg(0))
	fmt.Printf("所有的剩余命令行参数是: %v\n",flag.Args())
	fmt.Printf("解析后剩余的参数数量:%d\n",flag.NArg())
	fmt.Printf("实际解析到了%d个flag\n",flag.NFlag())

}
