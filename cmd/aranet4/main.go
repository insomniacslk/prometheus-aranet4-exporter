package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
	"sbinet.org/x/aranet4"
)

var (
	flagMacAddress = pflag.StringP("mac-address", "m", "", "Aranet4 device MAC address")
)

type deviceInfo struct {
	Name    string
	Version string
}

func info(dev *aranet4.Device) (string, string, error) {
	name, err := dev.Name()
	if err != nil {
		return "", "", fmt.Errorf("cannot get device name: %w", err)
	}
	version, err := dev.Version()
	if err != nil {
		return "", "", fmt.Errorf("cannot get device version: %w", err)
	}
	return name, version, nil
}

func main() {
	pflag.Parse()
	dev, err := aranet4.New(*flagMacAddress)
	if err != nil {
		log.Printf("Failed to create Aranet4 client: %v", err)
		return
	}
	defer dev.Close()
	name, version, err := info(dev)
	if err != nil {
		log.Printf("Warning: failed to get device info: %v", err)
		return
	}
	fmt.Printf("Name:        %s\n", name)
	fmt.Printf("Version:     %s\n", version)

	start := time.Now()
	data, err := dev.Read()
	duration := time.Since(start)
	if err != nil {
		log.Printf("Read failed: %v", err)
		return
	}
	fmt.Println(data)
	fmt.Printf("done in %s\n", duration)
}
