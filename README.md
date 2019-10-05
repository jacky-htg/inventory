# Inventory

Is an open source Inventory API
- Language : Golang
- Database : MySQL 8
- Architecture : Simple MVC
- Router : httprouter
- SQL : database/sql

## Disclaimer
This application is under develoment. You can follow the progress of project in [inventory project](https://github.com/jacky-htg/inventory/projects/1). You also can view of [milestones](https://github.com/jacky-htg/inventory/milestones).

## Features
- Multi companies
- Setting : The company can use the whole of features, or cherry pick part of features.
- Company registration and verifications 
- Multi regions. One region can be assigned to many branches.
- Multi branches/shops/warehouses
- Multi users, multi roles and multi access
- Role Base Access Control (RBAC)
- One user can be assigned multi roles
- One role can be assigned multi access  
- Master products
- Master product categories
- Master brands (brand of products)
- Master customers
- Master suppliers
- Master salesman
- Transaction of purchase
- Transaction of purchase return
- Transaction of good receiving
- Transaction of good receiving return
- Transaction of sales order
- Transaction of sales order return
- Transaction of delivery order
- Transaction of delivery order return
- Transaction of internal warehouse mutations
- Transaction of external warehouse mutations
- Transaction of stock opname
- Transaction of closing stocks
- Report of users
- Report of products
- Report of customers
- Report of suppliers
- Report of salesman
- Report of stock
- Report of product history (the history of product from receiving in warehouse until delivery to customer)
- Report of purchase
- Report of purchase return
- Report of good receiving
- Report of good receiving return
- Report of sales order
- Report of sales order return
- Report of delivery order
- Report of delivery order return
- Report of internal warehouse mutations
- Report of external warehouse mutations
- Report of stock opname

## Get Started
- git clone git@github.com:jacky-htg/inventory.git
- cp .env.example .env
- edit .env with your environment
- create database (the database name must be match with your environment)
- go mod init github.com/jacky-htg/inventory
- go run cmd/main.go migrate
- go run cmd/main.go seed
- go run cmd/main.go scan-access
- go test -v (To test all of API. For run this command, you need docker installed in your laptop)
- go run main.go

## API Testing
- Open your postman application
- Import file inventory.postman_collection.json
- Import file inventory.postman_environment.json
- Call GET /login request in auth directory. username: jackyhtg password:12345678
- Edit current value of token on inventory environment with token in result of login
- Test all request

## How to Add new Module
This application using golang simple framework. Life cycles is :
```
Request -> Middleware -> Controllers -> Models -> Response
``` 
Directory structure is :
```
> cmd
> controllers
> libraries
> middleware
> models
> payloads
    > request
    > response
> routing
> schema 
```
You can read sample of [add new master](https://github.com/jacky-htg/inventory/blob/master/master.md).

## API Documentation and Specification Program
You can read [API documentation and specification program](https://github.com/jacky-htg/inventory/wiki) in wiki inventory. 