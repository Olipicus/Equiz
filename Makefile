build:
	cd api && \
	GOOS=linux go build -o app
	docker build -t equiz ./api
	rm -f ./api/app

run:
	cd api && \
	GOOS=linux go build -o app
	docker-compose up -d
	rm -f ./api/app