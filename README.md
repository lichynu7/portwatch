# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected listeners.

---

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a baseline of allowed ports:

```bash
portwatch --allow 22,80,443 --interval 30s
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--allow` | _(none)_ | Comma-separated list of expected open ports |
| `--interval` | `60s` | How often to scan for open ports |
| `--alert` | `stdout` | Alert destination (`stdout`, `syslog`, or webhook URL) |
| `--verbose` | `false` | Enable verbose logging |

**Example output when an unexpected port is detected:**

```
[ALERT] 2024-01-15T10:32:05Z - Unexpected listener detected on port 4444 (PID: 9821)
```

Run as a background service:

```bash
portwatch --allow 22,80,443 --alert syslog &
```

---

## How It Works

`portwatch` periodically reads active TCP/UDP listeners from the system, compares them against your allowed list, and fires an alert whenever an unlisted port appears. It is designed to be minimal, dependency-free, and easy to integrate into existing monitoring pipelines.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

MIT © 2024 youruser