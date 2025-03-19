package main

import (
	"clickhouse-prometheus-exporter/internal/clickhouse"
	"clickhouse-prometheus-exporter/internal/config"
	"clickhouse-prometheus-exporter/internal/exporter"
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Определяем флаг для пути к конфигурационному файлу
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Создаем клиенты для каждого сервера ClickHouse
	clients := make(map[string]*clickhouse.Client)
	for _, server := range cfg.Servers {
		client, err := clickhouse.NewClient(server.DSN)
		if err != nil {
			log.Fatalf("Error creating ClickHouse client for server %s: %v", server.Name, err)
		}
		clients[server.Name] = client
		defer client.Close()
	}

	reg := prometheus.NewPedanticRegistry()

	// Создаем экспортер
	exp := exporter.NewExporter(cfg.Queries, clients, reg)

	// Регистрируем экспортер в Prometheus

	err = reg.Register(exp)
	if err != nil {
		panic(err)
	}

	// Запускаем HTTP-сервер для экспорта метрик
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting ClickHouse Prometheus Exporter on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
