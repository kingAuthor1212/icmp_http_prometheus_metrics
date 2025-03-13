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
    icmpPingSuccess = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "icmp_ping_success",
            Help: "1 if the ping was successful else 0.",
        },
        []string{"target"},
    )

    icmpPingResponseTime = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "icmp_ping_response_time",
            Help: "Response time of ICMP ping in seconds.",
        },
        []string{"target"},
    )
)

func init() {
    prometheus.MustRegister(icmpPingSuccess)
    prometheus.MustRegister(icmpPingResponseTime)
}

func ping(target string) {
    // Listen on all interfaces
    c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer c.Close()

    start := time.Now()

    // Prepare the ICMP echo request message
    message := icmp.Message{
        Type: ipv4.ICMPTypeEcho,
        Code: 0,
        Body: &icmp.Echo{
            ID:   os.Getpid() & 0xffff,
            Seq:  1,
            Data: []byte("Ping"),
        },
    }

    // Marshal the ICMP message
    msgBytes, err := message.Marshal(nil)
    if err != nil {
        fmt.Println("Error marshaling message:", err)
        return
    }

    // Resolve the target IP address
    addr, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        fmt.Println("Error resolving address:", err)
        return
    }

    // Send the ICMP message
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
        fmt.Printf("Ping to %s successful, duration: %v seconds\n", target, duration)
    } else {
        icmpPingSuccess.WithLabelValues(target).Set(0)
        fmt.Println("Ping failed:", err)
    }
}

// https://pkg.go.dev/net/http
func httpGet(url string) {
    startTime := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("HTTP GET request failed:", err)
        return
    }
    defer resp.Body.Close()

    duration := time.Since(startTime).Seconds()
    fmt.Println("HTTP GET response from",url,"in", duration)
}

func main() {
    target := "8.8.8.8" // Default target for ICMP ping
    if len(os.Args) > 1 {
        target = os.Args[1] // Accept target as command-line argument
    }

    go func() {
        for {
            ping(target)
            httpGet("https://www.google.com") // Perform HTTP GET request
            time.Sleep(5 * time.Second) // Ping every 5 seconds
        }
        }()
        

    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
