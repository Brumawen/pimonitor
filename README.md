# pimonitor

Microservice that monitors the Raspberry Pi it runs on and reports the metrics over MQTT.

Build with Go 1.23 or later, from the `src` directory. The Pi does not need Go installed if you copy a finished binary onto it.

## Build

On the Pi, the default build matches that machine:

```sh
cd src
go build -o pimonitor .
```

From another machine, set the target to the Pi's CPU. On the Pi, `uname -m` reports it:

| `uname -m` | Pi OS | Build |
| --- | --- | --- |
| `aarch64` | 64-bit (Pi 3, 4, 5, Zero 2 W) | `GOARCH=arm64` |
| `armv7l` | 32-bit (Pi 2 and later) | `GOARCH=arm GOARM=7` |
| `armv6l` | 32-bit (Pi 1, Pi Zero) | `GOARCH=arm GOARM=6` |

A binary built for the wrong architecture will not start. `CGO_ENABLED=0` keeps the build from needing a C cross-compiler.

64-bit example:

```sh
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o pimonitor .
```

32-bit Pi 2 or later:

```sh
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o pimonitor .
```

32-bit Pi 1 or Pi Zero:

```sh
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o pimonitor .
```

Copy the binary to the Pi, for example:

```sh
scp pimonitor pi@raspberrypi:~/pimonitor/
```

## Configure

Put `config.json` in the same directory as the binary. The service switches to that directory on startup and reads the file from there. `period` is the upload interval in minutes. If it is missing or `0`, the service uses 5. `mqttHost` is a broker URL.

```json
{
  "period": 5,
  "enableMqtt": true,
  "mqttHost": "tcp://192.168.1.10:1883",
  "mqttUsername": "user",
  "mqttPassword": "password"
}
```

## Run as a service

Install and start from the directory that holds the binary and `config.json`. The installed unit keeps running that copy of the binary, so install it from the directory where it will stay. These commands need root. The service name is `PiMonitor`, and it listens on port 20517.

```sh
sudo ./pimonitor -service install
sudo ./pimonitor -service start
```

Other actions are `stop`, `restart`, and `uninstall`:

```sh
sudo ./pimonitor -service stop
sudo ./pimonitor -service restart
sudo ./pimonitor -service uninstall
```

Logs from the last hour:

```sh
journalctl --no-pager -u PiMonitor -S "1 hour ago"
```

To run it in the foreground instead of as a service, start `./pimonitor` with no `-service` flag. `-p` sets the port for that foreground run only. The installed service always uses 20517.
