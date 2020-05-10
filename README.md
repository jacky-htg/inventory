# Inventory

Is an open source Inventory API
- Language : Golang
- Database : MySQL 8
- Architecture : Simple MVC
- Router : httprouter
- SQL : database/sql

## Disclaimer
This application only provides Inventory API, if you need the frontend application in desktop/web based/mobile android/mobile ios you can make it your self, or you can contact me by email : rijal.asep.nugroho@gmail.com

To see the direction of developing this application, you can follow the [kanban inventory project](https://github.com/jacky-htg/inventory/projects/1).

## Features
- [x] Multi companies
- [ ] Setting : The company can use the whole of features, or cherry pick part of features.
- [ ] Company registration and verifications 
- [x] Multi regions. One region can be assigned to many branches.
- [x] Multi branches/shops/warehouses
- [x] Master Shelves
- [x] Multi users, multi roles and multi access
- [x] Role Base Access Control (RBAC)
- [x] One user can be assigned multi roles
- [x] One role can be assigned multi access  
- [x] Master products
- [x] Master product categories
- [x] Master brands (brand of products)
- [x] Master customers
- [x] Master suppliers
- [x] Master salesman
- [x] Transaction of purchase
- [x] Transaction of purchase return
- [x] Transaction of good receiving
- [x] Transaction of good receiving return
- [x] Transaction of sales order
- [x] Transaction of sales order return
- [x] Transaction of delivery order
- [x] Transaction of delivery order return
- [ ] Transaction of internal warehouse mutations
- [ ] Transaction of external warehouse mutations
- [ ] Transaction of stock opname
- [x] Transaction of closing stocks
- [ ] Auto suggestions for purchasing order when the product stock is less than the minimum stock
- [ ] Report of users
- [ ] Report of products
- [ ] Report of customers
- [ ] Report of suppliers
- [ ] Report of salesman
- [ ] Report of stock
- [ ] Report of product history (the history of product from receiving in warehouse until delivery to customer)
- [ ] Report of purchase
- [ ] Report of purchase return
- [ ] Report of good receiving
- [ ] Report of good receiving return
- [ ] Report of sales order
- [ ] Report of sales order return
- [ ] Report of delivery order
- [ ] Report of delivery order return
- [ ] Report of internal warehouse mutations
- [ ] Report of external warehouse mutations
- [ ] Report of stock opname

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

## License
The license of application is GPL-3.0, You can use this apllication for commercial use, distribution or modification. But there is no liability and warranty. Please read the [inventory license](https://github.com/jacky-htg/inventory/blob/master/LICENSE) details carefully.