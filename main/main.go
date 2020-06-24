package main

import (
	"flag"
	"path/filepath"

	gate "github.com/hanjingo/gate"
	env "github.com/hanjingo/golib/env"
	file "github.com/hanjingo/golib/file"
	logger "github.com/hanjingo/golib/logger"
)

// for win:   go build -o gate.exe main.go
// for linux: go build -o gate main.go
func main() {
	pwd := env.GetCurrPath()
	confStr := flag.String("f", filepath.Join(pwd, "gate.json"), "输入网关的配置文件绝对路径")
	flag.Parse()

	//设置日志
	//关闭文件打印机
	log := logger.GetDefaultLogger()
	log.GetWriter(logger.DefaultFileWriterName).SetValid(false)

	//加载配置
	log.Info("加载配置文件...")
	conf := gate.NewGateConfig()
	if err := file.LoadJsonConfig(*confStr, conf); err != nil {
		panic(err)
	}
	g := gate.NewGate(conf)
	if g == nil {
		panic("建立网关失败")
	}
	g.Run()
	g.Wait()
}
