# prometheus-aranet4-exporter

This is an Aranet4 exporter for Prometheus. It requires an Aranet4 Home device.

It will export the metric `aranet4_measurement` with various fields:
* name (device name as shown in the Bluetooth LE scan)
* version (device's firmware version)
* humidity (percentage)
* pressure (in the unit configured via mobile app)
* temperature (in the unit configured via mobile app)
* co2 (in parts per million, or ppm)
* battery (percentage)
* quality ("green" if CO2 is between 0 and 999 ppm, "yellow" if between 1000 and
  1399 ppm, "red" if above 1400 ppm)
* interval (a string representing the measurement interval on the device)
* time (time string of when the last measurement was taken)

## Run it

```
go build
./prometheus-aranet4-exporter --mac-address <your Aranet4 MAC address>
```
