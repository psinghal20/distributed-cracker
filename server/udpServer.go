package main

import (
    "fmt"
    "net"
    "os"
)

type Worker struct {
    workerId int
    workerAddr *net.UDPAddr
    status int
}

var workers map[*net.UDPAddr]Worker = make(map[*net.UDPAddr]Worker)

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
        _, udpAddr, err := udpConn.ReadFromUDP(data)
        if err != nil {
            fmt.Println(err)
            return
        }
        go handleWorkerJoinRequest(udpAddr)
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
}

func handleWorkerFoundRequest(data string, udpAddr *net.UDPAddr) {
    fmt.Printf("Worker node %v found the password : %s\n", udpAddr, data)
}

func setUpNewWorker(udpAddr *net.UDPAddr) {
    newWorker := Worker{
        len(workers),
        udpAddr,
        0,
    }
    workers[udpAddr] = newWorker
    fmt.Printf("%v\n", workers)
}