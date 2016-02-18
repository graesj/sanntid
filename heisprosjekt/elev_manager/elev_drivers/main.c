#include "elev.h"
#include "fsm.h"
#include <stdio.h>
#include "timer.h"
#include "em.h"
#include "orders.h"



int main() {

    if (!elev_init()) {
        printf("Unable to initialize elevator hardware!\n");
        return 1;
    }

    while (fsm_startUp() != 1); 
    
    while(1){  
        em_isTimeOut(); 
    	em_buttonPressed();
        em_newFloor();	
        em_checkQueue();  
    }
    
    return 0;
}
