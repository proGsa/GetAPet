# Запуск интеграционных тестов для UserRepository
## Запустить 
 docker compose -f tests/docker-compose.integration.yml up --build --abort-on-container-exit --exit-code-from integration-tests

## Остановить 
 docker compose -f tests/docker-compose.integration.yml down -v