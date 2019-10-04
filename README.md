# Inventory

Is an open source Inventory API

## Get Started
- git clone git@github.com:jacky-htg/inventory.git
- cp .env.example .env
- edit .env with your environment
- create database (the database name must be match with your environment)
- go mod init github.com/jacky-htg/inventory
- go run cmd/main.go migrate
- go run cmd/main.go seed
- go run cmd/main.go scan-access
- go test (for run this command, you need docker installed in your laptop)
- go run main.go

## API Testing
- Open your postman application
- Import file inventory.postman_collection.json
- Import file inventory.postman_environment.json
- Call GET /login request in auth directory. username: jackyhtg password:12345678
- Edit current value of token on inventory environment with token in result of login
- Test all request

## How to Add new Module
You can read sample of [add new master](https://github.com/jacky-htg/inventory/blob/master/master.md).

## API Documentation and Specification Program
You can read [API documentation and specification program](https://github.com/jacky-htg/inventory/wiki) in wiki inventory. 