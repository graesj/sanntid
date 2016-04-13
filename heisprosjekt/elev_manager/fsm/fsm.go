package fsm

import (
	. "../.././structs"
	. "./driver"
)

func ProcessElevOrders(e Elevator, fromMain chan Message) {
	current_floor := -1
	engineTrouble := false
	lastFloor := e.Current_Floor
	engineCheck := time.NewTimer(3 * time.Second)
	engineCheck.Stop() //

	for {

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
				e.OnMotorError(fromMain)
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
							removeOrders(current_floor, BTN_UP, fromMain, e)
						}
						removeOrders(current_floor, BTN_CMD, fromMain, e)
						e.State = STATE_DOOROPEN
						break
					}
					if e.Planned_Dir == DIR_DOWN {
						e.Furthest_Floor = e.calcFurthestFloor(e.Planned_Dir)

						if e.Internal_orders[BTN_DOWN][current_floor] == 1 {
							if e.Furthest_Floor == current_floor {
								removeOrders(current_floor, BTN_DOWN, fromMain, e)
								e.Current_Dir = DIR_STOP
								e.State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Planned_Dir == DIR_UP {
						if e.Internal_orders[BTN_UP][current_floor] == 1 {
							removeOrders(current_floor, BTN_DOWN, fromMain, e)
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
							removeOrders(current_floor, BTN_DOWN, fromMain, e)
						}
						removeOrders(current_floor, BTN_CMD, fromMain, e)
						e.Current_Dir = DIR_STOP
						e.State = STATE_DOOROPEN
						break
					}
					if e.Planned_Dir == DIR_UP {
						e.Furthest_Floor = e.calcFurthestFloor(e.Planned_Dir)

						if e.Internal_orders[BTN_UP][current_floor] == 1 {
							if e.Furthest_Floor == current_floor {
								removeOrders(current_floor, BTN_UP, fromMain, e)
								e.Current_Dir = DIR_STOP
								e.State = STATE_DOOROPEN
								break
							}
						}
					} else if e.Planned_Dir == DIR_DOWN {
						if e.Internal_orders[BTN_DOWN][current_floor] == 1 {
							removeOrders(current_floor, BTN_UP, fromMain, e)
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
				e.checkForOrdersAndDirChange(fromMain)
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
	}
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

func removeOrders(floor int, button int, lampChan chan Message, e Elevator) { //HEI BEDRE NAVN DA

	switch button {
	
	case BTN_CMD:
		e.Internal_orders[BTN_CMD][floor] = 0
		ElevSetButtonLamp(button, floor, 0)
	
	case BTN_DOWN:
		e.Internal_orders[BTN_DOWN][floor] = 0
		UpdateButtonLamp(BTN_DOWN, floor, fromMain)

	case BTN_UP:
		e.Internal_orders[BTN_DOWN][floor] = 0
		UpdateButtonLamp(BTN_UP,floot,fromMain)
 }
