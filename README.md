# ICMP Ping and HTTP GET with Prometheus Metrics

This Go application sends ICMP pings to a specified target and performs HTTP GET requests to Google. It records metrics for both operations and exposes them via a Prometheus HTTP endpoint.

## Features

- ICMP Ping to a specified target.
- HTTP GET request to `https://www.google.com`.
- Exposes metrics for ICMP ping success and response time.
- Metrics are available at the `/metrics` endpoint.

## Prerequisites

- Go installed (version 1.11 or higher).

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/kingAuthor1212/icmp_http_prometheus_metrics.git
   
   cd icmp_http_prometheus_metrics
   ```

2. **Install dependencies:**

   Run the following command to get the necessary Go modules and libraries:

   ```bash
   go get golang.org/x/net/icmp
   go get github.com/prometheus/client_golang/prometheus
   go mod tidy
   ```

## Build

To build the application, run:

```bash
go build -o icmp_http_prometheus_metrics
```

This will create an executable named `icmp_http_prometheus_metrics`.

## Run

To run the application, use the following command:

```bash
./icmp_http_prometheus_metrics  <target> <port> 
```

Replace target and port with `hostname or IP address and port nunmber you want` and If no target and port are provided, it defaults to `8.8.8.8` and `8080`.

Example:

```bash
./icmp_http_prometheus_metrics 1.1.1.1 9090
```

The application will start an HTTP server on port provided else `8080` and expose the metrics at `/metrics`.

## Test

To test the application, follow these steps:

1. **Run the application** as described above.
2. Open your web browser or use a tool like `curl` to access the metrics endpoint:

   ```bash
   curl http://localhost:<port>/metrics
   ```
   ```browser
   http://localhost:<port>/metrics
   ```

You should see metrics related to ICMP ping success and response times.

## Stopping the Application

You can stop the application using `Ctrl + C` in the terminal where it's running.
