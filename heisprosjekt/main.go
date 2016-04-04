package main

import (
	. "./elev_manager"
	//. "./elev_manager/fsm"
	//. "./elev_manager/fsm/driver"

	"fmt"
	"time"
)

func main() {

	e := Em_makeElevManager()
	fmt.Println(e.Self_id)
	for {
		e.Em_checkAllFloorButtons()
		fmt.Println(e.Elevators[e.Self_id].Internal_orders)
		time.Sleep(time.Millisecond * 100)
	}

}
