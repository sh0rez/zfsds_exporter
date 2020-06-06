package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
)

const Namespace = "zfsds"

var Version = "dev"

func main() {
	var c Collector

	listen := pflag.String("listen", ":2112", "address:port to listen on")
	pflag.StringSliceVarP(&c.Props, "props", "p", []string{"used"}, "ZFS properties to query for each dataset")
	pflag.StringSliceVarP(&c.Datasets, "datasets", "d", nil, "ZFS datasets to include")
	showVersion := pflag.BoolP("version", "v", false, "Show version information")
	pflag.Parse()

	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if len(c.Datasets) == 0 {
		log.Fatalln("Please specify datasets with --datasets / -d")
	}
	if len(c.Props) == 0 {
		log.Fatalln("Please specify ZFS properties with --props / -p")
	}

	prometheus.MustRegister(c)
	log.Printf("Listening on %s", *listen)
	if err := http.ListenAndServe(*listen, promhttp.Handler()); err != nil {
		log.Fatalln(err)
	}
}

type Collector struct {
	Datasets []string
	Props    []string
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	for _, name := range c.Datasets {
		values, err := GetProps(name, c.Props)
		if err != nil {
			log.Println(err)
			continue
		}

		for k, v := range values {
			i, err := strconv.Atoi(v)
			if err != nil {
				log.Println(err)
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(Namespace+"_"+k, "", []string{"dataset"}, nil),
				prometheus.GaugeValue,
				float64(i),
				name,
			)
		}
	}
}
