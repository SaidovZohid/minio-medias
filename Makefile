run:
	go run cmd/main.go

swag:
	swag init -g api/api.go -o api/docs