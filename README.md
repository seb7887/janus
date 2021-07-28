# Janus

## Overview

**Janus** is an ingester service which consumes data coming from a RabbitMQ message queue. Its main purpose is to process all the raw data to be able to clasify it and store it in different databases. It also provides a query service to retrieve the processed information.

## How it works

The consumer processor listens to a RabbitMQ queue and when a new message is published, it will execute a new state or telemetry goroutine depending of the message type. The state data is stored in a MongoDB collection and the Telemetry/Log data is stored in different TimescaleDB tables.
There is a gRPC API to get aggregated data from both databases.

## Diagram

![diagram](./diagram.png)

## Getting Started

### Requirements

- RabbitMQ
- TimescaleDB
- MongoDB

### Environment variables

Create a `.env` file in the root directory based on the `env.sample` file.

### Run project locally

You can run Janus locally by executing `make run` in the root directory.
