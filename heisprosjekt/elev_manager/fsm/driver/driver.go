package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
#include "io.h"
*/
import "C"
import (
	. "../../.././message"
	. "../../.././structs"
)

func ElevGetButtonSignal(buttonType int, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(C.int(buttonType)), C.int(floor)))
}

func ElevGetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
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

func ElevSetButtonLamp(buttonType int, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(int(buttonType)), C.int(floor), C.int(value))
}

func CheckButtons(buttonChan chan Message) {
	for {
		for floor := 0; floor < 4; floor++ {
			for buttonType := 0; buttonType < 3; buttonType++ {
				if (floor == 0 && buttonType == 1) || (floor == N_FLOORS-1 && buttonType == 0) {

				} else {
					if ElevGetButtonSignal(buttonType, floor) == 1 {
						if buttonType == BTN_CMD {
							buttonMessage := Message{ID: BUTTON_INTERNAL, Floor: floor}
							buttonChan <- buttonMessage

						} else {
							buttonMessage := Message{ID: BUTTON_EXTERNAL, ButtonType: buttonType, Floor: floor}
							buttonChan <- buttonMessage

						}
					}
				}
			}
		}
	}
}

func UpdateButtonLamp(button int, floor int, lampChan chan Message) {
	switch button {
	case BTN_CMD:
		ElevSetButtonLamp(button, floor, 0)
	default:
		lampMessage := Message{ID: LAMP_MESSAGE, ButtonType: button, Floor: floor}
		lampChan <- lampMessage
	}
}

func UpdateFloorLight(floor int) {
	ElevSetFloorLamp(floor)
}

func TurnOffAllExternalLights() {
	for button := 0; button < 3; button++ {
		for floor := 0; floor < N_FLOORS; floor++ {
			if !(button == 0 && floor == N_FLOORS-1) && !(button == 1 && floor == 0) {
				ElevSetButtonLamp(button, floor, 0)
			}
		}
	}
}

func RunToFirstFloor() int {
	elevInit()
	if ElevGetFloorSensorSignal() == -1 {
		for ElevGetFloorSensorSignal() == -1 {
			ElevSetMotorDirection(DIR_UP)
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

func elevInit() {
	C.elev_init()
}
