# Speedtest Metrics Exporter

[![GitHub Release](https://img.shields.io/github/v/release/kutovoys/speedtest-exporter?style=flat&color=blue)](https://github.com/kutovoys/speedtest-exporter/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kutovoys/speedtest-exporter/build-publish.yml)](https://github.com/kutovoys/speedtest-exporter/actions/workflows/build-publish.yml)
[![DockerHub](https://img.shields.io/badge/DockerHub-kutovoys%2Fspeedtest--exporter-blue)](https://hub.docker.com/r/kutovoys/speedtest-exporter/)
[![GitHub License](https://img.shields.io/github/license/kutovoys/speedtest-exporter?color=greeen)](https://github.com/kutovoys/speedtest-exporter/blob/main/LICENSE)

Speedtest Metrics Exporter is a Prometheus exporter designed to collect and expose network performance metrics using speedtest.net servers. Its unique feature is the ability to test multiple servers in a single iteration, providing comprehensive network performance data across different geographical locations and service providers.

## Features

- **Multi-Server Testing**: Test multiple speedtest servers in a single iteration by specifying server IDs, enabling comprehensive network performance monitoring across different locations
- **Automatic Server Selection**: If no server IDs are specified, the exporter automatically selects the closest server for testing
- **Comprehensive Metrics**: Collects detailed metrics including download speed, upload speed, latency, jitter, and test duration
- **Server Information**: Provides rich server metadata including server name, sponsor/provider, and distance
- **Configurable Update Intervals**: Customize how frequently speed tests are performed
- **Optional BasicAuth Protection**: Secure your metrics endpoint with basic authentication
- **Instance Labeling**: Support for instance labels to distinguish between multiple exporter instances
- **Memory Optimization**: Implements garbage collection after each test to maintain optimal performance

## Metrics

The exporter provides the following metrics:

| Name                                 | Description                       |
| ------------------------------------ | --------------------------------- |
| `speedtest_download_bits_per_second` | Download speed in bits per second |
| `speedtest_upload_bits_per_second`   | Upload speed in bits per second   |
| `speedtest_latency`                  | Latency in milliseconds           |
| `speedtest_jitter`                   | Jitter in milliseconds            |
| `speedtest_test_duration`            | Test duration in seconds          |

Each metric includes the following labels:

- `server_id`: Speedtest server identifier
- `server_name`: Name of the speedtest server
- `server_sponsor`: Server sponsor/provider
- `server_distance`: Distance to server in kilometers
- `instance`: Custom instance label (optional)

## Configuration

The exporter can be configured using environment variables or command-line arguments:

| Environment Variable | Command-Line Argument | Required | Default                   | Description                                                                                                                   |
| -------------------- | --------------------- | -------- | ------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `METRICS_PORT`       | `--metrics-port`      | No       | `9090`                    | Port to expose metrics on                                                                                                     |
| `METRICS_PROTECTED`  | `--metrics-protected` | No       | `false`                   | Enable BasicAuth protection                                                                                                   |
| `METRICS_USERNAME`   | `--metrics-username`  | No       | `metricsUser`             | Username for BasicAuth                                                                                                        |
| `METRICS_PASSWORD`   | `--metrics-password`  | No       | `MetricsVeryHardPassword` | Password for BasicAuth                                                                                                        |
| `UPDATE_INTERVAL`    | `--update-interval`   | No       | `60`                      | Test interval in minutes                                                                                                      |
| `SERVER_IDS`         | `--server-ids`        | No       | `0`                       | Comma-separated list of speedtest server IDs. If not specified or set to 0, the closest server will be selected automatically |
| `INSTANCE`           | `--instance`          | No       | `""`                      | Instance label for metrics                                                                                                    |

## Usage

### Finding Server IDs

See [how to find Speedtest.net server IDs](https://www.dcmembers.com/skwire/how-to-find-a-speedtest-net-server-id/) for instructions on locating specific server IDs.

### CLI

```bash
# With specific server IDs
speedtest-exporter --server-ids=1234,5678 --update-interval=60 --metrics-port=9090

# With basic auth protection
speedtest-exporter --metrics-protected --metrics-username=custom_user --metrics-password=custom_pass

# With default server selection (closest server will be chosen automatically)
speedtest-exporter
```

### Docker

```bash
docker run -d \
  -e SERVER_IDS=1234,5678,9012 \
  -e UPDATE_INTERVAL=60 \
  -p 9090:9090 \
  kutovoys/speedtest-exporter
```

### Docker Compose

```yaml
services:
  speedtest-exporter:
    image: kutovoys/speedtest-exporter
    environment:
      - SERVER_IDS=1234,5678,9012
      - UPDATE_INTERVAL=60
      - METRICS_PROTECTED=true
      - METRICS_USERNAME=custom_user
      - METRICS_PASSWORD=custom_password
    ports:
      - "9090:9090"
```

### Prometheus Configuration

Add the following to your prometheus.yml:

```yaml
scrape_configs:
  - job_name: "speedtest"
    static_configs:
      - targets: ["localhost:9090"]
    scrape_interval: 60m
```

### Multi-Server Testing Example

To test multiple servers, specify their IDs in the SERVER_IDS environment variable or --server-ids argument:

```bash
docker run -d \
  -e SERVER_IDS=1234,5678,9012 \
  -p 9090:9090 \
  kutovoys/speedtest-exporter
```

This will:

1. Test server ID 1234 first
2. After completion, test server ID 5678
3. Finally, test server ID 9012
4. Wait for the configured update interval before starting the next iteration

The exporter maintains separate metrics for each server, allowing you to:

- Compare performance across different providers
- Monitor network quality to different geographical locations
- Identify regional network issues
- Track performance trends per server

## API Endpoints

- `/metrics` - Prometheus metrics endpoint
- `/health` - Health check endpoint
- `/` - Basic information page

## Contribute

Contributions to Speedtest Metrics Exporter are warmly welcomed. Whether it's bug fixes, new features, or documentation improvements, your input helps make this project better. Here's a quick guide to contributing:

1. **Fork & Branch**: Fork this repository and create a branch for your work.
2. **Implement Changes**: Work on your feature or fix, keeping code clean and well-documented.
3. **Test**: Ensure your changes maintain or improve current functionality, adding tests for new features.
4. **Commit & PR**: Commit your changes with clear messages, then open a pull request detailing your work.
5. **Feedback**: Be prepared to engage with feedback and further refine your contribution.

Happy contributing! If you're new to this, GitHub's guide on [Creating a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request) is an excellent resource.
