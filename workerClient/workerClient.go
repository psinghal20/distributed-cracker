package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
)

var conn net.Conn

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
    joinRequest()
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("Task received: %s\n", data)
    conn.Close()
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
