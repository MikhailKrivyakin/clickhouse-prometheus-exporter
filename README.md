
# ClickHouse Prometheus Exporter

  

Этот проект представляет собой экспортер метрик из ClickHouse в Prometheus. Он позволяет собирать данные из ClickHouse и экспортировать их в формате, совместимом с Prometheus, для дальнейшего мониторинга и визуализации.

  

---

  

## Оглавление  

  

1. [Особенности](#особенности)  

2. [Требования](#требования)  

3. [Установка](#установка)  

4. [Конфигурация](#конфигурация)  

5. [Запуск](#запуск)  

6. [Примеры запросов](#примеры-запросов)  

7. [Метрики](#метрики)  

8. [Docker](#docker)  

9. [Лицензия](#лицензия)  

  

---

  

## Особенности  

  

- Поддержка нескольких серверов ClickHouse.  

- Динамическое создание метрик на основе конфигурации.  

- Поддержка пользовательских запросов и лейблов.  

- Экспорт метрик в формате, совместимом с Prometheus.  

- Поддержка SSL/TLS для безопасного подключения к ClickHouse.  

  

---

  

## Требования

  

- Go 1.19 или выше.  
  
- ClickHouse (версия 21.x или выше).  

- Prometheus (для сбора и визуализации метрик).  

  

---

  

## Установка  

  

1. Клонируйте репозиторий:  
 
```bash

git clone https://github.com/yourusername/clickhouse-prometheus-exporter.git

cd clickhouse-prometheus-exporter
```
  

2. Установите зависимости:  
   ```bash
   go mod download
   ```
 3. Соберите проект:  
     ```bash
     go build -o clickhouse-prometheus-exporter ./cmd/exporter/main.go
     ```
 
 ## Конфигурация
 Конфигурация проекта задаётся в файле `config/queries.yaml`. Пример конфигурации:  
 ```yaml
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
 ```
 ### Параметры конфигурации  
 -   **servers**:  
    
    -   `name`: Имя сервера (используется как лейбл).  
        
    -   `dsn`: DSN для подключения к ClickHouse.  
        
-   **queries**:  
    
    -   `name`: Название запроса (для логирования).  
        
    -   `query`: SQL-запрос к ClickHouse.  
        
    -   `metric_name`: Имя метрики в Prometheus.  
        
    -   `help`: Описание метрики.  
        
    -   `type`: Тип метрики (например,  `gauge`).  
        
    -   `value_column`: Столбец, используемый как значение метрики.  
        
    -   `labels`: Столбцы, используемые как лейблы.  

## Запуск
1. Запустите экспортер  
```bash
./clickhouse-prometheus-exporter -config /path/to/config.yaml
```
По умолчанию используется файл `config/queries.yaml`, если флаг `-config` не указан.  
2. Метрики будут доступны по адресу:  
```
http://localhost:8080/metrics
```

## Docker
Вы можете запустить проект в Docker-контейнере:  

1.  Соберите Docker-образ:
    
    ```bash    
    docker build -t clickhouse-prometheus-exporter .
    ```
2.  Запустите контейнер:  
    
    ```bash
    docker run -d -p 8080:8080 -v config.yaml:/ect/clickhouse-prometheus-exporter/config.yaml --name clickhouse-exporter clickhouse-prometheus-exporter
    ```