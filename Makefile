#application build
GOOS?=linux
GOARCH?=amd64

APP=go-url-shortener
MAIN_SRC=./cmd/main.go

#deploy
DOCKER_IMAGE?=go-url-shortener-svc

clean:
	rm -f ${APP}

build_local: clean
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o ${APP} ${MAIN_SRC}

docker_build:
	docker build -t ${DOCKER_IMAGE} .

deploy_redis:
	docker-compose up -d redis
down_redis:
	docker-compose rm -s -v -f redis #NOTE: -f flag is used here only to skip prompt (to avoid user interaction)

deploy_svc:
	docker-compose up -d app
down_svc:
	docker-compose rm -s -v -f app #NOTE: -f flag is used here only to skip prompt (to avoid user interaction)

deploy_all:
	docker-compose up --build -d --rm
down_all:
	docker-compose down

deploy_cleanup:
	docker-compose rm -s -v -f






