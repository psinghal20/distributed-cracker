package main

import (
    "fmt"
    "net"
    "os"
)

type Worker struct {
    workerId int
    workerConn net.Conn
    status bool
}

var workers []Worker

func (w Worker) isBusy() bool {
    return w.status
}

func workerServer(port string) {
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
        go handleWorkerJoinRequest(c)
    }
}

func handleWorkerJoinRequest(conn net.Conn) {
    fmt.Printf("Worker join request from: %s\n", conn.RemoteAddr().String())
    setUpNewWorker(conn)
    fmt.Printf("Worker %s joined the network!\n", conn.RemoteAddr().String())
}

func setUpNewWorker(conn net.Conn) {
    newWorker := Worker{
        len(workers),
        conn,
        false,
    }
    workers = append(workers, newWorker)
    fmt.Printf("%v\n", workers)
}