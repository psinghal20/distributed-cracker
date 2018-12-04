package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
)

type Job struct {
    reqConn net.Conn
    hash string
    len string
}

var jobs []Job

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
        go handleNewJobRequest(c)
    }
}

func handleNewJobRequest(conn net.Conn) {
    fmt.Printf("New crack request from %s\n", conn.RemoteAddr().String())

    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Println("Couldn't set up the new job: ", err)
        return
    }
    jobParams := strings.Split(strings.TrimSuffix(data, "\n"), ":")
    setUpNewJob(conn, jobParams)
    conn.Write([]byte("Working on you cracking request kindly wait!"))
}

func setUpNewJob(conn net.Conn, jobParams []string) {
    job := Job{
        conn,
        jobParams[0],
        jobParams[1],
    }
    jobs = append(jobs, job)
}
