package main

import (
    "fmt"
    "net"
    "encoding/json"
    "os"
    "sync"
)

var conn net.Conn

// Packet struct used to receive tasks from server
type Packet struct {
    Hash string
    Start string
    End string
}

var receivedPacket Packet
var mutex = &sync.Mutex{}

func main() {
    arguments := os.Args
    if len(arguments) != 3 {
        fmt.Println("Please provide correct arguments!")
        os.Exit(1)
    }

    serverAddr := arguments[1] + ":" + arguments[2]
    var err error
    conn, err = net.Dial("udp", serverAddr)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer conn.Close()
    joinRequest()
    for {
        readData()
    }
}

// Send network join request to the server
func joinRequest() {
    response := make([]byte, 1024)
    conn.Write([]byte("JOIN"))
    size, err := conn.Read(response)
    if err != nil {
        fmt.Println("Node failed to join the network as worker!", err)
        os.Exit(1)
    }
    if string(response[0:size]) == "1" {
        fmt.Println("Node joined the network as a worker!")
    }
}

// Read data from the UDP connection
// Check if the received data is Health check request or assigned task
func readData() {
    buf := make([]byte, 1024)
    size, err := conn.Read(buf);
    if err != nil {
        fmt.Println("Couldn't read the packet!", err)
        os.Exit(1)
    }
    if string(buf[:size]) == "CHECK" {
        respondPoll()
    } else {
        processPacket(buf[:size])
    }
}

// Respond to health check poll
func respondPoll() {
    mutex.Lock()
    defer mutex.Unlock()
    _, err := conn.Write([]byte("ACK"))
    if err != nil {
        fmt.Println("Failed to poll the server!", err)
        os.Exit(1)
    }
}

// Process the task packet received
func processPacket(buf []byte) {
    err := json.Unmarshal(buf, &receivedPacket);
    if err != nil {
        fmt.Println("Couldn't Unmarshal packet!", err)
        os.Exit(1)
    }
    flag = false
    resultFound = false
    completed = false
    // Execute the received task request in another goroutine
    // while listening for health check requests.
    go executeQuery()
}

func notFoundResponse() {
    if _, err := conn.Write([]byte("NOT FOUND")); err != nil {
        fmt.Println("Couldn't send the result server!", err)
        os.Exit(1)
    }
}

func foundResponse() {
    if _, err := conn.Write([]byte(fmt.Sprintf("FOUND:%s", result))); err != nil {
        fmt.Println("Couldn't send the result server!", err)
        os.Exit(1)
    }
}
