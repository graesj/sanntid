package main

import (
	. "./elev_manager"
	. "./elev_manager/fsm/driver"
	. "./message"
	. "./network"
	//. "./utilities"
	. "fmt"
	"time"
)

func main() {
	e := Em_makeElevManager()
	buttonChan := make(chan Message, 100)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go e.Em_processElevOrders()
	go Manager(fromMain, toMain)
	go CheckButtons(buttonChan)

	
	broadcastTicker := time.NewTicker(1*time.Second).C

	for {

		select {
		case message := <-toMain:

			switch message.ID {

			case REMOVE_ELEVATOR:
				//e.RemoveElevator(message.Target)

				/*case GENERAL_UPDATE:
				if message.Source != e.Id {
					e.updateElevators(message)
				}
				*/
				//case CALCULATE_COST:
				//	e.calculateCostOfOrder()

				//case BUTTON_EXTERNAL:
				//	if(e.isMaster()){
				//En funksjon som ber alle kalkulere kosten for å ta oppdraget.
				//	}

			case BUTTON_EXTERNAL:
				Println("hhhhhhhhei")
				e.Em_AddExternalOrders(message.Floor, message.ButtonType)

			case ELEVATOR_DATA:
				if message.Elevator.Self_id == e.Self_id {
					Println("Mottok egen melding")
				} else {
					Println("Mottok ny melding")
					_, present := e.Elevators[message.Elevator.Self_id]

					if present { //Update the elevatordata
						e.Em_elevatorUpdate(message.Elevator)
					} else {
						e.Em_newElevator(message.Elevator)
					}
				}
			}

		case buttonMessage := <-buttonChan:

			if buttonMessage.ID == BUTTON_EXTERNAL {
				fromMain <- buttonMessage
				Println("shalla")

			} else if buttonMessage.ID == BUTTON_INTERNAL {
				e.Em_AddInternalOrders(buttonMessage.Floor)
				Println("bais")
			}

		case <- broadcastTicker:
			BroadcastElevatorInfo(*e.Elevators[e.Self_id], fromMain)
			Println(e.Self_id)
			Println(e.Elevators[e.Self_id].Dir)
			Println(e.Elevators[e.Self_id].Floor)
			Println(e.Elevators[e.Self_id].Internal_orders)
		}
	}
}