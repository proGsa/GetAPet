# Запуск интеграционных тестов 
## Запустить (из папки backend)
docker compose -f tests/docker-compose.integration.yml up --build --abort-on-container-exit --exit-code-from integration-tests

## Остановить 
docker compose -f tests/docker-compose.integration.yml down -v

# Запуск юнит тестов 
## Запустить (из папки backend)
docker compose -f tests/docker-compose.unit.yml up --build --abort-on-container-exit --exit-code-from unit-tests
## Остановить 
docker compose -f tests/docker-compose.unit.yml down -v