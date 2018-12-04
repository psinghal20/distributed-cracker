package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
)

func main() {
    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide a port number")
        os.Exit(1)
    }
    PORT := ":" + arguments[1]
    listener, err := net.Listen("tcp4", PORT)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer listener.Close()

    for {
        c, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            return
        }
        go handleConnection(c)
    }
}

func handleConnection(conn net.Conn) {
    fmt.Printf("accepting request from %s\n", conn.RemoteAddr().String())
    for {
        data, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            return
        }
        data = strings.TrimSpace(string(data))
        if data == "quit" {
            fmt.Printf("Closing Connection from %s\n", conn.RemoteAddr().String())
            break
        }
        fmt.Println(data)
        conn.Write([]byte("Hey you just connected to me!"))
    }
    conn.Close()
}