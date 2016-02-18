#include "timer.h"
#include <stdio.h>
#include <time.h>


#define TRUE 1
#define FALSE 0

static clock_t g_startTime;
static int enabled;

void timer_start(){
	if (enabled != TRUE){
		g_startTime = clock();
		enabled = TRUE;
	}
}

void resetEnabled(){
	enabled = FALSE;
}

int timer_isTimeOut(){
	if (enabled == FALSE){
		return 0;
	}
	if (((clock() - g_startTime) / CLOCKS_PER_SEC) >= 3 && enabled == TRUE){
		return 1;
	}
	else{
		return 0;
	}
}