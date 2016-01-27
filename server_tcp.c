//Server
#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <netdb.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>

//Define buffer size and UDP port number
#define BUFLEN 1024
#define NPACK 10
#define PORT 20029 //20000 + labplass

//Function for error handling
void diep(char *s)
    {
      perror(s);
      exit(1);
    }
   
int main(void)
    {

      //variabler
      struct sockaddr_in si_me, si_other; //server and client socket
      struct hostent *hostp; //Client host info
      char *hostaddrp; //Dotted decimal host addr string
      int s, i, slen=sizeof(si_other); 
      int n;

      //receive buffer
      char buf[BUFLEN];
       
   
      if ((s=socket(AF_INET, SOCK_STREAM, 0))==-1)
         diep("ERROR socket");

      //zero out the struct
      memset((char *) &si_me, 0, sizeof(si_me));

      si_me.sin_family = AF_INET;
      si_me.sin_port = htons(PORT);
              si_me.sin_addr.s_addr = htonl(INADDR_ANY);
      if (bind(s, (struct sockaddr *) &si_me, sizeof(si_me))==-1)
          diep("bind");

      listen(s,10);
  
      while(1)
      {
          usleep(500*1000);

          //Wait to establish connection
          printf("Waiting for client...\n");
          i = accept(s, (struct sockaddr*)NULL, NULL);
          if (i < 0)
            diep("ERROR on accept");

          //Obatin client information
          hostp = gethostbyaddr((char *)&si_other.sin_addr.s_addr, 
                             sizeof(si_other.sin_addr.s_addr), AF_INET);

          if (hostp == NULL)
            diep("ERROR on gethostbyaddr");

          hostaddrp = inet_ntoa(si_other.sin_addr);
          if (hostaddrp == NULL)
            diep("ERROR on inet_ntoa\n");

          printf("Server established connection with %s (%s) \n",hostp->h_name,hostaddrp);

          //Recieve data from client
          bzero(buf,BUFLEN);
          n = read(i,buf,BUFLEN);

          if(n < 0)
            diep("ERROR reading from socket");
          printf("Server recieved %d bytes: %s\n",n,buf);

          //Echo to client
          n = write(i,buf,strlen(buf));
          if (n < 0)
            diep("ERROR echoing to client");

          close(i);
      }
      
      return 0;
    }