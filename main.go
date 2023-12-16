package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"sbinet.org/x/aranet4"
)

var (
	flagPath       = pflag.String("p", "/metrics", "HTTP path where to expose metrics to")
	flagListen     = pflag.StringP("listen-address", "l", ":9111", "Address to listen to")
	flagMacAddress = pflag.StringP("mac-address", "m", "", "Aranet4 Home devices's MAC address")
	flagInterval   = pflag.DurationP("interval", "i", 1*time.Minute, "Interval between sensor readings, expressed as a Go duration string")
)

func makeGauge(name, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aranet4_" + name,
			Help: "Aranet4 - " + help,
		},
		[]string{"name", "version", "interval"},
	)
}

var (
	humidityGauge    = makeGauge("humidity", "humidity (percentage)")
	pressureGauge    = makeGauge("pressure", "pressure (in the unit configured via Aranet4 app)")
	temperatureGauge = makeGauge("temperature", "temperature (in the unit configured via Aranet4 app)")
	co2Gauge         = makeGauge("co2", "CO2 (parts per million)")
	batteryGauge     = makeGauge("battery", "battery (percentage)")
)

func collector(mac net.HardwareAddr) {
	log.Printf("Collector started")
	ctx := context.Background()
	for {
		dev, err := aranet4.New(ctx, strings.ToUpper(mac.String()))
		if err != nil {
			log.Printf("Failed to connect to Aranet4 device: %v", err)
			time.Sleep(*flagInterval)
			continue
		}
		name := dev.Name()
		log.Printf("Name: %s", name)
		version, err := dev.Version()
		if err != nil {
			dev.Close()
			log.Printf("Failed to get Aranet4 device version: %v", err)
			time.Sleep(*flagInterval)
			continue
		}
		log.Printf("Version: %v", version)
		data, err := dev.Read()
		if err != nil {
			dev.Close()
			log.Printf("Failed to read Aranet4 device: %v", err)
			time.Sleep(*flagInterval)
			continue
		}
		dev.Close()
		log.Printf("Aranet4 reading for device '%s': %+v", name, data)
		humidityGauge.WithLabelValues(name, version, data.Interval.String()).Set(float64(data.H))
		pressureGauge.WithLabelValues(name, version, data.Interval.String()).Set(float64(data.P))
		temperatureGauge.WithLabelValues(name, version, data.Interval.String()).Set(float64(data.T))
		co2Gauge.WithLabelValues(name, version, data.Interval.String()).Set(float64(data.CO2))
		batteryGauge.WithLabelValues(name, version, data.Interval.String()).Set(float64(data.Battery))

		time.Sleep(*flagInterval)
	}
}

func main() {
	pflag.Parse()

	mac, err := net.ParseMAC(*flagMacAddress)
	if err != nil {
		log.Fatalf("Invalid MAC address: %v", err)
	}

	// register all gauges
	if err := prometheus.Register(humidityGauge); err != nil {
		log.Fatalf("Failed to register Aranet4 humidity gauge: %v", err)
	}
	if err := prometheus.Register(pressureGauge); err != nil {
		log.Fatalf("Failed to register Aranet4 pressure gauge: %v", err)
	}
	if err := prometheus.Register(temperatureGauge); err != nil {
		log.Fatalf("Failed to register Aranet4 temperature gauge: %v", err)
	}
	if err := prometheus.Register(co2Gauge); err != nil {
		log.Fatalf("Failed to register Aranet4 co2 gauge: %v", err)
	}
	if err := prometheus.Register(batteryGauge); err != nil {
		log.Fatalf("Failed to register Aranet4 battery gauge: %v", err)
	}

	// start collector
	go collector(mac)

	server := http.Server{Addr: *flagListen}
	http.Handle(*flagPath, promhttp.Handler())
	log.Printf("Starting server on %s", *flagListen)
	log.Fatal(server.ListenAndServe())
}
