package main

import (
	"fmt"
	"time"

	"github.com/hanjingo/network"

	"github.com/hanjingo/gate/cli/golang"
)

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
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
		}()

		time.Sleep(time.Duration(10) * time.Millisecond)
	}

}
