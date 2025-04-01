package main

import (
	"flag"
	"net/http"
	"os"
	"time"

)


func main() {

    target := flag.String("ping", "8.8.8.8", "Ping Address")
    PORT := flag.String("port", "8080", "Port Address")
    flag.Parse()

    if !checkPort(*PORT){
      os.Exit(1)
    }

    go func() {
        for {
            ping(*target)// Send ICMP ping to the target
            httpGet("https://www.google.com") // Perform HTTP GET request to Google
            time.Sleep(2 * time.Second)  // Wait for 2 seconds before the next iteration
        }
        }()
            
    server := &http.Server{
        Addr: ":"+*PORT,
        Handler: routes(),
    }
    logger.Info("Starting server on :"+*PORT)
    err := server.ListenAndServe()
    logger.Error(err.Error())
    os.Exit(1)
}
