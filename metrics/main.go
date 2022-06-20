package main

import (
	"log"
	"os"
	"time"

	"github.com/rcrowley/go-metrics"
)

//for:github.com/rcrowley/go-metrics监控观测使用
type Metrics struct {
	metrics.Meter
	metrics.Histogram
}

func (me Metrics) Incr(size int64, time time.Duration) {
	me.Mark(size)
	me.Histogram.Update(int64(time.Microseconds()))
}

var PoMt Metrics

func init() {
	inertMeterMT := metrics.NewMeter()
	metrics.Register("insertSize bytes/s:", inertMeterMT)
	s := metrics.NewUniformSample(1024 * 4)
	inertTimeMT := metrics.NewHistogram(s)
	metrics.Register("insertLatency ms:", inertTimeMT)
	PoMt = Metrics{
		inertMeterMT,
		inertTimeMT,
	}
	go metrics.Log(metrics.DefaultRegistry,
		1*time.Second,
		log.New(os.Stdout, "metrics: ", log.Lmicroseconds))

}

func main() {
	var j int64
	j = 1
	for true {
		j++
		a := time.Now()
		time.Sleep(time.Millisecond)
		PoMt.Incr(j, time.Since(a))
	}
}

/*
metrics: 15:52:27.440643 meter insertSize bytes/s:
metrics: 15:52:27.440668   count:        32429430
metrics: 15:52:27.440675   1-min rate:    1896198.11
metrics: 15:52:27.440680   5-min rate:    1692838.57
metrics: 15:52:27.440683   15-min rate:   1657607.62
metrics: 15:52:27.440687   mean rate:     3242948.47
metrics: 15:52:27.441350 histogram insertLatency ms:
metrics: 15:52:27.441359   count:            8052
metrics: 15:52:27.441368   min:              1012
metrics: 15:52:27.441376   max:              4981
metrics: 15:52:27.441384   mean:             1240.26
metrics: 15:52:27.441400   stddev:            157.19
metrics: 15:52:27.441404   median:           1268.00
metrics: 15:52:27.441408   75%:              1301.00
metrics: 15:52:27.441412   95%:              1333.00
metrics: 15:52:27.441416   99%:              1382.00
metrics: 15:52:27.441420   99.9%:            3245.82
*/
