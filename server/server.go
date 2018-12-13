package main

import (
    "fmt"
    "net"
    "os"
    "bufio"
    "strings"
    "sync"
    "strconv"
)

// Structure carrying data for individual task
// sent to each worker. A worker works on a single
// task at a time.
type Task struct {
    jobId string
    workerId string
    status int
    start string
    end string
}

// Struct for each received job
type Job struct {
    jobId string
    reqConn net.Conn
    hash string
    len int
}

var (
    // Map of all the received jobs
    jobs map[string]Job = make(map[string]Job)

    // Map of all the split tasks
    tasks map[string]Task = make(map[string]Task)
    wg sync.WaitGroup
)

func main() {
    arguments := os.Args
    if len(arguments) != 3 {
        fmt.Println("Please provide a port number")
        os.Exit(1)
    }
    // Ports for request server(TCP) & worker server(UDP)
    REQUEST_PORT := ":" + arguments[1]
    WORKER_PORT := ":" + arguments[2]

    // Wait for both goroutines to end.
    wg.Add(2)
    go requestServer(REQUEST_PORT)
    go udpServer(WORKER_PORT)
    wg.Wait()
}

// Start TCP server to listen to cracking requests
func requestServer(port string) {
    listener, err := net.Listen("tcp4", port)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Defered calls to close TCP server and
    // signify end of goroutine to waitgroup
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

    // Split job parameters received as HASH:LEN
    jobParams := strings.Split(strings.TrimSuffix(data, "\n"), ":")
    conn.Write([]byte("Working on you cracking request kindly wait!\n"))
    setUpNewJob(conn, jobParams)
}

func setUpNewJob(conn net.Conn, jobParams []string) {
    passLen, _ := strconv.Atoi(jobParams[1])
    job := Job{
        conn.RemoteAddr().String(),
        conn,
        jobParams[0],
        passLen,
    }
    jobs[conn.RemoteAddr().String()] = job
    splitJob(job)
    distributeTask()
}
// Global variables used for spliting job
// dic : set of alphabets
// start : start string for tasks, "a" * LEN
// counter used to split strings in batches
// flag: used to update start for each task
var dic ="abcdefghijklmnopqrstuvwxyz"
var start string
var counter = 0
var flag = false

func splitJob(job Job) {
    start = strings.Repeat("a", job.len)
    permuteStrings("", job.len, job)
}

func permuteStrings(prefix string, k int, job Job) {
    if k == 0 {
        counter++
        if counter == 5000 || prefix == strings.Repeat("z", job.len) {
            tasks[job.jobId+":"+start] = Task{
                job.jobId,
                NoWorker,
                UnassignedTask,
                start,
                prefix,
            }
            flag = true
            counter = 0
        }
        if flag {
            start = prefix
            flag = false
        }
        return
    }
    for i := 0; i < 26; i++ {
        newPrefix := prefix + string(dic[i])
        permuteStrings(newPrefix, k - 1, job)
    }
}

// Distribute tasks on first-come, first-serve basis
// to first free worker in map
func distributeTask() {
    fmt.Println("Distributing tasks")
    for taskId, task := range tasks {
        for workerId, worker := range workers {
            if task.status == UnassignedTask && worker.status == FreeWorker {
                task.workerId = workerId
                task.status = AssignedTask
                tasks[taskId] = task
                worker.status = BusyWorker
                worker.taskId = taskId
                workers[workerId] = worker
                sendWorkerTask(tasks[taskId], workers[workerId])
            }
        }
    }
}

// Delete the job with given jobId and all associated tasks
func removeJob(jobId string) {
    delete(jobs, jobId)
    for taskId, task := range tasks {
        if jobId == task.jobId {
            delete(tasks, taskId)
        }
    }
}

// Send the final received by workers to the requesting client
func sendResultToClient(result string, jobId string) {
    job, ok := jobs[jobId]
    if !ok {
        return
    }
    _, err := job.reqConn.Write([]byte(fmt.Sprintf("%s", result)))
    if err != nil {
        fmt.Println(err)
    }
    removeJob(jobId)
}