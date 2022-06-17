package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"sbinet.org/x/aranet4"
)

var (
	flagPath          = pflag.String("p", "/metrics", "HTTP path where to expose metrics to")
	flagListen        = pflag.StringP("listen-address", "l", ":9101", "Address to listen to")
	flagMacAddress    = pflag.StringP("mac-address", "m", "", "Path to speedtest-cli")
	flagSleepInterval = pflag.DurationP("interval", "i", 30*time.Minute, "Interval between speedtest executions, expressed as a Go duration string")
)

func aranet4measurement(mac string) (string, string, *aranet4.Data, error) {
	dev, err := aranet4.New(mac)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create Aranet4 client: %w", err)
	}
	defer dev.Close()
	name, err := dev.Name()
	if err != nil {
		return "", "", nil, fmt.Errorf("cannot get device name: %w", err)
	}
	version, err := dev.Version()
	if err != nil {
		return "", "", nil, fmt.Errorf("cannot get device version: %w", err)
	}
	data, err := dev.Read()
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to read measurement: %w", err)
	}
	return name, version, &data, nil
}

// NewAranet4Collector returns a new Aranet4Collector object.
func NewAranet4Collector(macString string) (*Aranet4Collector, error) {
	mac, err := net.ParseMAC(macString)
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address: %w", err)
	}
	return &Aranet4Collector{
		mac: mac,
	}, nil
}

// Aranet4Collector is a prometheus collector for Aranet4 measurements.
type Aranet4Collector struct {
	mac net.HardwareAddr
}

var measurementDesc = prometheus.NewDesc(
	"aranet4_measurement",
	"Aranet4 CO2 measurement",
	[]string{"name", "version", "humidity", "pressure", "temperature", "co2", "battery", "quality", "interval", "time"},
	nil,
)

// Describe implements prometheus.Collector.Describe for WeatherCollector.
func (ac *Aranet4Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(ac, ch)
}

// Collect implements prometheus.Collector.Collect for WeatherCollector.
func (ac *Aranet4Collector) Collect(ch chan<- prometheus.Metric) {
	log.Printf("Reading measurement from Aranet4 device with MAC '%s'", ac.mac)
	name, version, data, err := aranet4measurement(*flagMacAddress)
	if err != nil {
		// if it fails, skip
		log.Printf("Failed to read measurement from Aranet4 device with mac '%s': %v", ac.mac, err)
	} else {
		// update values
		ch <- prometheus.MustNewConstMetric(
			measurementDesc,
			prometheus.GaugeValue,
			1,
			name,
			version,
			strconv.FormatFloat(data.H, 'f', -1, 64),
			strconv.FormatFloat(data.P, 'f', -1, 64),
			strconv.FormatFloat(data.T, 'f', -1, 64),
			strconv.FormatInt(int64(data.CO2), 10),
			strconv.FormatInt(int64(data.Battery), 10),
			string(data.Quality),
			data.Interval.String(),
			data.Time.String(),
		)
	}
}

func main() {
	pflag.Parse()

	wc, err := NewAranet4Collector(*flagMacAddress)
	if err != nil {
		log.Fatalf("Failed to create Aranet4 collector: %v", err)
	}
	if err := prometheus.Register(wc); err != nil {
		log.Fatalf("Failed to register Aranet4 collector: %v", err)
	}

	http.Handle(*flagPath, promhttp.Handler())
	log.Printf("Starting server on %s", *flagListen)
	log.Fatal(http.ListenAndServe(*flagListen, nil))
}
