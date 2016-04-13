package message

import (
	. "../structs"
)

const (

	//Id konstanter

	SELF_ID         = 1 //Messages containing the uniqe elevator id
	NEW_ELEVATOR    = 2
	BUTTON_INTERNAL = 3 //Constant for messages containing orders from internal panel
	BUTTON_EXTERNAL = 4 //Constant for messages containing orders from external panels
	REMOVE_ELEVATOR = 5
	ELEVATOR_DATA   = 6
	ORDER_COMMAND   = 7
	LAMP_MESSAGE    = 8
	GET_UP_TO_DATE  = 9
)

type Message struct {
	Source     int
	Floor      int
	ButtonType int
	Elevator   Elevator
	Target     int
	ID         int
}
