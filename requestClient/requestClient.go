package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    arguments := os.Args
    if len(arguments) != 5 {
        fmt.Println("Please provide correct arguments!")
        os.Exit(1)
    }

    serverAddr := arguments[1] + ":" + arguments[2]

    conn, err := net.Dial("tcp", serverAddr)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    HASH := arguments[3]
    LENGTH := arguments[4]
    // Send the request to the server, as "HASH:LEN"
    _, err = conn.Write([]byte(HASH + ":" + LENGTH + "\n"))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    result := make([]byte, 1024)
    // Read acknowledgment of received task
    size, err := conn.Read(result)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Print(string(result[:size]))
    // Read the result from the server.
    size, err = conn.Read(result)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("The cracked password is : %s", string(result[:size]))
    conn.Close()
}
