#ifndef __INCLUDE_FSM
#define __INCLUDE_FSM


int fsm_startUp();
void fsm_evButtonPressed(int buttonType, int etg);
void fsm_evStopButtonPressed();
void fsm_stopAndOpenDoors();
void fsm_updateFloor(int etg);
void fsm_goToFloor(int etg);
void fsm_closeDoors();
void fsm_setCurrentDirection(int dir);
void fsm_updateState(int state);
int fsm_getCurrentState();
int fsm_getCurrentDirection();
int fsm_getCurrentFloor();


#endif