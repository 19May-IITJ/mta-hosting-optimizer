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

Hosting service serves the /hostnames endpoint. It takes the threshold value (minimum number of active IPs on hostname) via a env variable MTA_THRESHOLD (default value is 1). Apart from that 2 env vairables also need to set which are NATS_URI for building connection with NATS server and HOSTINGSERVICE_PORT which is the port on which the service will service.

## Config Service

Config Service is an interface layer between the DB (currently a JSON File) and Hosting Service. It maintains a cache record of the of DB on local for providing operation of updating the status of MTA IP on given hostname via API handle /refresh and updating records at hosting service end communicated via NATS.

## Installation

Clone the library in you local,
Step 1:- setup NATS on your local or remote or via docker (docker pull nats)
Step 2:- set the env variable of the service you want to host/serve first
for Hosting service MTA_THRESHOLD, NATS_URI*, HOSTINGSERVICE_PORT* need to be set (_ marked are compulsory env no default value)
for Config Service NATS_URI_, CONFIGSERVICE_PORT*, DBPATH* needs to be set (There is a default value for DBPATH but only work once you set it in repo of your own and build the project)
Step 3:- There is a make file included in the project
run command on shell
make bin (for making binary name of binary is fixed - "mta")
make config
make hosting
(NOTE:- if you running multiple instance of hosting service then env of HOSTINGSERVICE_PORT need to be set differently for each instance if multiple instances are on same machine)

## Usage

    Please hit the localhost:8080 on your local for with endpoints /hostname -> for uncovering server's hostname
    having MTAs hosted less than the threshold provided in env varibale MTA_THRESHOLD.
    Also you can hit /refresh to refresh the data set.
