package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"sbinet.org/x/aranet4"
)

var (
	flagMacAddress = pflag.StringP("mac-address", "m", "", "Aranet4 device MAC address")
)

func info(dev *aranet4.Device) (string, string, error) {
	name := dev.Name()
	version, err := dev.Version()
	if err != nil {
		return "", "", fmt.Errorf("cannot get device version: %w", err)
	}
	return name, version, nil
}

func main() {
	pflag.Parse()
	if *flagMacAddress == "" {
		log.Fatalf("Missing MAC address, see -m/--mac-address")
	}
	ctx := context.Background()
	dev, err := aranet4.New(ctx, strings.ToUpper(*flagMacAddress))
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
