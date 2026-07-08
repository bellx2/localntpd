# localntpd

[日本語版はこちら / Japanese version](README_ja.md)

A simple NTP server that serves your local PC's clock. Written in Go, it runs on
multiple platforms (Windows / macOS / Linux) and can be registered as a system
service.

## Features

- Serves the local machine's time over the NTP protocol (UDP)
- Cross-platform single binary (Windows / macOS / Linux)
- Installable as an OS system service (via [kardianos/service](https://github.com/kardianos/service))
- No upstream synchronization — acts as a standalone local clock (stratum, reference ID `LOCL`)

## Build

```bash
go build -o localntpd .
```

### Windows build (with icon and version info)

To embed an icon and version metadata into the `.exe`, generate the Windows
resource with [`go-winres`](https://github.com/tc-hib/go-winres) before building.
`go build` links the generated `rsrc_windows_*.syso` automatically.

```bash
go install github.com/tc-hib/go-winres@v0.3.3   # first time only
go-winres make --arch amd64,arm64                # generate .syso from winres/
GOOS=windows GOARCH=amd64 go build -o localntpd.exe .
```

Requires Go 1.26+. The resource definition lives in `winres/winres.json` and
the icons in `winres/icon16.png` / `icon32.png` / `icon48.png` / `icon256.png`
(source: `winres/icon.svg`). Replace them and re-run `go-winres make` to
customize. The `.syso` files are only linked on Windows builds, so they do not
affect macOS/Linux builds.

## Usage

```bash
localntpd [command] [options]
```

### Commands

| Command     | Description                          |
|-------------|--------------------------------------|
| `run`       | Run in the foreground (default)      |
| `install`   | Register as a system service         |
| `uninstall` | Remove the service                   |
| `start`     | Start the service                    |
| `stop`      | Stop the service                     |
| `restart`   | Restart the service                  |
| `status`    | Show the service status              |
| `help`      | Show help                            |

### Options

| Option           | Default | Description                          |
|------------------|---------|--------------------------------------|
| `-addr string`   | `:123`  | Listen address (e.g. `:123`, `0.0.0.0:123`) |
| `-stratum uint`  | `2`     | Stratum (1–15)                       |

### Examples

```bash
localntpd run -addr :12345              # Run on a non-privileged port
localntpd install                       # Register as a service (requires administrator/root)
localntpd install -addr :12345          # Register with a custom listen address
localntpd start
```

## Notes

- Binding to port 123 (the standard NTP port) requires administrator/root
  privileges. This applies to `install` and to running with the default `:123`.
- For development and testing, use a non-privileged port such as `-addr :12345`.
- Options passed to `install` (`-addr`, `-stratum`) are stored as service
  startup arguments.

## License

See the repository for license information.
