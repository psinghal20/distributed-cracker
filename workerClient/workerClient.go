package main

import (
    "fmt"
    "net"
    "encoding/json"
    "os"
)

var conn net.Conn

type Packet struct {
    Hash string
    Start string
    End string
}

var receivedPacket Packet

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
        flag = false
        resultFound = false
        completed = false
        readPacket()
        executeQuery()
    }
}

func joinRequest() {
    response := make([]byte, 1024)
    conn.Write([]byte("JOIN"))
    size, _ := conn.Read(response)
    if string(response[0:size]) == "1" {
        fmt.Println("Node joined the network as a worker!")
    } else {
        fmt.Println("Node failed to join the network as worker!")
        os.Exit(1)
    }
}

func notFoundResponse() {
    conn.Write([]byte("NOT FOUND"))
}

func foundResponse() {
    conn.Write([]byte(fmt.Sprintf("FOUND:%s", result)))
}

func readPacket() {
    buf := make([]byte, 1024)
    size, err := conn.Read(buf);
    if err != nil {
        fmt.Println("Couldn't read the packet!", err)
    }
    err = json.Unmarshal(buf[:size], &receivedPacket);
    if err != nil {
        fmt.Println(err)
    }
}
