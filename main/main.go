package main

import (
	"flag"
	"path/filepath"

	"github.com/hanjingo/gate"
	"github.com/hanjingo/logger"
	"github.com/hanjingo/util"
)

// for win:   go build -o gate.exe main.go
// for linux: go build -o gate main.go
func main() {
	pwd := util.GetCurrPath()
	confStr := flag.String("f", filepath.Join(pwd, "gate.json"), "输入网关的配置文件绝对路径")
	flag.Parse()

	//设置日志
	//关闭文件打印机
	log := logger.GetDefaultLogger()
	log.GetWriter(logger.DefaultFileWriterName).SetValid(false)

	//加载配置
	log.Info("加载配置文件...")
	conf := gate.NewGateConfig()
	if err := util.LoadJsonConfig(*confStr, conf); err != nil {
		panic(err)
	}
	g := gate.NewGate(conf)
	if g == nil {
		panic("建立网关失败")
	}
	g.Run()
	g.Wait()
}
