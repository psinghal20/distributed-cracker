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
    taskId string
}

type Packet struct {
    Hash string
    Start string
    End string
}

var workers map[string]Worker = make(map[string]Worker)
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
    distributeTask()
}

func handleWorkerNotFoundRequest(udpAddr *net.UDPAddr) {
    fmt.Println("Worker Couldn't find the password")
    workerId := udpAddr.String()
    taskId := workers[workerId].taskId
    task, ok := tasks[taskId]
    freeWorker(workerId)
    if !ok {
        return
    }
    jobId := task.jobId
    if checkStatusOfJob(jobId) {
        sendResultToClient("Password Not Found!", jobId)
    }
}

func handleWorkerFoundRequest(data string, udpAddr *net.UDPAddr) {
    fmt.Printf("Worker node %v found the password : %s\n", udpAddr, data)
    workerId := udpAddr.String()
    taskId := workers[workerId].taskId
    task, ok := tasks[taskId]
    if !ok {
        return
    }
    jobId := task.jobId
    sendResultToClient(strings.Split(data, ":")[1], jobId)
    freeWorker(workerId)
}

func setUpNewWorker(udpAddr *net.UDPAddr) {
    workers[udpAddr.String()] = Worker{
        udpAddr,
        FreeWorker,
        NoTask,
    }
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

func freeWorker(workerId string) {
    worker := workers[workerId]
    worker.status = FreeWorker
    worker.taskId = NoTask
    workers[workerId] = worker
    distributeTask()
}

func removeTask(taskId string) {
    delete(tasks, taskId)
}

func checkStatusOfJob(jobId string) bool {
    for _, task := range tasks {
        if jobId == task.jobId && task.status != CompletedTask {
            return false
        }
    }
    return true
}