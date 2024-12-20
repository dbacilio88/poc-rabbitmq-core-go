
docker:
	docker run -d --hostname my-rabbit --name some-rabbit -p 8080:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest rabbitmq:3-management

stop:
	docker container stop some-rabbit
	docker container rm some-rabbit

test:
	go test -v -cover ./...

.PHONY: test