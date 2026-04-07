.PHONY: up down logs clean start

# Menjalankan database dan cache di background
up:
	docker-compose up -d

# Mematikan container
down:
	docker-compose down

# Melihat log dari container
logs:
	docker-compose logs -f

# Mematikan container dan menghapus volume (Reset Total)
clean:
	docker-compose down -v

start:
	go run cmd/api/main.go