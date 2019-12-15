package main

import (
	"fmt"

	"github.com/hanjingo/network"

	"github.com/hanjingo/gate/cli/golang"
)

func main() {
	cli := golang.NewGateCli()
	conf := &network.ConnConfig{
		WaitTimeout:  3000,
		WChanCapa:    100,
		RChanCapa:    100,
		ReadBufCapa:  4096,
		WriteBufCapa: 4096,
		ReadLen:      1024,
		NoDelay:      true,
	}
	fmt.Println("拨号结果:", cli.Dial("", "ws://127.0.0.1:10187", "", conf))
}
