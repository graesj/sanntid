package elev_manager

import (
	. "./network"
	. "./message"
	. "./fsm"
)

type elev_manager struct{
	master int
	self_id int
	external_orders[N_FLOORS*2-2]  int //This is where the orders from the floor panels are put. These orders are broadcasted to all elevators, 
	elevators map[int] *Elevator //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
	message_id int
}

func Em_checkButtons() int {
	return Hello()
}

func Em_makeElevManager() elev_manager{
	var e elev_manager

	//set all parameters of the struct
	return e
}


/*
Check all buttons to see if any action must be performed, and initialize an action
*/
func Em_checkAllFloorButtons(fromMain chan){
	for floor := 0; floor < 4; floor++ {
		for buttonType := 0; buttonType < 3; buttonType++{
			if ElevGetButtonSignal(buttonType, floor) {
				Em_handleFloorButtonPressed(buttonType, floor, fromMain chan)
			}
		}
	}
}


/*
Perform an action based on a button being pressed
*/
func (e * elev_manager) Em_handleFloorButtonPressed(buttonType int, floor int) { 


	if buttonType != 2 {
			message := Message[ID: FLOOR_BUTTON_PRESSED, DIR: UP, FLOOR: 3]
			fromMain <- message

	}
	else {

		//The elevator panel has been used (inside the elevator), and the internal orders of this elevator should be updated. 
		switch floor {
		case 0:
			e.elevators[self_id].internal_orders[floor] = 1
		case 3:
			e.elevators[self_id].internal_orders[floor+2] = 1
		default:
			e.elevators[self_id].internal_orders[floor] = 1
			e.elevators[self_id].internal_orders[floor+2] = 1
		}
	}
}
