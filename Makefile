build:
	cd api && \
	GOOS=linux go build -o app
	docker build -t equiz ./api