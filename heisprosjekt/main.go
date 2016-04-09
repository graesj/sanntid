package main

import (
	. "./elev_manager"
	. "./elev_manager/fsm/driver"
	. "./message"
	. "./network"
	. "./structs"
	//. "./utilities"
	. "fmt"
	"time"

//"net"
)

func main() {

	e := Em_makeElevManager()
	buttonChan := make(chan Message, 100)
	LampChan := make(chan Message, 100)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go e.Em_processElevOrders(LampChan)
	go Manager(fromMain, toMain)
	go CheckButtons(buttonChan)

	broadcastTicker := time.NewTicker(100 * time.Millisecond).C

	for {

		select {
		case message := <-toMain:

			switch message.ID {

			case REMOVE_ELEVATOR:
				e.ConnectionTimeout(message.Source, fromMain)

			case BUTTON_EXTERNAL:
				ElevSetButtonLamp(message.ButtonType, message.Floor, 1)
				if e.Em_isMaster() {
					Println("Mottok knapp og er master....")
					assignID := e.Determine_target_elev(message.ButtonType, message.Floor)
					message.Source = e.Self_id
					message.ID = ORDER_COMMAND
					message.Target = assignID
					Print("Beste id ble beregnet til å være: ")
					Println(message.Target)
					fromMain <- message
					time.AfterFunc(300*time.Millisecond, func() { e.CheckIfOrderIsTaken(message, fromMain) })

				}

			case NEW_ELEVATOR:
				e.Em_newElevator(message.Elevator)

			case ELEVATOR_DATA:
				if message.Elevator.Self_id == e.Self_id {
					//Println("Mottok egen melding")
				} else {
					//Println("Mottok annen heismelding...")
					_, present := e.Elevators[message.Elevator.Self_id]

					if present { //Update the elevatordata
						e.Em_elevatorUpdate(message.Elevator)
						//	Println("Det var en oppdatering")
					} else {
						e.Em_newElevator(message.Elevator)
						//Println("Det var en ny heis :D")
					}
				}
			case LampID:
				ElevSetButtonLamp(message.ButtonType, message.Floor, 0)
			}

		case buttonMessage := <-buttonChan:

			if buttonMessage.ID == BUTTON_EXTERNAL {
				fromMain <- buttonMessage

				//en bool som sier at ordre har blitt sent
				//kanskje ha en funksjon som sjekker at master plukker den opp
				Println("Sender eksterntknappetrykk ut")

			} else if buttonMessage.ID == BUTTON_INTERNAL {
				e.Em_AddInternalOrders(buttonMessage.Floor, BTN_CMD)
				Println("bais")
			}
		case LampMessage := <-LampChan:
			fromMain <- LampMessage

		case <-broadcastTicker:
			BroadcastElevatorInfo(*e.Elevators[e.Self_id], fromMain)
			//Println(e.Self_id)
			///Println(e.Elevators[e.Self_id].Current_Dir)
			//Println(e.Elevators[e.Self_id].Current_Floor)
			//Println(e.Elevators[e.Self_id].State)
			//Println(e.Elevators[e.Self_id].Planned_Dir)
			Println(e.Elevators[e.Self_id].Internal_orders)

		}
	}
}
