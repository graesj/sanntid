package fsm

import (
	. "../.././structs"
	. "./driver"
)




func RunToFirstFloor() int {
	ElevInit() 

	if ElevGetFloorSensorSignal() == -1 {
		for ElevGetFloorSensorSignal() == -1 {
		}
	}
	ElevSetMotorDirection(DIR_STOP)

	if ElevGetFloorSensorSignal() != 0 {
		for ElevGetFloorSensorSignal() != 0 {
			ElevSetMotorDirection(DIR_DOWN)
		}
	}
	ElevSetMotorDirection(DIR_STOP)

	return 1
}

/*
func ProcessElevOrders(elev Elevator) {
	elev.Internal_orders[]
}
*/
