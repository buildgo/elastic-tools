

## How to use commands
```
docker-compose up -d  
go run main.go --action info --host http://localhost:9200
go run main.go --action insert --host http://localhost:9200 --index test-elastic --batch 3 --size 10 --period 1 
```