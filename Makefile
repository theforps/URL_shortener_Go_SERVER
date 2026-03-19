test:
	docker-compose -f docker-compose.test.yml up -d
	go clean -testcache
	go test ./tests -v
	docker-compose -f docker-compose.test.yml down