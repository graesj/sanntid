#include <stdio.h>
#include "em.h"
#include "elev.h"
#include "fsm.h"
#include "timer.h"
#include "orders.h"


#define stop -1
#define up 0
#define down 1

void em_buttonPressed(){

	if(elev_get_stop_signal()){
		fsm_evStopButtonPressed();
	}

	for (int etg = 0; etg <4; etg++){
		for (int buttonType = 0; buttonType < 3; buttonType++){

			if(etg == 0 && buttonType == 1){			//hindrer assertion fra elev.h
				buttonType++;
			}

			if(etg == 3 && buttonType == 0){			// 			-||-
				buttonType++;
			}

			if(elev_get_button_signal(buttonType, etg) == 1){
				fsm_evButtonPressed(buttonType,etg);
			}

		}
	}
}

void em_isTimeOut(){
	if(timer_isTimeOut() == 1){
		fsm_closeDoors();
		resetEnabled();	
	}
}

void em_newFloor(){
	int etg = elev_get_floor_sensor_signal();
	if (etg != -1){
		fsm_updateFloor(etg);
	}
}

void em_checkQueue(){

	int sum = 0;

	if (fsm_getCurrentDirection() == up){
		
		for (int etg = fsm_getCurrentFloor() + 1; etg < 4; etg++){
			
			if (elevList[etg] == 1){
				fsm_goToFloor(etg);
				sum++;
				break;
			}
		}
			if(sum == 0){
				int sum2 = 0;
			 	for (int etg = 7; etg > 3; etg--){
		 		if (elevList[etg] == 1){
					sum2++;
					fsm_goToFloor(etg - 4);
					break;
				}
			}

				if(sum2 == 0){
					fsm_setCurrentDirection(stop);
				}
			}
		}
		
	else if (fsm_getCurrentDirection() == down){
		int sum = 0;
		
		for (int etg = (fsm_getCurrentFloor() + 3); etg > 3; etg--){
			
			if (elevList[etg] == 1){
				fsm_goToFloor(etg - 4);
				sum++;
				break;
			}
		}

		if(sum == 0){																			
			int sum2 = 0;
			for (int etg = 0; etg < 4; etg++){
		 		if (elevList[etg] == 1){
					sum2++;
					fsm_goToFloor(etg);
					break;
				}
			}

		 	if(sum2 == 0){
		 		fsm_setCurrentDirection(stop);
			}
		}																					
	}

	else if (fsm_getCurrentDirection() == stop){
		
		int check = 0;
		for (unsigned int etg = 0; etg < 8; etg++){
			if (elevList[etg] == 1){
				
				if (etg < 4){
					check++;
					fsm_goToFloor(etg);
				}
				else{
					check++;
					fsm_goToFloor(etg - 4);
				}
			}
		}
			
		if(check == 0 && fsm_getCurrentState() != 3){
			fsm_updateState(2); 	//setter e_state til state_wait hvis kølisten er tom
			printf("State: Wait\n");
		}
	}
}
