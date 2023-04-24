package common

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

func PrometheusBoot(port int) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
		if err != nil {
			log.Fatalln("启动Prometheus失败")
			return
		}
		log.Println("启动Prometheus成功: " + strconv.Itoa(port))
	}()
}
