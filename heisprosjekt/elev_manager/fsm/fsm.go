package fsm

import (
	. "../.././structs"
	. "./driver"
	. "../.././message"
	"time"
	. "fmt"
	"os"
	"sync"

)

func ProcessElevOrders(elev * Elevator, fromMain chan Message) {
	current_floor := -1
	engineTrouble := false
	engineCheck := time.NewTimer(3 * time.Second)
	engineCheck.Stop() //
	mutex := &sync.Mutex{}
	mutex.Lock()
	e := elev
	mutex.Unlock()
	lastFloor := e.Current_Floor


	for {
		//Copy to ensure new data up to date jalla
		mutex.Lock()
		e = elev
		mutex.Unlock()

		switch e.State {
		case STATE_IDLE:
		/*	select {

				case <-engineCheck.C:
				default:
				}

				e.Current_Dir = DIR_STOP
				e.Planned_Dir = DIR_STOP
			*/
			for floor := 0; floor < N_FLOORS; floor++ {
				for buttonType := 0; buttonType < 3; buttonType++ {
					if e.Internal_orders[buttonType][floor] == 1 {
						engineCheck.Reset(3 * time.Second)
						e.State = STATE_RUNNING
					}
				}
			}
		case STATE_RUNNING:

			select {
			case <-engineCheck.C:
				Print(lastFloor)
				Print(" = ")
				Println(e.Current_Floor)
				if lastFloor == e.Current_Floor {
					engineTrouble = true
				} else {
					lastFloor = e.Current_Floor
					engineCheck.Reset(3 * time.Second)
				}

			default:
			}

			if engineTrouble {
				ElevSetMotorDirection(DIR_STOP)
				e.ErrorType = ERROR_MOTOR
				onMotorError(fromMain,e)
				break
			}

			switch e.Current_Dir {
			case DIR_UP:

				ElevSetMotorDirection(DIR_UP)
				current_floor = ElevGetFloorSensorSignal()
				if current_floor != -1 {
					UpdateFloorLight(current_floor)

					e.Current_Floor = current_floor
					if e.Internal_orders[BTN_CMD][current_floor] == 1 {
						e.Current_Dir = DIR_STOP
						if e.Planned_Dir == DIR_UP {
							time.AfterFunc(2*time.Second,func() {removeOrders(current_floor, BTN_UP, fromMain, elev)})
						}
						time.AfterFunc(2*time.Second,func() {removeOrders(current_floor, BTN_CMD, fromMain, elev)})
						e.State = STATE_DOOROPEN
						break
					}
					if e.Planned_Dir == DIR_DOWN {
						e.Furthest_Floor = calcFurthestFloor(e.Planned_Dir, e)

						if e.Internal_orders[BTN_DOWN][current_floor] == 1 {
							if e.Furthest_Floor == current_floor {
								time.AfterFunc(2*time.Second,func() {removeOrders(current_floor, BTN_DOWN, fromMain, elev)})
								e.Current_Dir = DIR_STOP
								e.State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Planned_Dir == DIR_UP {
						if e.Internal_orders[BTN_UP][current_floor] == 1 {
							time.AfterFunc(2*time.Second,func(){removeOrders(current_floor, BTN_UP, fromMain, elev)})
							e.Current_Dir = DIR_STOP
							e.State = STATE_DOOROPEN
							break
						}
					}
				}
			case DIR_DOWN:
				ElevSetMotorDirection(DIR_DOWN)

				current_floor = ElevGetFloorSensorSignal()
				if current_floor != -1 {
					UpdateFloorLight(current_floor)
					e.Current_Floor = current_floor
					if e.Internal_orders[BTN_CMD][current_floor] == 1 {
						if e.Planned_Dir == DIR_DOWN {
							time.AfterFunc(2*time.Second,func(){removeOrders(current_floor, BTN_DOWN, fromMain, elev)})
						}
						time.AfterFunc(2*time.Second,func() {removeOrders(current_floor, BTN_CMD, fromMain, elev)})
						e.Current_Dir = DIR_STOP
						e.State = STATE_DOOROPEN
						break
					}
					if e.Planned_Dir == DIR_UP {
						e.Furthest_Floor = calcFurthestFloor(e.Planned_Dir, e)

						if e.Internal_orders[BTN_UP][current_floor] == 1 {
							if e.Furthest_Floor == current_floor {
								time.AfterFunc(2*time.Second,func() {removeOrders(current_floor, BTN_UP, fromMain, elev)})
								e.Current_Dir = DIR_STOP
								e.State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Planned_Dir == DIR_DOWN {
						if e.Internal_orders[BTN_DOWN][current_floor] == 1 {
							time.AfterFunc(2*time.Second,func(){removeOrders(current_floor, BTN_DOWN, fromMain, elev)})
							e.Current_Dir = DIR_STOP
							e.State = STATE_DOOROPEN
							break
						}
					}
				}
			case DIR_STOP:
				ElevSetMotorDirection(DIR_STOP)
				boll := true
				for floor := 0; floor < N_FLOORS; floor++ {
					for buttonType := 0; buttonType < 3; buttonType++ {
						if e.Internal_orders[buttonType][floor] == 1 {
							boll = false
						}
					}
				}
				if boll {
					lastFloor = e.Current_Floor
					engineCheck.Stop()
					e.Current_Dir = DIR_STOP
					e.Planned_Dir = DIR_STOP
					e.State = STATE_IDLE
					break
				}

				e.Furthest_Floor = -1
				//e.Planned_Dir = DIR_STOP
				checkForOrdersAndDirChange(e, fromMain)
			}

		case STATE_DOOROPEN:
			ElevSetMotorDirection(DIR_STOP)
			ElevSetDoorOpenLamp(1)
			engineCheck.Stop()

			doorTimeout := time.NewTimer(3 * time.Second)
			<-doorTimeout.C

			ElevSetDoorOpenLamp(0)
			doorTimeout.Stop()
			e.State = STATE_RUNNING
			lastFloor = e.Current_Floor
			engineCheck.Reset(3 * time.Second)
			ElevSetMotorDirection(DIR_STOP)
		}
		
		mutex.Lock()
		elev.State = e.State
		elev.Current_Dir = e.Current_Dir
		elev.Planned_Dir = e.Planned_Dir
		e.Furthest_Floor = e.Furthest_Floor
		mutex.Unlock()

		time.Sleep(10*time.Millisecond)
	}
}

func calcFurthestFloor(planned_dir int, e *Elevator) int {
	switch planned_dir {
	case DIR_UP:
		for furthest_floor := 0; furthest_floor < e.Current_Floor; furthest_floor++ {
			if e.Internal_orders[BTN_UP][furthest_floor] == 1 {
				return furthest_floor
			}
		}
	case DIR_DOWN:
		for furthest_floor := N_FLOORS - 1; furthest_floor > e.Current_Floor; furthest_floor-- {
			if e.Internal_orders[BTN_DOWN][furthest_floor] == 1 {
				return furthest_floor
			}
		}
	case DIR_STOP:
		Println("THIS IS THE WRONG PLANNED DIR BRUH")
		return -1
	}
	return e.Furthest_Floor
}

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

func removeOrders(floor int, button int, fromMain chan Message, e *Elevator) { //HEI BEDRE NAVN DA
	mutex := &sync.Mutex{}
	mutex.Lock()
	Println(floor)
	Println(button)
	switch button {
	
	case BTN_CMD:
		e.Internal_orders[BTN_CMD][floor] = 0
		ElevSetButtonLamp(button, floor, 0)
	
	case BTN_DOWN:
		e.Internal_orders[BTN_DOWN][floor] = 0
		if floor != 0 {
			UpdateButtonLamp(BTN_DOWN, floor, fromMain)
		}

	case BTN_UP:
		e.Internal_orders[BTN_UP][floor] = 0
		if floor != N_FLOORS-1 {
			UpdateButtonLamp(BTN_UP,floor,fromMain)
		}
 }
 	Println(e.Internal_orders)
 	mutex.Unlock()
}

func checkForOrdersAndDirChange(elev *Elevator, fromMain chan Message) {
	foundFloor := 0
	mutex := &sync.Mutex{}
	mutex.Lock()
	e := elev
	mutex.Unlock()

	switch e.Planned_Dir {
	case DIR_UP:
		for floor:= 0; floor < N_FLOORS; floor++ {
			if e.Internal_orders[BTN_UP][floor] == 1 {
				foundFloor = 1
				if e.Current_Floor > floor {
					e.Current_Dir = DIR_DOWN
					break
				} else if e.Current_Floor < floor{
					e.Current_Dir = DIR_UP
					break
				} else{
					time.AfterFunc(2*time.Second,func() {removeOrders(e.Current_Floor, BTN_UP, fromMain, elev)})
					e.Current_Dir = DIR_STOP
					e.State = STATE_DOOROPEN
					break
				}
			}
		}
		if foundFloor == 0 {
			e.Planned_Dir = DIR_STOP
		}
	case DIR_DOWN:
		for floor:= N_FLOORS-1; floor >= 0; floor-- {
			if e.Internal_orders[BTN_DOWN][floor] == 1 {
				foundFloor = 1
				if e.Current_Floor > floor {
					e.Current_Dir = DIR_DOWN
					break
				} else if e.Current_Floor < floor {
					e.Current_Dir = DIR_UP
					break
				} else {
					time.AfterFunc(2*time.Second,func() {removeOrders(e.Current_Floor, BTN_UP, fromMain, elev)})
					e.Current_Dir = DIR_STOP
					e.State = STATE_DOOROPEN
					break
				}
			}
		}
		if foundFloor == 0 {
			e.Planned_Dir = DIR_STOP
		}
	case DIR_STOP:
		Println("DIR STOP")
		for buttonType := BTN_CMD; buttonType >= 0; buttonType-- {
			for floor := 0; floor < N_FLOORS; floor++ {
				if e.Internal_orders[buttonType][floor] == 1 {
					foundFloor = 1
					if buttonType == BTN_UP {
						Println("Setting Planned_Dir = DIR_UP")
						e.Planned_Dir = DIR_UP
					} else if buttonType == BTN_DOWN {
						Println("Setting Planned_Dir = DIR_DOWN")
						e.Planned_Dir = DIR_DOWN						
					} else if buttonType == BTN_CMD{
						if e.Current_Floor > floor {
							e.Planned_Dir = DIR_DOWN
						} else if e.Current_Floor < floor {
							Println("DETTE SKJER")
							e.Planned_Dir = DIR_UP
						}
					}
					if e.Current_Floor > floor {
						e.Current_Dir = DIR_DOWN
					} else if e.Current_Floor < floor {
						Println("CORRREECTO")
						e.Current_Dir = DIR_UP
					} else {
						time.AfterFunc(2*time.Second,func() {removeOrders(e.Current_Floor, buttonType, fromMain, elev)})
						e.Current_Dir = DIR_STOP
						e.State = STATE_DOOROPEN
					}
					buttonType = 0 //Asserting break out of nested loop
					break
				}
			}
		}
	}
}

func onMotorError(fromMain chan Message, e *Elevator) {
	fromMain <- Message{ID: REMOVE_ELEVATOR, Source: e.Self_id, Elevator: *e}
	Println("Motortrøbbel. Ring Roger.")
	Print("Feil av type: ")
	Println(e.ErrorType)
	os.Exit(1)

}

func AddExternalOrders(e * Elevator, floor int, buttonType int) {
	var mutex =  &sync.Mutex{}
	mutex.Lock()
	if buttonType == BTN_UP {
		e.Internal_orders[BTN_UP][floor] = 1
		e.External_orders[BTN_UP][floor] = 1

	} else if buttonType == BTN_DOWN {
		e.Internal_orders[BTN_DOWN][floor] = 1
		e.External_orders[BTN_DOWN][floor] = 1

	}
	mutex.Unlock()
}


func AddInternalOrders(e *Elevator, floor int, button int) {
	mutex := &sync.Mutex{}
	mutex.Lock()

	switch button {
	case BTN_UP:
		e.Internal_orders[BTN_UP][floor] = 1
	case BTN_DOWN:
		e.Internal_orders[BTN_DOWN][floor] = 1
	case BTN_CMD:
		ElevSetButtonLamp(BTN_CMD, floor, 1)
		e.Internal_orders[BTN_CMD][floor] = 1
	}
	mutex.Unlock()
	Println(e.Internal_orders)
}