package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
)

func main() {
    arguments := os.Args
    if len(arguments) != 3 {
        fmt.Println("Please provide correct arguments!")
        os.Exit(1)
    }

    serverAddr := arguments[1] + ":" + arguments[2]

    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("Node joined the network as a worker!")
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("Task received: %s\n", data)
    conn.Close()
}