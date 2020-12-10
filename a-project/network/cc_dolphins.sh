#!/bin/bash

if [ $# -ne 2 ]; then
	echo "Arguments are missing. ex) ./cc_tea.sh instantiate 1.0.0"
	exit 1
fi

instruction=$1
version=$2

set -ev

#chaincode install
docker exec cli peer chaincode install -n dolphins -v $version -p github.com/dolphins
#chaincode instatiate
docker exec cli peer chaincode $instruction -n dolphins -v $version -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 3
#chaincode invoke add diver
# {Id: args[0], Name: args[1], Bdate: args[2], Gender: args[3], Btype: args[4]}
docker exec cli peer chaincode invoke -n dolphins -C mychannel -c '{"Args":["addDiver","user1","Anna","04-18-1983","F","BO+"]}'
sleep 3

#chaincode invoke add level
# {id, Levelname:args[1], Org:args[2], Instid: args[3]}
docker exec cli peer chaincode invoke -n dolphins -C mychannel -c '{"Args":["addLevel","user1","openwater","PADI","songjh"]}'
sleep 3

#chaincode invoke add course
# {id, levelname, course}
docker exec cli peer chaincode invoke -n dolphins -C mychannel -c '{"Args":["addCourse","user1","open water","safety"]}'
sleep 3

#chaincode invoke add testresult
# {id, levelname, testresult}
docker exec cli peer chaincode invoke -n dolphins -C mychannel -c '{"Args":["addTestResult","user1","open water","passed"]}'
sleep 3

#chaincode query user1
docker exec cli peer chaincode query -n dolphins -C mychannel -c '{"Args":["getLevel","user1"]}'
docker exec cli peer chaincode query -n dolphins -C mychannel -c '{"Args":["getHistoryForKey","user1"]}'


echo '-------------------------------------END-------------------------------------'
