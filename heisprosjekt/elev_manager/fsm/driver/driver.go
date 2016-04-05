package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/

import "C"

func ElevInit() {
	C.elev_init()
}

func ElevGetButtonSignal(buttonType int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(C.int(buttonType)), C.int(floor)))
}

func ElevGetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func ElevGetStopSignal() int {
	return int(C.elev_get_stop_signal())
}

func ElevGetObstructionSignal() int {
	return int(C.elev_get_obstruction_signal())
}

func ElevSetMotorDirection(motorDirection int) {
	(C.elev_set_motor_direction(C.elev_motor_direction_t(C.int(motorDirection))))
}

func ElevSetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func ElevSetFloorLamp(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func ElevSetStopLamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func CheckButtons(fromMain chan, e Elevator){
	for floor := 0; floor < 4; floor++ {
		for buttonType := 0; buttonType < 3; buttonType++{
			if ElevGetButtonSignal(buttonType, floor) {
				if (buttonType == 2){
					//Put de i e sine interne ordre
					
				}

				else {
					message := Message{Id: BUTTON, Dir: buttonType, Floor: floor}
					fromMain <- message

				}
			}
		}
	}
}

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
