package main

import (
	. "./elev_manager"
	. "./elev_manager/fsm/driver"
	. "./message"
	. "./network"

	. "fmt"
	//"time"
)

func main() {
	e := Em_makeElevManager()
	buttonChan := make(chan Message, 100)
	fromMain := make(chan Message, 100)
	toMain := make(chan Message, 100)

	go e.Em_processElevOrders()
	go Manager(fromMain, toMain)
	go CheckButtons(buttonChan)

	//msg := Message{Source: 1, Floor: 1, Target: 1, ID: 1}
	i := 0

	for {
		i = i + 1

		select {
		case message := <-toMain:

			switch message.ID {

			case NEW_ELEVATOR:
				Println("ny heis")
				//e.newElevator(message) //Skal legge til den nye heisen, og sjekke hvem som er master

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

			}

		//case msg <- toMain:
		//	fromMain <- msg
		//msg.ID = i

		case buttonMessage := <-buttonChan:

			if buttonMessage.ID == BUTTON_EXTERNAL {
				fromMain <- buttonMessage
				Println("shalla")

			} else if buttonMessage.ID == BUTTON_INTERNAL {
				e.Em_AddInternalOrders(buttonMessage.Floor)
				Println("bais")
			}

		}
	}
}
