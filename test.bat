@echo off
echo ######## TEST DB UP ########
echo docker-compose -f docker-compose.test.yml up -d

echo ######## RUN TESTS ########
go clean -testcache
go test ./tests -v

echo ######## TEST DB DOWN ########
echo docker-compose -f docker-compose.test.yml down
pause