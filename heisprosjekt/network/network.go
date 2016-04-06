package network

import (
	. "../message"
	. "./UDP"
	"time"
	//"fmt"
//	. ".././elev_manager"
	//. ".././elev_manager/fsm"
	. "../structs"
)

var con_timer map[int]*time.Timer

func BroadcastElevatorInfo(e Elevator, UDPsend chan Message) {

	UDPsend <- Message{Source: e.Self_id, ID: ELEVATOR_DATA, Elevator: e}
}

func Manager(fromMain chan Message, toMain chan Message) {

	sendChan := make(chan Message, 50)
	recieveChan := make(chan Message, 50)

	go UDPsend(sendChan)
	go UDPlisten(recieveChan)

	con_timer = make(map[int]*time.Timer)

	for {
		select {
		case message := <-recieveChan:

			if message.ID == SELF_ID {
				_, present := con_timer[message.Source]

				if present { //The ip_key already has a running Timer
					con_timer[message.Source].Reset(3 * time.Second)

				} else { //new elevator
					con_timer[message.Source] = time.AfterFunc(3*time.Second, func() { remove_elev(message.Source, toMain) })
					message.ID = NEW_ELEVATOR

					toMain <- message
				}
			} else {
				toMain <- message
			}

		case message := <-fromMain:

			sendChan <- message
		}
	}

}

func remove_elev(ip_key int, toMain chan Message) {
	m := Message{Source: ip_key, ID: REMOVE_ELEVATOR}
	toMain <- m
	delete(con_timer, ip_key)
}
