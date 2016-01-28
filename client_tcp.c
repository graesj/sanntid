//Client


#include <stdio.h> //printf
#include <string.h> //memset
#include <stdlib.h> //exit(0);
#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <netdb.h>
#include <unistd.h>

#define SERVER "129.241.187.146"
#define BUFLEN 1024  //Max length of buffer
#define PORT 20029   //The port on which to send data
 
void die(char *s)
{
    perror(s);
    exit(1);
}
 
int main(int argc, char *argv[])
{
    struct sockaddr_in si_other;
    struct hostent *server;
    int portnr;
    char *hostname; 
    int s,n,i, slen=sizeof(si_other);
    char buf[BUFLEN];

    memset(buf,'\0',BUFLEN);

    strcpy(buf,argv[1]);
    strcat(buf," ");

    for(int k = 2; k < argc; k++){
        strcat(buf,argv[k]);
        strcat(buf," ");
    }

    if ( (s=socket(AF_INET, SOCK_STREAM, 0)) == -1)
    {
        die("socket");
    }
    
    memset((char *) &si_other, 0, sizeof(si_other));
    si_other.sin_family = AF_INET;
    si_other.sin_port = htons(PORT);

    si_other.sin_addr.s_addr = inet_addr(SERVER);
    /*
    if (inet_aton(SERVER , &si_other.sin_addr) == 0) 
    {
        fprintf(stderr, "inet_aton() failed\n");
        exit(1);
    }*/
    
    //Connect to server:
    if (connect(s,(struct sockaddr *)&si_other,sizeof(si_other)) < 0)
        die("ERROR connecting");

   
    
    //Send message to server
    n = write(s,buf,BUFLEN);
    if (n < 0)
        die("ERROR sending to server");
    
    //Recieve echo
    memset(buf,'\0',BUFLEN);
    n = read(s, buf,BUFLEN);
    
    if (n < 0)
        die("ERROR receiving from server");

    printf("Echo from server %s\n",buf);

    close(s);

    

    return 0;
}
