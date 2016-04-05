package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"
import (
	"time"
	. "../../.././message"
)

const (
	N_FLOORS = 4
	DIR_UP   = 1
	DIR_DOWN = -1
	DIR_STOP = 0

	//STATES (Golang doesn't have enums)
	STATE_IDLE     = 0
	STATE_RUNNING  = 1
	STATE_DOOROPEN = 2
	STATE_STOP     = 3
)

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

func StopAndOpenDoor() {
	ElevSetMotorDirection(DIR_STOP)
	ElevSetDoorOpenLamp(1)
	//wait 3 seconds
	time.Sleep(time.Second * 2) //use something else than sleep
	ElevSetDoorOpenLamp(0)
}

func CheckButtons(buttonChan chan Message) {
	for floor := 0; floor < 4; floor++ {
		for buttonType := 0; buttonType < 3; buttonType++{
			if ElevGetButtonSignal(buttonType, floor) == 1 {
				if (buttonType == 2){
					
					buttonMessage := Message{ID: BUTTON_INTERNAL, Floor: floor}
					buttonChan <- buttonMessage
					time.Sleep(250*time.Millisecond)
					
				} else {

					buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: buttonType, Floor: floor}
					buttonChan <- buttonMessage
					time.Sleep(250*time.Millisecond)

				}
			}
		}
	}
}
