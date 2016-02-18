#include <stdio.h>
#include "fsm.h"
#include "elev.h"
#include "timer.h"
#include "orders.h"

#define TRUE 1
#define FALSE 0

typedef enum {
	state_startup,
	state_running,
	state_wait ,
	state_doorOpen,
	state_stopButton
} elevState;

static elevState e_state = state_wait;

typedef enum{
	firstFloor = 0,
	secFloor = 1,
	thirdFloor = 2,
	fourthFloor = 3,
}floor;

typedef enum{
	up = 0,
	down = 1,
	stop = -1
}dir;

static dir currentDirection = stop;
static floor currentFloor = firstFloor;

static int directionBeforeStopButtonPressed = -1;	//disse to brukes bare til å behandle stopknappen
static int stopButtonPressedBetweenFloors = FALSE;

int fsm_startUp(){
	e_state = state_startup;
	printf("State: Start up\n");

	for (int i = 0; i<8; i++){
		elevList[i] = 0;
	}

	 while(elev_get_floor_sensor_signal() != 0){
           
            elev_set_motor_direction(DIRN_DOWN);
            if (elev_get_stop_signal()){
            	elev_set_stop_lamp(1);
            }
            else{
            	elev_set_stop_lamp(0);
            }
	}

    elev_set_motor_direction(DIRN_STOP);

    currentDirection = stop;
    currentFloor = firstFloor;
    e_state = state_wait; 
    printf("State: Wait\n");

    return 1;
}

void fsm_updateFloor(int etg){
		currentFloor = etg;

		elev_set_floor_indicator(currentFloor);
	
		if (e_state != state_doorOpen){

			if (currentDirection == up){
				if(elevList[currentFloor] == 1){
					fsm_stopAndOpenDoors();
				}
			}
			else if (currentDirection == down){
				if(elevList[currentFloor + 4] == 1){
					fsm_stopAndOpenDoors();
				}
			}
			else if (currentDirection == stop){
				if(elevList[currentFloor] == 1 || elevList[currentFloor + 4] == 1){
					fsm_stopAndOpenDoors();
				} 
			}
		}
}

void fsm_closeDoors(){
	elev_set_door_open_lamp(0);
	e_state = state_running;
	printf("State: Running\n");
}

void fsm_stopAndOpenDoors(){

	e_state = state_doorOpen;
	printf("State: Door open\n");
	elevList[currentFloor] = 0;
	elevList[currentFloor + 4] = 0;
	elev_set_door_open_lamp(1);
	elev_set_motor_direction(DIRN_STOP);
	for(int button = 0; button < 3; button++){
		if (currentFloor == 0 && button == down) { button++;}
		if (currentFloor == 3 && button == up) { button++;}
		elev_set_button_lamp(button, currentFloor, 0);
	}
	timer_start();
}

void fsm_setCurrentDirection(int dir){ 
	currentDirection = dir;
}

void fsm_evButtonPressed(int buttonType, int etg){
	
	switch(e_state){

		case state_wait: {
			e_state = state_running;
			printf("State: Running\n");
			ord_putInQueue(buttonType, etg, currentDirection, currentFloor);
			elev_set_button_lamp(buttonType, etg, 1);
			break;
		}

		case state_running: {
			ord_putInQueue(buttonType, etg, currentDirection, currentFloor);
			elev_set_button_lamp(buttonType, etg, 1);
			break;
		}

		case state_doorOpen:{
			if(etg != currentFloor){
			ord_putInQueue(buttonType, etg, currentDirection, currentFloor);
			elev_set_button_lamp(buttonType, etg, 1);
		}
			break;
		}

		case state_stopButton:{
			break;
		}
		
		case state_startup:{
			break;
		}
	}
}

int fsm_getCurrentState(){
	return e_state;
}

void fsm_evStopButtonPressed(){

	e_state = state_stopButton;
	printf("State: Stop Button Pressed\n");
	directionBeforeStopButtonPressed = currentDirection;
	elev_set_motor_direction(0);

	elev_set_stop_lamp(1);

	for (int etg = 0; etg <4; etg++){
		for (int buttonType = 0; buttonType < 3; buttonType++){

			if(etg == 0 && buttonType == 1){			//hindrer assertion fra elev.h
				buttonType++;
			}

			if(etg == 3 && buttonType == 0){			// 			-||-
				buttonType++;
			}
			elev_set_button_lamp(buttonType, etg, 0);
		}
	}

	if(elev_get_floor_sensor_signal() != -1){
		fsm_stopAndOpenDoors();
	}
	else{
		stopButtonPressedBetweenFloors = TRUE;
	}

	ord_deleteQueue();

	while(elev_get_stop_signal()){
		//Hindrer programmet i å kjøre
	}
	elev_set_stop_lamp(0);
}

void fsm_updateState(int state){
	e_state = state;
}

int fsm_getCurrentFloor(){
	return currentFloor;
}

int fsm_getCurrentDirection(){
	return currentDirection;
}

void fsm_goToFloor(int etg){

	if (stopButtonPressedBetweenFloors == TRUE){				//spesialtilfelle med stoppknappen

		if (etg == currentFloor){
			if (directionBeforeStopButtonPressed == up){
				while (elev_get_floor_sensor_signal() != etg){
					elev_set_motor_direction(DIRN_DOWN);
				}
			}

			else if (directionBeforeStopButtonPressed == down){
				while (elev_get_floor_sensor_signal() != etg){
					elev_set_motor_direction(DIRN_UP);
				}
			}
		}
		directionBeforeStopButtonPressed = -1;
		stopButtonPressedBetweenFloors = FALSE;
	}

	else {

		switch(e_state) {

			case state_running: {
				
				if (currentFloor > etg){
					elev_set_motor_direction(DIRN_DOWN);
					currentDirection = down;
				}
				else if (currentFloor < etg){
					elev_set_motor_direction(DIRN_UP);
					currentDirection = up;
				}
				else {
					elev_set_motor_direction(DIRN_STOP);
					currentDirection = stop;
				}
				break;
			}

			case state_wait: {

				printf("State: Running\n");
				
				if (currentFloor > etg){
					elev_set_motor_direction(DIRN_DOWN);
					currentDirection = down;
					e_state = state_running;
				}
				else if (currentFloor < etg){
					elev_set_motor_direction(DIRN_UP);
					currentDirection = up;
					e_state = state_running;
				}
				else {
					elev_set_motor_direction(DIRN_STOP);
					currentDirection = stop;
					e_state = state_running;
				}
				break;
			}

			case state_doorOpen: {		
				break;
			}

			case state_startup: {
				break;
			}	

			case state_stopButton: {
				break;
			}
		}
	}
}