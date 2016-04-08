#!/bin/bash
clear
echo "Username set to 'Student':"

echo "Type in the last byte of the IP to the elvator you want to connect to:"
read IP
echo "Connecting to 129.241.187.$IP"
scp -rq ~/PENISSNAP/sanntid/heisprosjekt/main student@129.241.187.$IP:~/Videos/hv_kjentmann
ssh student@129.241.187.$IP "Videos/hv_kjentmann"

echo Elevator script started


