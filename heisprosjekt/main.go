package main

import (
	. "./message"
	"./network"
	"elev_manager"
	"fmt"
	"time"
)

func main() {

	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)
	go network.Manager(fromMain, toMain)
	go checkButtons(fromMain)

	mes := Message{Source: 1, Floor: 1, Target: 1, ID: 1, IP: 1}
	i := 0
	for {
		i = i + 1
		select {
		case m := <-toMain:
			fmt.Println(m.ID)
		case <-time.Tick(500 * time.Millisecond):
			fromMain <- mes
			mes.ID = i
		}

	}
}
