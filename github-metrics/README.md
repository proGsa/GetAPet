# GitHub Development Velocity Dashboard

Локальный стенд для оценки скорости разработки по репозиторию GitHub.

Что поднимается:
- `Elasticsearch` (хранение метрик),
- `Kibana` (визуализация),
- скрипт загрузки данных по `issues`, `pull requests`, `commits`.

## 1) Требования

- Docker + Docker Compose
- Python 3.10+
- GitHub Personal Access Token:
  - для публичного репозитория достаточно `public_repo` (classic) или `Contents: Read + Pull requests: Read + Issues: Read` (fine-grained)
  - для приватного репозитория используйте `repo` (classic) или эквивалентные `Read`-права (fine-grained)

## 2) Переменные окружения

Скопируйте файл и заполните значения:

```bash
cp github-metrics/.env.example github-metrics/.env
```

Пример:

```env
GITHUB_OWNER=proGsa
GITHUB_REPO=GetAPet
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
ELASTIC_URL=http://localhost:9200
LOOKBACK_DAYS=180
```

## 3) Запуск ELK

Из корня проекта:

```bash
docker compose --profile metrics up -d elasticsearch kibana
```

Проверка:
- Elasticsearch: `http://localhost:9200`
- Kibana: `http://localhost:5601`

## 4) Загрузка метрик из GitHub

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r github-metrics/requirements.txt
set -a; source github-metrics/.env; set +a
python github-metrics/load_github_metrics.py
```

Скрипт создаёт и заполняет индексы:
- `github_issues`
- `github_pull_requests`
- `github_commits`

## 5) Настройка в Kibana

1. Откройте `http://localhost:5601`.
2. Перейдите в **Management -> Stack Management -> Data Views**.
3. Создайте Data Views:
   - `github_issues*` (time field: `created_at`)
   - `github_pull_requests*` (time field: `created_at`)
   - `github_commits*` (time field: `committed_at`)
4. В разделе **Discover** можно смотреть сырые данные и строить фильтры.
5. В разделе **Dashboard** создайте три дашборда:
   - `Overview` (общий поток задач/PR/коммитов),
   - `GitHub Issues`,
   - `GitHub Pull Requests`.

## 6) Остановка

```bash
docker compose --profile metrics down
```
