# Terminal 1: Users service
```sh
cd services/users
go run main.go
```

# Terminal 2: Products service
```sh
cd services/products
go run main.go
```

# Terminal 3: Orders service
```sh
cd services/orders
go run main.go
```

# Terminal 4: API Gateway
```sh
cd gateway
go run main.go
```

# Testing
```sh
# Get user
curl http://localhost:8080/users?id=1

# Create order
curl -X POST http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id": "1", "product_id": "101", "amount": 99.99}'
```
