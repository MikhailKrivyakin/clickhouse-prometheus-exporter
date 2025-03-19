package exporter

import (
	"clickhouse-prometheus-exporter/internal/clickhouse"
	"clickhouse-prometheus-exporter/internal/config"
	"log"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Query = config.Query

type Exporter struct {
	queries []Query
	clients map[string]*clickhouse.Client
	metrics map[string]*prometheus.GaugeVec
}

var (
	metrics map[string]*prometheus.GaugeVec
)

func init() {
	metrics = make(map[string]*prometheus.GaugeVec)
}

func NewExporter(queries []Query, clients map[string]*clickhouse.Client, reg *prometheus.Registry) *Exporter {
	metrics := make(map[string]*prometheus.GaugeVec)
	for _, q := range queries {
		metrics[q.MetricName] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: q.MetricName,
				Help: q.Help,
			},
			append([]string{"server"}, q.Labels...),
		)

	}
	return &Exporter{
		queries: queries,
		clients: clients,
		metrics: metrics,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	for serverName, client := range e.clients {
		wg.Add(1)
		go func(serverName string, client *clickhouse.Client) {
			defer wg.Done()
			e.collectMetricsForServer(serverName, client, ch)
		}(serverName, client)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()
	for name, _ := range e.metrics {
		e.metrics[name].Collect(ch)
	}
}

// collectMetricsForServer собирает метрики для одного сервера
func (e *Exporter) collectMetricsForServer(serverName string, client *clickhouse.Client, ch chan<- prometheus.Metric) {
	for _, q := range e.queries {
		results, err := client.Query(q.Query)
		if err != nil {
			log.Printf("Error executing query %s on server %s: %v", q.Name, serverName, err)
			continue
		}

		// Используем map для отслеживания уникальных комбинаций лейблов
		seenLabels := make(map[string]bool)

		for _, result := range results {
			// Извлекаем лейблы
			labelValues := []string{serverName}
			for _, label := range q.Labels {
				if value, ok := result[label].(string); ok {
					labelValues = append(labelValues, value)
				} else {
					log.Printf("Label %s not found in query result for server %s", label, serverName)
					labelValues = append(labelValues, "unknown") // Используем "unknown" для отсутствующих лейблов
				}
			}

			// Проверяем, была ли уже обработана эта комбинация лейблов
			labelKey := strings.Join(labelValues, ",")
			if seenLabels[labelKey] {
				log.Printf("Duplicate labels for metric '%s': %v", q.MetricName, labelValues)
				continue
			}
			seenLabels[labelKey] = true

			// Извлекаем значение метрики и приводим его к float64
			if value, ok := result[q.ValueColumn]; ok {
				var floatValue float64
				switch v := value.(type) {
				case float64:
					floatValue = v
				case int:
					floatValue = float64(v)
				case int64:
					floatValue = float64(v)
				case uint64:
					floatValue = float64(v)
				case float32:
					floatValue = float64(v)
				default:
					log.Printf("Unsupported type for value column '%s': %T", q.ValueColumn, value)
					continue
				}

				// Устанавливаем значение метрики
				e.metrics[q.MetricName].WithLabelValues(labelValues...).Set(floatValue)
			} else {
				log.Printf("Value column '%s' not found in query result for server %s", q.ValueColumn, serverName)
				// Устанавливаем значение по умолчанию 0, если столбец не найден
				e.metrics[q.MetricName].WithLabelValues(labelValues...).Set(0)
			}
		}
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.metrics {
		m.Describe(ch)
	}
}
