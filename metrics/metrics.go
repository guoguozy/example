package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

var (
	//Cpu percent usage
	Cpupercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "Cpupercent",
		Help:      "Cpu percent used.",
	})
	//system memory usage
	Diskpercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "Diskpercent",
		Help:      "Disk memory used.",
	})	
	//System memory usage
	Memorypercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "Memorypercent",
		Help:      "System memory used.",
	})
	//system time(minutes)
	Time = prometheus.NewGauge(
		prometheus.GaugeOpts{
		Name:      "time",
		Help:      "system time(minutes).",
	})

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "request_total",
			Help:      "Number of request processed by this service.",
		}, []string{},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "request_latency_seconds",
			Help:      "Time spent in this service.",
			Buckets:   []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 60.0, 120.0, 300.0},
		}, []string{},
	)
)

// AdmissionLatency measures latency / execution time of Admission Control execution
// usual usage pattern is: timer := NewAdmissionLatency() ; compute ; timer.Observe()
type RequestLatency struct {
	histo *prometheus.HistogramVec
	start time.Time
}

func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(Cpupercent)
	prometheus.MustRegister(Diskpercent)
	prometheus.MustRegister(Memorypercent)
	prometheus.MustRegister(Time)
}


// NewAdmissionLatency provides a timer for admission latency; call Observe() on it to measure
func NewAdmissionLatency() *RequestLatency {
	return &RequestLatency{
		histo: requestLatency,
		start: time.Now(),
	}
}

// Observe measures the execution time from when the AdmissionLatency was created
func (t *RequestLatency) Observe() {
	(*t.histo).WithLabelValues().Observe(time.Now().Sub(t.start).Seconds())
}


// RequestIncrease increases the counter of request handled by this service
func RequestIncrease() {
	requestCount.WithLabelValues().Add(1)
	
	hour:=time.Now().Hour()
	minute:=time.Now().Minute()
	//transfer 
	minutes:=float64(minute)/60
	hours:=float64(hour)
	Time.Set(hours*60+minutes)

	mem_,_ :=mem.VirtualMemory()
	Memorypercent.Set(mem_.UsedPercent)
	
	percent, _:= cpu.Percent(time.Second, false)
	Cpupercent.Set(percent[0])

	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	Diskpercent.Set(diskInfo.UsedPercent)
	
}
