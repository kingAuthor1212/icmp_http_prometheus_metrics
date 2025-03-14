package main

import (
    "fmt"
    "net"
    "net/http"
    "os"
    "time"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


var (
    /* icmpPingSuccess is a gauge that indicates whether the ICMP ping was successful.
    It records a value of 1 if the ping was successful else 0.
    https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#NewGaugeVec*/
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
        fmt.Println("Error listening:", err)
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
        fmt.Println("Error marshaling message:", err)
        return
    }

    addr, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }

    _, err = c.WriteTo(msgBytes, addr)
    if err != nil {
        fmt.Println("Error writing:", err)
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
        fmt.Println("Ping to",target ,"in", duration, "seconds")
    } else {
        icmpPingSuccess.WithLabelValues(target).Set(0)
        fmt.Println("Ping failed", err)
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
        fmt.Println("HTTP GET request failed:", err)
        return
    }
    defer resp.Body.Close()

    duration := time.Since(startTime).Seconds()
    fmt.Println("HTTP GET response from",url,"in", duration, "seconds")
}

func main() {
    target := "8.8.8.8"
    if len(os.Args) > 1 {
        target = os.Args[1]
    }

    go func() {
        for {
            ping(target)// Send ICMP ping to the target
            httpGet("https://www.google.com") // Perform HTTP GET request to Google
            time.Sleep(2 * time.Second)  // Wait for 2 seconds before the next iteration
        }
        }()
        

    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
