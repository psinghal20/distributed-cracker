package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Worker struct used to denote each worker node
type Worker struct {
	workerAddr *net.UDPAddr
	status     int
	taskId     string
}

// Task packet sent to the worker nodes
type Packet struct {
	Hash  string
	Start string
	End   string
}

// Map of all the connected workers with UDP addr as the key
var workers map[string]Worker = make(map[string]Worker)

// Global UDP conn object used for UDP server
var udpConn *net.UDPConn
var mutex = &sync.Mutex{}

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
	// Defered calls to close UDP server and
	// signify end of goroutine to waitgroup
	defer wg.Done()
	defer udpConn.Close()
	data := make([]byte, 1024)
	// Worker channel used to notify about health check responses
	workerChan := make(chan bool)

	// Run a seperate goroutine to check the health of each worker
	go checkHealthOfWorkers(workerChan)
	for {
		size, udpAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleUDPPacket(data[0:size], udpAddr, workerChan)
	}
}

func handleUDPPacket(data []byte, udpAddr *net.UDPAddr, workerChan chan bool) {
	switch res := string(data[:]); res {
	case "JOIN":
		handleWorkerJoinRequest(udpAddr) // Worker wants to join the network
	case "NOT FOUND":
		handleWorkerNotFoundRequest(udpAddr) // Worker couldn't find the passwd in the given range
	case "ACK":
		workerChan <- true // Response to health check request
	default:
		handleWorkerFoundRequest(res, udpAddr) // Worker found the passwd successfully!
	}
}

func handleWorkerJoinRequest(udpAddr *net.UDPAddr) {
	fmt.Printf("Worker join request from: %v\n", udpAddr)
	udpConn.WriteToUDP([]byte("1"), udpAddr)
	fmt.Printf("Worker %v joined the network!\n", udpAddr)
	setUpNewWorker(udpAddr)
	// After a new node joins the network, again distribute the free tasks.
	distributeTask()
}

func handleWorkerNotFoundRequest(udpAddr *net.UDPAddr) {
	fmt.Println("Worker Couldn't find the password")
	workerId := udpAddr.String()
	taskId := workers[workerId].taskId
	task, ok := tasks[taskId]
	freeWorker(workerId)
	if !ok {
		// If no task with received taskId is present in tasks map,
		// job must have been completed, do nothing.
		return
	}
	// Mark task as completed.
	task.status = CompletedTask
	tasks[taskId] = task
	jobId := task.jobId
	// Check if all the tasks for the job are completed
	// If yes, send "Not found" message to client
	if checkStatusOfJob(jobId) {
		sendResultToClient("Password Not Found!", jobId)
	}
	distributeTask()
}

func handleWorkerFoundRequest(data string, udpAddr *net.UDPAddr) {
	fmt.Printf("Worker node %v found the password : %s\n", udpAddr, data)
	workerId := udpAddr.String()
	taskId := workers[workerId].taskId
	task, ok := tasks[taskId]
	if !ok {
		// If no task with received taskId is present in tasks map,
		// job must have been completed, do nothing.
		return
	}
	// Mark task as completed
	task.status = CompletedTask
	tasks[taskId] = task
	jobId := task.jobId
	freeWorker(workerId)
	// Send the password to the client, data receieved as "FOUND:passwd"
	sendResultToClient(strings.Split(data, ":")[1], jobId)
	distributeTask()
}

func setUpNewWorker(udpAddr *net.UDPAddr) {
	workers[udpAddr.String()] = Worker{
		udpAddr,
		FreeWorker,
		NoTask,
	}
}

func sendWorkerTask(task Task, worker Worker) {
	// Send each task as packet struct using JSON marshalling
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

// Marks the worker free, removing taskId and updating its status
func freeWorker(workerId string) {
	worker := workers[workerId]
	worker.status = FreeWorker
	worker.taskId = NoTask
	workers[workerId] = worker
}

// Remove the task with given taskId
func removeTask(taskId string) {
	delete(tasks, taskId)
}

// Check status of the job, by checking if all the tasks are completed.
func checkStatusOfJob(jobId string) bool {
	for _, task := range tasks {
		if jobId == task.jobId && task.status != CompletedTask {
			return false
		}
	}
	return true
}

// Check the health of workers, by polling them in intevals of 10 sec
// Waits 5 sec for health check response
// Uses worker channel to receieve the notification for health check response
func checkHealthOfWorkers(workerChan chan bool) {
	for {
		time.Sleep(10 * time.Second)
		for workerId, worker := range workers {
			if _, err := udpConn.WriteToUDP([]byte("CHECK"), worker.workerAddr); err != nil {
				fmt.Println("Failed to communicate with worker!", err)
				mutex.Lock()
				unassignTask(worker.taskId)
				delete(workers, workerId)
				mutex.Unlock()
				continue
			}
			select {
			case <-workerChan:
				continue
			case <-time.After(5 * time.Second):
				mutex.Lock()
				unassignTask(worker.taskId)
				delete(workers, workerId)
				mutex.Unlock()
			}
		}
	}
}

// Mark the task as unassigned.
func unassignTask(taskId string) {
	task := tasks[taskId]
	task.status = UnassignedTask
	task.workerId = NoWorker
}
