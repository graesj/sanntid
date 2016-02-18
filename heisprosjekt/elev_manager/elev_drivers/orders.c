#include "orders.h"


#define stop -1
#define up 0
#define down 1

void ord_putInQueue(int buttonType, int etg, int dir, int currentFloor){

	//Køen er en liste med 8 elementer. De fire første er bestillinger som skal behandles på vei opp,
	//de fire siste er bestillinger som skal behandles på vei ned. 
	
	if (buttonType != 2){

		if (buttonType == up){
				if (etg == 0){
					elevList[etg + 4] = 1;
				}
				else{
				elevList[etg] = 1;
				}
			}

		else if (buttonType == down){ 
			if (etg == 3){
				elevList[etg] = 1;
			}
			else{
			elevList[etg + 4] = 1;
			}
		}
	}

	if (buttonType == 2){

		if (dir == up){

			if (etg > currentFloor){
				elevList[etg] = 1;
			}
			else if (etg == currentFloor){
				elevList[etg] = 1;
				elevList[etg + 4] = 1;
			}
			else{
				elevList[etg + 4] = 1;
			}
		}

		else if (dir == down){

			if (etg > currentFloor){
				elevList[etg] = 1;
			}
			else if (etg == currentFloor){
				elevList[etg] = 1;
				elevList[etg + 4] = 1;
			}
			else{
				elevList[etg + 4] = 1;
			}
		}

		else if (dir == stop){

			if(etg > currentFloor){
				elevList[etg] = 1;
			}

			else{
				elevList[etg + 4] = 1;
			}
		}
	}
}

void ord_deleteQueue(){
	for (int i = 0; i < 8; i++){
		elevList[i] = 0;
	}
}