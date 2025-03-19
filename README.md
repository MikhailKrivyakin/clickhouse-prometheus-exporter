ClickHouse Prometheus Exporter
Этот проект представляет собой экспортер метрик из ClickHouse в Prometheus. Он позволяет собирать данные из ClickHouse и экспортировать их в формате, совместимом с Prometheus, для дальнейшего мониторинга и визуализации.

Оглавление
Особенности

Требования

Установка

Конфигурация

Запуск

Примеры запросов

Метрики

Docker

Лицензия

Особенности
Поддержка нескольких серверов ClickHouse.

Динамическое создание метрик на основе конфигурации.

Поддержка пользовательских запросов и лейблов.

Экспорт метрик в формате, совместимом с Prometheus.

Поддержка SSL/TLS для безопасного подключения к ClickHouse.

Требования
Go 1.19 или выше.

ClickHouse (версия 21.x или выше).

Prometheus (для сбора и визуализации метрик).

Установка
Клонируйте репозиторий:

bash
Copy
git clone https://github.com/yourusername/clickhouse-prometheus-exporter.git
cd clickhouse-prometheus-exporter
Установите зависимости:

bash
Copy
go mod download
Соберите проект:

bash
Copy
go build -o clickhouse-prometheus-exporter ./cmd/exporter/main.go
Конфигурация
Конфигурация проекта задаётся в файле config/queries.yaml. Пример конфигурации:

yaml
Copy
servers:
  - name: "server1"
    dsn: "tcp://localhost:9000?username=default&password=&database=default"
    ca_cert_path: "/path/to/ca.crt"  # Путь к CA-сертификату (опционально)

queries:
  - name: "total_users"
    query: "SELECT COUNT(*) as total, department FROM users GROUP BY department"
    metric_name: "clickhouse_total_users"
    help: "Total number of users in the database by department"
    type: "gauge"
    value_column: "total"
    labels: ["department"]
Параметры конфигурации
servers:

name: Имя сервера (используется как лейбл).

dsn: DSN для подключения к ClickHouse.

ca_cert_path: Путь к CA-сертификату (для SSL/TLS, опционально).

queries:

name: Название запроса (для логирования).

query: SQL-запрос к ClickHouse.

metric_name: Имя метрики в Prometheus.

help: Описание метрики.

type: Тип метрики (например, gauge).

value_column: Столбец, используемый как значение метрики.

labels: Столбцы, используемые как лейблы.

Запуск
Запустите экспортер:

bash
Copy
./clickhouse-prometheus-exporter -config /path/to/config.yaml
По умолчанию используется файл config/queries.yaml, если флаг -config не указан.

Метрики будут доступны по адресу:

Copy
http://localhost:8080/metrics
Примеры запросов
Топ-10 баз данных по размеру
sql
Copy
SELECT
    database AS db_name,
    sum(bytes) AS db_size_bytes
FROM system.tables
GROUP BY database
ORDER BY db_size_bytes DESC
LIMIT 10;
Количество активных сессий
sql
Copy
SELECT
    user,
    COUNT(*) AS active_sessions
FROM system.processes
GROUP BY user;
Метрики
Метрики создаются динамически на основе конфигурации. Пример метрики:

Copy
clickhouse_total_users{server="server1", department="HR"} 10
clickhouse_total_users{server="server1", department="IT"} 20
Docker
Вы можете запустить проект в Docker-контейнере:

Соберите Docker-образ:

bash
Copy
docker build -t clickhouse-prometheus-exporter .
Запустите контейнер:

bash
Copy
docker run -d -p 8080:8080 --name clickhouse-exporter clickhouse-prometheus-exporter
Лицензия
Этот проект распространяется под лицензией MIT. Подробности см. в файле LICENSE.

Авторы
Ваше имя

Если у вас есть вопросы или предложения, создайте issue или свяжитесь со мной.

Этот README.md содержит всю необходимую информацию для запуска и использования вашего проекта. Вы можете адаптировать его под свои нужды. 😊