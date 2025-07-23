# Subscription Service

A simple REST API service to manage user subscriptions.

## Features

- Create, update, delete subscriptions  
- Get a specific subscription (by subscription ID), get all subscriptions 
- Calculate the total subscription price for a certain period with filters by user ID and Service name  
- Swagger API documentation (`/swagger/index.html`)

## Tech Stack

- Go 1.21+  
- PostgreSQL  
- Chi Router  
- Swaggo (Swagger UI)  
- Docker / Docker Compose

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/zapevnik/test-task-for-EffectiveMobile.git
cd test-task-for-EffectiveMobile
```
### 2. Configure the application

**config.yaml:**
```
server:
  host: <your-server-host>
  port: "<your-server-port>"

db:
  host:     <your-db-host>
  port:     <your-db-port>
  user:     <your-db-user>
  password: <your-db-password>
  name:     <your-db-name>

log_level: debug
```

**docker-compose:**
```
environment:
  CONFIG_PATH: /app/config/<your-config-file>
```
⚠️ **Replace all placeholders with your actual database credentials.**


### 3. Run with Docker Compose
```bash
cd deploy
docker-compose up --build
```
With my config service will be available at:
http://localhost:8080

Swagger UI is available at:
http://localhost:8080/swagger/index.html


