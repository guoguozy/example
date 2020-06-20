# 云上的应⽤开发、部署和运维
---
修改业务逻辑代码，增加自定义的prometheus Exporter
## 代码逻辑
- 此处我定义了四个指标：Cpupercent，Diskpercent，Memorypercent，time 。分别是Cpu占用率，Disk占用率，内存占用率，系统时间（分钟）。
```go
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
    ······
    )
```
- import github.com/shirou/gopsutil
```
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)
```
- 注册指标
```
func Register() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestLatency)
	prometheus.MustRegister(Cpupercent)
	prometheus.MustRegister(Diskpercent)
	prometheus.MustRegister(Memorypercent)
	prometheus.MustRegister(Time)
}
```
- 赋值
使用github.com/shirou/gopsutil 的资源获取本机的Cpu等信息，注意其赋值方式以及**所取参数的类型**。
注意：代码中不可定义无用变量，否则会报错
```
hour:=time.Now().Hour()
minute:=time.Now().Minute()
//transfer 
minutes:=float64(minute)/60
hours:=float64(hour)
Time.Set(hours*60+minutes)

mem_,_ :=mem.VirtualMemory()
memorypercent.Set(mem_.UsedPercent)
	
percent, _:= cpu.Percent(time.Second, false)
Cpupercent.Set(percent[0])

parts, _ := disk.Partitions(true)
diskInfo, _ := disk.Usage(parts[0].Mountpoint)
Diskpercent.Set(diskInfo.UsedPercent)
```

## 文件结构
````
├── Dockerfile                   制作镜像所使用
├── README
├── deploy                       部署资源对象时使用的配置文件
│   ├── deployment.yaml          云服务的Deployment配置文件
│   ├── metrics_service.yaml     
│   ├── prometheus.config.yml    prometheus抓取目标配置文件
│   ├── prometheus.deploy.yml    prometheus部署所使用Deployment
│   ├── prometheus.rbac.yml      prometheus权限配置文件
│   └── service.yaml
├── go.mod                       依赖管理
├── go.sum                       依赖管理
├── metrics                      Exporter
│   └── metrics.go
├── metrics_version            
│   └── main.go
└── without_metrics            
    └── main.go
