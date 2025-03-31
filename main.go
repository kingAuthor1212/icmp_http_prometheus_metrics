package main

import (
	"strconv"

	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    /* icmpPingSuccess is a gauge that indicates whether the ICMP ping was successful.
    It records a value of 1 if the ping was successful else 0.
    https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#NewGaugeVec
    */
    icmpPingSuccess = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "icmp_ping_success",
            Help: "1 if the ping was successful else 0.",
        },
        []string{"target"},
    )
     /* icmpPingResponseTime measures the response time of the ICMP ping in seconds.
     This gauge reflects the time taken for the ping response.*/
    icmpPingResponseTime = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "icmp_ping_response_time",
            Help: "Response time of ICMP ping in seconds.",
        },
        []string{"target"},
    )

    logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

)

// init initializes the metrics for Prometheus monitoring and registers them
func init() {
    prometheus.MustRegister(icmpPingSuccess)
    prometheus.MustRegister(icmpPingResponseTime)
}


// check if vaild port
func checkPort(portStr string) bool {
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return false
    }
    return port >= 0 && port <= 65535
}

func main() {
    route := http.NewServeMux()
    target := "8.8.8.8"
    PORT := "8080"

    if len(os.Args) > 1 {
        target = os.Args[1]
        if len(os.Args)>2 && checkPort(os.Args[2]){PORT = os.Args[2]}
    }

    go func() {
        for {
            ping(target)// Send ICMP ping to the target
            httpGet("https://www.google.com") // Perform HTTP GET request to Google
            time.Sleep(2 * time.Second)  // Wait for 2 seconds before the next iteration
        }
        }()
        
    route.Handle("/metrics", promhttp.Handler())
    
    server := &http.Server{
        Addr: ":"+PORT,
        Handler: route,
    }
    logger.Info("Starting server on :"+PORT)
    err := server.ListenAndServe()
    logger.Error(err.Error())
    os.Exit(1)
}
