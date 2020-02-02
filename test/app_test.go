package main

import (
	"fmt"
	"time"
	"regexp"
//	"reflect"
	"testing"
	"net/http"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)
var (
	exapp = "http://example-app:3000/metrics"
	exappname = "cadvisor-test_example-app_1"
        cadvisor = "http://cadvisor:8080/metrics"

	ignoreMetricRe  = regexp.MustCompile(`^container_.*$`)
	ignoreLabelRe = regexp.MustCompile(`^container_.*$|^image$|^id$|^name$`)
        metricNameRe = regexp.MustCompile(`^app_(.*)$`)
)

func extendMetric(m model.Metric, labels model.LabelSet) {
	for k, v := range labels {
		m[k] = v
	}
}

func fetchMetrics(t *testing.T, url string, ch chan *model.Sample) {
	defer close(ch)
	currentTime := time.Now()
        res, err := http.Get(url)
        if err != nil {
		fmt.Println(err)
                t.Fatal(err)
	}
	defer res.Body.Close()
       	if res.StatusCode != http.StatusOK {
		fmt.Println(err)
		t.Fatal("server returned HTTP status", res.Status)
	}
	sdec := expfmt.SampleDecoder{
		Dec: expfmt.NewDecoder(res.Body, expfmt.ResponseFormat(res.Header)),
		Opts: &expfmt.DecodeOptions{
			Timestamp: model.TimeFromUnixNano(currentTime.UnixNano()),
		},
	}
	decSamples := make(model.Vector, 0, 50)
	for {
//		fmt.Println("---------------", len(decSamples))
		if err = sdec.Decode(&decSamples); err != nil {
			break
		}
		for _, sample := range decSamples {
			ch <- sample
		}
	}
}
func fetchAppMetrics(t *testing.T, res chan model.Vector){
	defer close(res)
	ch := make(chan *model.Sample)
	resVector := make(model.Vector, 0, 50)
	go fetchMetrics(t, exapp, ch)
	for sample := range ch {
		resVector = append(resVector, sample)
	}
	res <- resVector
}
func fetchAppMetricsFromCadvisor(t *testing.T, res chan model.Vector, name string){
	defer close(res)
        ch := make(chan *model.Sample)
	resVector := make(model.Vector, 0, 50)
        go fetchMetrics(t, cadvisor, ch)
        for sample := range ch {
		metName := string(sample.Metric[model.MetricNameLabel])
		if !ignoreMetricRe.MatchString(metName) {
			if val, ok := sample.Metric["name"]; ok && string(val) == name {
                                s := model.Sample{}
                                s.Value = sample.Value
                                s.Timestamp = sample.Timestamp
                                s.Metric = make(model.Metric)
                                s.Metric[model.MetricNameLabel] = sample.Metric[model.MetricNameLabel]
				for labelName, labelValue  := range sample.Metric {
					if sm := metricNameRe.FindStringSubmatch(string(labelName)); len(sm) > 0 {
						s.Metric[model.LabelName(sm[1])] = labelValue
					}
				}
//                                fmt.Printf("--------------- %v\n", s)
				resVector = append(resVector, &s)
			}
		}
        }
	res <- resVector
}
func metricsToFingerprintSet(samples model.Vector) model.FingerprintSet {
	fset := make(model.FingerprintSet)
	for _, sample := range samples {
		fprint := sample.Metric.Fingerprint()
		fset[fprint] = struct{}{}
	}
	return fset
}
func testMetricLabelsEqual(t *testing.T) bool {
	cmc := make(chan model.Vector)
	amc := make(chan model.Vector)
	go fetchAppMetrics(t, amc)
	go fetchAppMetricsFromCadvisor(t, cmc, exappname)
	cm, am := <-cmc, <-amc
	cf := metricsToFingerprintSet(cm)
	af := metricsToFingerprintSet(am)
	return cf.Equal(af)
}

func TestAppMetrics(t *testing.T) {
	fmt.Println("start test")
	fails := 0
	attempts := 50
	for i := 0; i < attempts; i++ {
		status := testMetricLabelsEqual(t)
		if status {
			fmt.Print(".")
		} else {
			fails++
			fmt.Print("f")
		}
		time.Sleep(time.Second/4)
	}
	if fails > 0 {
		t.Fatalf("Cadvisor and App metrics not equal in %v of %v cases", fails, attempts)
	}
	fmt.Print("\n")
	fmt.Println("ok")
}
