# logpipe

> Streaming log aggregator with filtering and forwarding rules for containerized environments

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Installation

```bash
go install github.com/yourorg/logpipe@latest
```

Or pull the Docker image:

```bash
docker pull yourorg/logpipe:latest
```

---

## Usage

Define your forwarding rules in a `logpipe.yaml` config file:

```yaml
inputs:
  - type: docker
    containers: ["api-*", "worker-*"]

filters:
  - level: error
  - exclude: "healthcheck"

outputs:
  - type: stdout
  - type: http
    url: https://logs.example.com/ingest
    headers:
      Authorization: Bearer $LOG_TOKEN
```

Then run:

```bash
logpipe --config logpipe.yaml
```

Logs matching your filter rules are streamed and forwarded in real time. Use `--dry-run` to validate your config without forwarding.

```bash
logpipe --config logpipe.yaml --dry-run
```

---

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `LOG_TOKEN` | Bearer token used for HTTP output authorization | *(none)* |
| `LOGPIPE_CONFIG` | Path to config file (overridden by `--config` flag) | `logpipe.yaml` |
| `LOGPIPE_LOG_LEVEL` | Internal log verbosity (`debug`, `info`, `warn`, `error`) | `info` |

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the [MIT License](LICENSE).
