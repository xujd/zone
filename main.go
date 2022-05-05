package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"
	"zone.com/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 解析参数
	op := app.NewServerOption()
	op.AddFlags(pflag.CommandLine)
	pflag.Parse()

	// 启动服务
	if err := app.Run(op); err != nil {
		fmt.Fprintf(os.Stderr, "run app failed, %v\n", err)
		os.Exit(1)
	}
}
