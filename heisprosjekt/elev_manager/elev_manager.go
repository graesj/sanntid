package elev_manager

import (
	. "./fsm"
	. "./message"
	. "./network"
)

type elev_manager struct {
	master          int
	self_id         int
	external_orders [N_FLOORS*2 - 2]int //This is where the orders from the floor panels are put. These orders are broadcasted to all elevators,
	elevators       map[int]*Elevator   //creates a hash table with 'int' as a keyType, and '*Elevator' as a valueType
	message_id      int
}

func Em_checkButtons() int {
	return Hello()
}

func Em_makeElevManager() elev_manager {
	var e elev_manager

	//set all parameters of the struct
	return e
}
