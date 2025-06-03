# Time Tracker Service — FreelanceX

## Overview
A MVP service where freelancer can manage their working hours in minutes using gRPC.

## Tech Stack
- Go (Golang)
- gRPC
- PostgreSQL
- Protocol Buffers

## Folder Structure

.
├── config/
├── handler/
├── model/
├── proto/
├── repository/
├── service/
├── main.go
└── go.mod


## Setup

### 1. Clone & Navigate
```bash
git clone https://github.com/Prototype-1/freelancex_timeTrancker_service.git
cd freelancex_project/time_tracker_service
```

## Install Dependencies

go mod tidy

### Create .env

PORT=50054
DB_URL=postgres://username:password@localhost:5432/timetracker

### Run Migrations

go run scripts/migrate.go

## Start the Service

go run main.go

### Proto Definitions

    proto/time_tracker/time_tracker.proto

#### Maintainers

aswin100396@gmail.com