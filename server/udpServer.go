package main

import (
    "fmt"
    "net"
    "os"
    "encoding/json"
    "strings"
)

type Worker struct {
    workerAddr *net.UDPAddr
    status int
    taskId int
}

type Packet struct {
    Hash string
    Start string
    End string
}

var workers []Worker
var udpConn *net.UDPConn

func (w Worker) isBusy() int {
    return w.status
}

func udpServer(port string) {
    addr, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        fmt.Println("Couldn't resolve the udp address", err)
        os.Exit(1)
    }
    udpConn, err = net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer wg.Done()
    defer udpConn.Close()
    data := make([]byte, 1024)
    for {
        size, udpAddr, err := udpConn.ReadFromUDP(data)
        if err != nil {
            fmt.Println(err)
            return
        }
        go handleUDPPacket(data[0:size], udpAddr)
    }
}

func handleUDPPacket(data []byte, udpAddr *net.UDPAddr) {
    switch res := string(data[:]); res {
    case "JOIN":
        handleWorkerJoinRequest(udpAddr)
    case "NOT FOUND":
        handleWorkerNotFoundRequest(udpAddr)
    default:
        handleWorkerFoundRequest(res, udpAddr)
    }
}

func handleWorkerJoinRequest(udpAddr *net.UDPAddr) {
    fmt.Printf("Worker join request from: %v\n", udpAddr)
    setUpNewWorker(udpAddr)
    udpConn.WriteToUDP([]byte("1"), udpAddr)
    fmt.Printf("Worker %v joined the network!\n", udpAddr)
}

func handleWorkerNotFoundRequest(udpAddr *net.UDPAddr) {
    fmt.Println("Worker Couldn't find the password")
    freeWorker(udpAddr)
}

func handleWorkerFoundRequest(data string, udpAddr *net.UDPAddr) {
    fmt.Printf("Worker node %v found the password : %s\n", udpAddr, data)
    sendResultToClient(strings.Split(data, ":")[1], udpAddr)
    freeWorker(udpAddr)
}

func setUpNewWorker(udpAddr *net.UDPAddr) {
    newWorker := Worker{
        udpAddr,
        FreeWorker,
        NoTask,
    }
    workers = append(workers, newWorker)
}

func sendWorkerTask(task Task, worker Worker) {
    packet := Packet{
        jobs[task.jobId].hash,
        task.start,
        task.end,
    }
    data, err := json.Marshal(packet)
    if err != nil {
        fmt.Println("Couldn't marshal packet", err)
    }
    udpConn.WriteToUDP(data, worker.workerAddr)
}

func freeWorker(udpAddr *net.UDPAddr) {
    for inWorker, _ := range workers {
        if workers[inWorker].workerAddr.String() == udpAddr.String() {
            workers[inWorker].status = FreeWorker
            workers[inWorker].taskId = NoTask
        }
    }
    distributeTask()
}