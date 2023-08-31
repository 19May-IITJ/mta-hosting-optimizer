# Project Name

MTA-hosting-optimizer
Service that uncovers the inefficient servers (hostname) hosting only few active MTAs

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)

## Introduction

    Mail transfer agents (MTAs) each having a dedicated public IP address. These MTAs are hosted on the server having unique hostname. This application is used to uncover all those hostname/servers having number of MTAs less than the defined threshold value.

## Features

    It have two HTTP API, with endpoint named as /hostname and /refersh. With initial service up there is no need to refresh data for the first time, application do its of it's own. After that whenever you want you can refersh the data.

## Installation

    Clone the library in you local,
    Step 1:- go test -cover ./... to firstly check all the test cases passing or not.
    Step 2:- go build -o <location binary you want to keep> <relative path of main.go file>
    Step 3:- Set the env variable MTA_THRESHOLD and DBPATH both are optional if not provided the deafult values will be set.

## Usage

    Please hit the localhost:8080 on your local for with endpoints /hostname -> for uncovering server's hostname having MTAs hosted less than the threshold provided in env varibale MTA_THRESHOLD.
