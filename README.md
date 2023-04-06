# prometheus-aranet4-exporter

This is an Aranet4 exporter for Prometheus. It requires an Aranet4 Home device.

It exports the following metrics:
* `aranet4_co2`
* `aranet4_temperature`
* `aranet4_pressure`
* `aranet4_humidity`
* `aranet4_battery`

## Run it

First run: pair it with `bluetoothctl` (or similar method):

```
# bluetoothctl
[bluetooth]# power on
Changing power on succeeded
[bluetooth]# scan on
Discovery started
... (list of discovered devices)
[bluetooth]# pair <mac address>
... (insert PIN shown on device)
[bluetooth]# trust <mac address>
...
```

then build and run it:

```
go build
./prometheus-aranet4-exporter --mac-address <your Aranet4 MAC address>
```
