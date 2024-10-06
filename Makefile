build:
	docker-compose build --no-cache

run:
	docker start omnihr_db
	docker start omnihr_redis
	REDIS_HOST=localhost \
	POSTGRES_DB=omnihr_db \
	POSTGRES_USER=postgres \
	POSTGRES_PASSWORD=password \
	POSTGRES_PORT=5435 \
	JWT_SECRET_KEY=RnNiTjM3d1FrcXZWSm1FcUNXeWVOSUZtWUpFS0hEcnM9SEURE \
	API_SECRET_KEY=a3hBckpHT3h3d0tQRFNLVHZGQlRVR0hQMjR5WmJvQkQ9QkhHTQ== \
	POSTGRES_HOST=localhost \
	go run cmd/server/main.go

up:
	docker-compose up

down:
	docker compose down

restart:
	docker compose restart

clean:
	docker stop go-rest-api-template
	docker stop dockerPostgres
	docker rm go-rest-api-template
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm omnihr-coding-test-backend
	rm -rf .dbdata
