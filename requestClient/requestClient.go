package main

import (
    "fmt"
    "net"
    "bufio"
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
    _, err = conn.Write([]byte(HASH + ":" + LENGTH + "\n"))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Print(data)
    data, err = bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("The cracked password is : %s", data)
    conn.Close()
}
