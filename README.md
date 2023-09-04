# Project Name

MTA-hosting-optimizer
Service that uncovers the inefficient servers (hostname) hosting only few active MTAs

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Introduction

Mail transfer agents (MTAs) each having a dedicated public IP address.
These MTAs are hosted on the server having unique hostname.
This project is used to uncover all those hostname/servers having number of MTAs less than the defined threshold value.

## Features

It have two services named as a hosting service and config service each having a HTTP API, with endpoint named as /hostnames and /refersh respectively. Both the two services communicate and exchange data among themseleves using NATS.

## Hosting Service

Hosting service serves the /hostnames endpoint. It takes the threshold value (minimum number of active IPs on hostname) via a env variable MTA_THRESHOLD (default value is 10). Apart from that 2 env vairables also need to set which are NATS_URI for building connection with NATS server and HOSTINGSERVICE_PORT which is the port on which the service will service.

## Config Service

Config Service is an interface layer between the DB (currently a JSON File) and Hosting Service. It maintains a cache record of the of DB on local for providing operation of updating the status of MTA IP on given hostname via API handle /refresh and updating records at hosting service end communicated via NATS.

## Installation

Clone the repo and run "docker compose up —build -d". The docker compose available is for single instance of Hosting service and config service. You can now hit localhost:8010/hostnames for getting the servers names having MTA less than the threshold value set as an env variable in Dockerfile.hosting which is currently 10. You can change the value and re-run the docker compose command. </br > 
Also if you want then you can change the data/ip config by your self via calling the API Handle of config service at localhost:8020/refresh and passing the payload </br >
    e.g:- 
    [
        {   
         "ipAddresses": "127.0.0.1",
         "hostname": "mta-prod-1",
         "status": true
        }
    ].

This will change IP config in local and same will be change in JSON file after 30 sec of API call provided no further request to change any other IP config data is not received in meantime. The project is a single repo multi module code base, both the micro service are design and developed in single repo.

After analyzing the project please run
    docker-compose down --rmi all </br >
Note \*\* If you are working via docker compose then above details for build and installation are fine, for further understand please refer to design docs send over email.

## Usage

Please hit the localhost:8010/hostname on your local for uncovering server's hostname having MTAs hosted less than the threshold provided in env varibale MTA_THRESHOLD set in Dockerfile.hosting.
Also you can hit /refresh to refresh the data set at endpoint localhost:8020/refersh with above demo payload to refersh data if need of any servers's mta IP.
