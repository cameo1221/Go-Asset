# Go Asset Management REST API 🛠️

A **production-ready boilerplate REST API for managing assets** built with Golang, PostgreSQL, Docker and Docker Compose. Easy to extend and customize to your needs.

![image](https://github.com/cameo1221/Go-Asset/assets/78523086/03e91ad4-5d21-4124-85a1-7a4b12528e13)


**Highlights:**

✅  Idiomatic Go project structure  

✅  PostgreSQL database dockerized

✅  Repository and service abstraction layers

✅  Graceful error handling middleware  

✅  Request validation middleware

✅  Basic CRUD routes for assets

✅  Integration + unit test coverage

✅  Docker + docker-compose files  

✅  VSCode devcontainer for turn-key development

## 🚀 Getting Started 

### Requirements
* **Go >1.19** 
* **Docker**
* **Docker Compose**  
* **Postman**

### Project Structure and Workflow
* ├── docker-compose.yml
* ├── main.go  
* ├── db/
* │    └──db.go
* ├── handlers/
* │    └── asset_handlers.go
* │    └── admin_handlers.go
* │    └── admin_session_handlers.go
* │    └── employee_handlers.go
* │    └── employeeAsset_handlers.go
* ├── models/ 
* │     └── asset.go
* │     └── admin.go
* │     └── session.go
* │     └── employee.go
* │     └── employee-asset.go
* ├── Middleware/
* │   └── middleware.go
* └── README.md
* **main.go: Program entrypoint that sets up router and http server**
* **db/: SQL scripts for database schema**
* **handlers/: Request route handlers for each resource**
* **models/: Models representing database entities**
* **middleware/: Middleware package for Json header**
