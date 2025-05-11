.PHONY: default
default:
	echo 'please input command'

.PHONY: run stop down
run:
	docker-compose up --build -d
stop:
	docker-compose stop
down:
	docker-compose down
.PHONY: test
test:
	go test ./services
