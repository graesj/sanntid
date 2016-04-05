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
		fmt.Println(e.Elevators[e.Self_id].Dir)
		go e.Em_checkAllFloorButtons()
		go e.Em_processElevOrders()
		fmt.Println(e.Elevators[e.Self_id].Internal_orders)
		time.Sleep(time.Millisecond * 100)
	}

}
