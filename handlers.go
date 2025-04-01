package main

import(
	"log"
	"net"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"strconv"
	"time"
	"net/http"
	"os"
	"log/slog"
	"github.com/prometheus/client_golang/prometheus"
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

/* ping sends an ICMP echo request (ping) to the specified target IP address or hostname.
It measures the time taken for the request and updates Prometheus metrics accordingly.

Parameters:
- target: A string representing the target IP address or hostname to ping.
ICMP packet handling, https://pkg.go.dev/golang.org/x/net/icmp
package to resolve IP addresses, https://pkg.go.dev/net#ResolveIPAddr
package for measuring durations, https://pkg.go.dev/time
*/
func ping(target string) {
    c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")  
    if err != nil {
        log.Println("Error listening:", err)
        return
    }
    defer c.Close()

    start := time.Now()

    message := icmp.Message{
        Type: ipv4.ICMPTypeEcho,
        Code: 0,
        Body: &icmp.Echo{
            ID:   os.Getpid() & 0xffff,
            Seq:  1,
            Data: []byte("Ping"),
        },
    }

    msgBytes, err := message.Marshal(nil)
    if err != nil {
        log.Println("Error marshaling message:", err)
        return
    }

    addr, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        log.Println("Error resolving address:", err)
        return
    }

    _, err = c.WriteTo(msgBytes, addr)
    if err != nil {
        log.Println("Error writing:", err)
        icmpPingSuccess.WithLabelValues(target).Set(0)
        return
    }

    reply := make([]byte, 1024)
    c.SetDeadline(time.Now().Add(2 * time.Second))
    _, _, err = c.ReadFrom(reply)
    duration := time.Since(start).Seconds()

    if err == nil {
        icmpPingSuccess.WithLabelValues(target).Set(1)
        icmpPingResponseTime.WithLabelValues(target).Set(duration)
        log.Println("Ping to",target ,"in", duration, "seconds")
    } else {
        icmpPingSuccess.WithLabelValues(target).Set(0)
        log.Println("Ping failed", err)
    }
}

/* httpGet performs an HTTP GET request to the specified URL.
It measures the time taken for the request and prints the response time.

Parameters:
- url: A string representing the target URL for the GET request.
Go net/http package https://pkg.go.dev/net/http
*/
func httpGet(url string) {
    startTime := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        log.Println("HTTP GET request failed:", err)
        return
    }
    defer resp.Body.Close()

    duration := time.Since(startTime).Seconds()
    log.Println("HTTP GET response from",url,"in", duration, "seconds")
}

// check if vaild port
func checkPort(portStr string) bool {
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return false
    }
    return port >= 0 && port <= 65535
}