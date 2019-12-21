package main

import (
	"flag"
	"net/http"
	"log"
	"os"
	"time"
	"math/rand"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	running = true
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func randomValues(metrics []prometheus.Gauge) {
	for {
		if !running {
			break
		}
		for _, metric := range metrics {
			metric.Set(rand.Float64())
		}
		time.Sleep(time.Millisecond * 500)
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())

	bind := ""
	randMetricCount := 0
	customLabel:= ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":3000", "The socket to bind to.")
	flagset.StringVar(&customLabel, "cl", "custom_label", "The custom label name.")
	flagset.IntVar(&randMetricCount, "c", 50, "The count of random metrics.")
	flagset.Parse(os.Args[1:])
	r := prometheus.NewRegistry()


	randomMetrics := make([]prometheus.Gauge, 0, 50)
	for i := 0; i < randMetricCount; i++ {
		rmetric := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "ex_app_rand_metric_" + randSeq(10),
				Help: "Example Random Metric",
				ConstLabels: map[string]string{
					"static_label": strconv.Itoa(rand.Intn(10)),
					"rand_label_" + randSeq(2): strconv.Itoa(rand.Intn(10)),
					"rand_label_" + randSeq(1): strconv.Itoa(rand.Intn(10)),
					customLabel: strconv.Itoa(rand.Intn(10)),
				},
			})
		randomMetrics = append(randomMetrics, rmetric)
		r.MustRegister(rmetric)
	}
	http.Handle("/", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	go randomValues(randomMetrics)
	log.Fatal(http.ListenAndServe(bind, nil))
	running = false
}
