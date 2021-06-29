# Janus

## Overview

**Janus** is an ingester services which consumes data coming from a RabbitMQ message queue. Its main purpose is to process all the raw data to be able to clasify it and store it in different databases. It also provides a query service to retrieve the processed information.

## Diagram

![diagram](./diagram.png)

## TODO

- Query service (gRPC Server)
- Dockerize
- CI/CD
- Docs

## Query Types

- StateQuery
  - deviceId
- StateQuerySubscription
- StatesQuery -> can filter by node
  - nodeId
- TimelineQuery
  - filter
    - type string
    - dimension string
    - value object
    - values object
    - lower object
    - upper object
    - fields filter[]
  - granularity string
  - interval filter
  - aggregations[]
    - type string
    - name string
    - field string
- TimelineSegmentsQuery
  - groupByField
    - dimension
    - name
    - bucketRanges[]
      - name
      - lower
      - upper
  - orderBys[]
    - type
    - dimension
    - aggregation
    - direction
  - limit int

## Response Types

- StateQueryResponse
  - meter & generator
- StatesQueryResponse
  - meter & generator []
- StateQueryStreamResponse
  - meter & generator
- TimelineQueryResponse
  - items map[string]object
  - name
  - count
  - total
- TimelineSegmentsQuery
  - items map[string]segment
    - segment
      - name
      - count
  - total
