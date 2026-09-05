# TELUS T3200 WiFi Toggle

A small Go utility that uses [chromedp](https://github.com/chromedp/chromedp) to automate the web interface of the **TELUS T3200** router and enable or disable its WiFi radio.

The original purpose of this project was to experiment with browser automation using Go and `chromedp`, but it is also useful as a lightweight way to automate router configuration from an embedded Linux device.

<img width="1456" height="720" alt="telus_router_page" src="https://github.com/user-attachments/assets/d71ccb84-68db-4587-aa16-eec37a5c4950" />

## Why Go + chromedp?

Many small embedded devices have limited CPU and memory resources. Running browser automation through Python, Node.js, or similar scripting environments can require a relatively large runtime and dependency footprint.

This project instead uses Go:

* Small standalone executable
* No Python or Node.js runtime required
* Easy cross-compilation for ARM devices
* `chromedp` communicates with Chrome/Chromium through the Chrome DevTools Protocol
* Suitable for automation from cron, systemd, scripts, or other embedded applications

`chromedp` is specifically designed to drive Chrome/Chromium through the DevTools Protocol and supports headless operation. See the [chromedp documentation](https://github.com/chromedp/chromedp) for more information.

> **Important:** The executable itself is small, but it still requires a compatible Chrome/Chromium installation on the target system. The browser is the main runtime dependency.

## Supported Router

This utility was developed specifically for the:

**TELUS T3200**

It interacts with the router's web administration interface and is **not intended to be a generic router WiFi controller**.

The current implementation expects the router to be reachable at:

```text
http://10.25.73.1
```

The router's web UI is automated using DOM selectors specific to this firmware. Firmware updates or changes to the router's web interface may therefore break the program.

## Requirements

### Target device

* Linux
* ARMv7 (`armhf`) or ARM64 (`aarch64`) supported
* Network access to the T3200 router
* Chrome or Chromium
* Router administrator password

### Development machine

* Go 1.16 or newer
* Chrome/Chromium for actually running the automation

The current `go.mod` specifies Go 1.16.

## Building

Clone the repository:

```bash
git clone https://github.com/reganchan/telus_t3200_wifi_toggle.git
cd telus_t3200_wifi_toggle
```

Download dependencies:

```bash
go mod download
```

Build:

```bash
go build -o wifi .
```

## Usage

The program requires the router administrator password.

### Check current WiFi status

```bash
./wifi --pass='YOUR_PASSWORD'
```

The program will log whether the WiFi radio is currently enabled.

### Enable WiFi

```bash
./wifi --pass='YOUR_PASSWORD' enable
```

### Disable WiFi

```bash
./wifi --pass='YOUR_PASSWORD' disable
```

The program logs its status as JSON using `logrus`, which makes the output convenient for consumption by scripts or logging systems.

## Example

```text
$ ./wifi --pass='secret' disable

{"enable":false,"level":"info","msg":"changing wifi status"}
{"enabled":false,"level":"info","msg":"wifi enable status"}
```

Running without `enable` or `disable` simply logs the current state.

## Chrome / Chromium

`chromedp` normally runs Chrome in headless mode, which makes it suitable for devices without a graphical desktop.

For example, on Debian/Ubuntu:

```bash
sudo apt install chromium
```

Depending on the distribution, the executable may instead be named `chromium-browser` or `google-chrome`.

Verify that a browser is available:

```bash
chromium --version
```

If `chromedp` cannot find the browser automatically, the code can be configured to explicitly specify the Chrome/Chromium executable.

### Embedded devices

For very small systems, a lightweight Chromium/headless-shell installation may be preferable. `chromedp` also documents using its headless-shell image for headless environments.

Keep in mind that **Go eliminates the Python/Node.js runtime dependency, but it does not eliminate the need for a browser** because this program is automating the router's web UI.

## ARM builds

Pre-built binaries can be cross-compiled from a normal x86-64 development machine.

### ARM64

For 64-bit ARM Linux devices:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
    go build -o wifi-linux-arm64 .
```

### ARMv7

For 32-bit ARMv7 Linux devices:

```bash
GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
    go build -o wifi-linux-armv7 .
```

Go officially supports Linux ARMv5/v6/v7 and ARM64, with `GOARM=7` used for ARMv7 and `GOARCH=arm64` for ARM64.

## GitHub Actions

The repository can automatically build Linux ARM artifacts whenever code is pushed.

Recommended artifacts:

| Target | Go settings                     | Artifact           |
| ------ | ------------------------------- | ------------------ |
| ARMv7  | `GOOS=linux GOARCH=arm GOARM=7` | `wifi-linux-armv7` |
| ARM64  | `GOOS=linux GOARCH=arm64`       | `wifi-linux-arm64` |
| x86-64 | `GOOS=linux GOARCH=amd64`       | `wifi-linux-amd64` |

A GitHub Actions workflow can upload these binaries as build artifacts so they can be downloaded directly from the Actions run.

For releases, the same binaries can also be attached automatically to a GitHub Release.

## How it works

The program performs roughly the following sequence:

1. Start a headless Chrome/Chromium instance through `chromedp`.
2. Connect to the T3200 web interface.
3. Log in using the supplied administrator password.
4. Navigate to the wireless settings page.
5. Inspect the current WiFi radio state.
6. Optionally click the enable/disable radio button.
7. Apply the new configuration.
8. Read the resulting WiFi state.
9. Log the result as JSON.

The browser automation is implemented directly against the T3200's HTML/DOM elements. The current source uses selectors such as `admin_password`, `btn_login`, `id_wl_on`, `id_wl_off`, and `btn_apply`.

## Security considerations

The router administrator password is supplied as a command-line argument:

```bash
./wifi --pass='YOUR_PASSWORD'
```

Be aware that command-line arguments can potentially be visible to other users through process inspection or shell history.

For unattended automation, consider supplying the password through a mechanism appropriate for your environment rather than putting it directly into a command line.

For example, a wrapper script could retrieve it from a protected configuration file and invoke the program.

## Limitations

* Designed specifically for the TELUS T3200.
* Depends on the router's web UI and its DOM structure.
* Router firmware changes may break the selectors.
* Requires Chrome/Chromium on the target device.
* The program currently assumes the router is accessible at `10.25.73.1`.
* This is browser automation rather than an official router API.

## Project Status

This project started primarily as an experiment with **Go + chromedp browser automation**.

It also serves as a practical example of using Go to replace heavier scripting environments on resource-constrained Linux/ARM devices.

If the T3200 exposes no suitable API for a particular automation task, browser automation can sometimes provide a practical alternative.

## License

See the repository for licensing information.
