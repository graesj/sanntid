
//Server
#include <arpa/inet.h>
//#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
//#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>

//Define buffer size and UDP port number
#define BUFLEN 512
#define NPACK 10
#define PORT 20004 //20000 + labplass

//Function for error handling
void diep(char *s)
	{
      perror(s);
       exit(1);
     }
   
     int main(void)
     {

       struct sockaddr_in si_me, si_other;
       int s, i, slen=sizeof(si_other); //receive buffer
       char buf[BUFLEN];
   
       if ((s=socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP))==-1)
         diep("socket");
   	   //zero out the struct
       memset((char *) &si_me, 0, sizeof(si_me));
       si_me.sin_family = AF_INET;
       si_me.sin_port = htons(PORT);
       si_me.sin_addr.s_addr = htonl(INADDR_ANY);
       if (bind(s, (struct sockaddr *) &si_me, sizeof(si_me))==-1)
           diep("bind");
   
   /*    for (i=0; i<NPACK; i++) {
         if (recvfrom(s, buf, BUFLEN, 0, (struct sockaddr *) &si_other, &slen)==-1)
           diep("recvfrom()");
         printf("Received packet from %s:%d\nData: %s\n\n", 
                inet_ntoa(si_other.sin_addr), ntohs(si_other.sin_port), buf);
       }
   */
       while(1)
    {
    	sleep(500);
        printf("Waiting for data...");
        fflush(stdout);
         
        //try to receive some data, this is a blocking call
        if ((recvfrom(s, buf, BUFLEN, 0, (struct sockaddr *) &si_other, &slen)) == -1)
        {
            diep("recvfrom()");
        }
         
        //print details of the client/peer and the data received
        printf("Received packet from %s:%d\n", inet_ntoa(si_other.sin_addr), ntohs(si_other.sin_port));
        printf("Data: %s\n" , buf);
         printf("Received packet from %s:%d\nData: %s\n\n", 
                inet_ntoa(si_other.sin_addr), ntohs(si_other.sin_port), buf);
        //now reply the client with the same data
        /*if (sendto(s, buf, recvfrom, 0, (struct sockaddr*) &si_other, slen) == -1)
        {
            die("sendto()");
        }*/
    }
       close(s);
       return 0;
    }