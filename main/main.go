package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/hanjingo/gate"
	"github.com/hanjingo/util"
)

// for win:   go build -o gate.exe main.go
// for linux: go build -o gate main.go
func main() {
	wg := new(sync.WaitGroup)
	pwd := util.GetCurrPath()
	confStr := flag.String("f", filepath.Join(pwd, "gate.json"), "输入网关的配置文件绝对路径")
	flag.Parse()

	//加载配置
	fmt.Println("加载配置文件...")
	conf := gate.NewGateConfig()
	if err := util.LoadJsonConfig(*confStr, conf); err != nil {
		panic(err)
	}
	g := gate.NewGate(conf)
	if g == nil {
		panic("建立网关失败")
	}
	g.Run(wg)
	wg.Wait()
}
