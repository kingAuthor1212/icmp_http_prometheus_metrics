package main

import (
	"flag"
	"net/http"
	"os"
	"time"

)


func main() {


    // get ping address 
    target := flag.String("ping", "8.8.8.8", "Ping Address")

    // get port number
    PORT := flag.String("port", "8080", "Port Address")
    flag.Parse()

    // check if within range of port address
    if !checkPort(*PORT){
      os.Exit(1)
    }

    // gorout with for
    go func() {
        for {
            ping(*target)// Send ICMP ping to the target
            httpGet("https://www.google.com") // Perform HTTP GET request to Google
            time.Sleep(2 * time.Second)  // Wait for 2 seconds before the next iteration
        }
        }()
    // http server instance
    server := &http.Server{
        Addr: ":"+*PORT,
        Handler: routes(),
    }
    logger.Info("Starting server on :"+*PORT)
    if err := server.ListenAndServe(); err != nil {
        logger.Error("Failed to start server: " + err.Error())
        os.Exit(1)
    }
}
