run:
	nunu run

start:
	docker compose -f docker-compose.dev.yml up

build:
	docker compose -f docker-compose.dev.yml build

serve:
	docker compose -f docker-compose.dev.yml up -d

publish:
	docker buildx build --platform linux/amd64,linux/arm64 -t peterroe/ogimg:latest --push .
