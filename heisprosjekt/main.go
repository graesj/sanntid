package main

import (
	. "./elev_manager"
	. "./elev_manager/fsm"
	. "./elev_manager/fsm/driver"

	//"fmt"
	//"time"
)

func main() {

	Em_checkButtons()
	Fsm_createElev()
	ElevSetStopLamp(1)
}
