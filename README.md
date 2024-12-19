# Secure Distributed K-core Decomposition
The project aims to improve the privacy and security of multiple client models using methods like asymmetric encryption in a distributed system. These methods can be applied to decentralized social networks where a central server does not exist and information is exchanged at the client level.

## Features
There is no central server that connects to every single user(client), and the clients only have their information and connection to their neighbors (e.g., friends, parents) where the data transferred is encrypted.

## Scope
● Data Gathering: Gether the data definition of the graph from files and construct it in a distributed system

● Development: Simulate the application in Golang which simulates the simultaneous of multiple nodes communication and releasing its core number

● Evaluation: While assessing the accuracy of the application, the termination is required to be detected as long as all the nodes stop releasing

## Expected outcomes:
● Simulate the Distributed K-core Decomposition algorithm in Golang

● Add the secured Homomorphic Encryption (Comparison & Addition) operation to the algorithm

● Decentralized Termination Detection while all the nodes' communication is over

## How to run the application
go build EECS6412

### golang version
1.23.3
