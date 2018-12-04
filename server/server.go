package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
    "sync"
)

type Job struct {
    jobId int
    reqConn net.Conn
    hash string
    len string
}

var jobs []Job
var wg sync.WaitGroup

func main() {
    arguments := os.Args
    if len(arguments) == 2 {
        fmt.Println("Please provide a port number")
        os.Exit(1)
    }
    REQUEST_PORT := ":" + arguments[1]
    WORKER_PORT := ":" + arguments[2]
    wg.Add(2)
    go requestServer(REQUEST_PORT)
    go workerServer(WORKER_PORT)
    wg.Wait()
}

func requestServer(port string) {
    listener, err := net.Listen("tcp4", port)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer wg.Done()
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
    conn.Write([]byte("Working on you cracking request kindly wait!\n"))
}

func setUpNewJob(conn net.Conn, jobParams []string) {
    job := Job{
        len(jobs),
        conn,
        jobParams[0],
        jobParams[1],
    }
    jobs = append(jobs, job)
    fmt.Printf("%+v\n",jobs)
}
