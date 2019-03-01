#!/bin/bash

peer chaincode invoke -n 190222 -C renting -c '{"args":["delHouse"]}' -o orderer0.house.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer0.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgrenting.house.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgrenting.house.com/peers/peer0.orgrenting.house.com/tls/ca.crt
peer chaincode invoke -n 190222 -C renting -c '{"args":["delOrder"]}' -o orderer0.house.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/house.com/orderers/orderer0.house.com/msp/tlscacerts/tlsca.house.com-cert.pem --peerAddresses peer0.orgrenting.house.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/orgrenting.house.com/peers/peer0.orgrenting.house.com/tls/ca.crt