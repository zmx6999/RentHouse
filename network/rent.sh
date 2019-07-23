#!/bin/bash

docker exec cli peer chaincode invoke -n rent -C rent -c '{"args":["clearUser"]}' -o orderer0.house.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer0.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgrent.house.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgrent.house.com/peers/peer0.orgrent.house.com/tls/ca.crt
