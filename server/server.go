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

type Task struct {
    jobId int
    workerId int
    status int
    start string
    end string
}

type Job struct {
    jobId int
    reqConn net.Conn
    hash string
    len int
}

var jobs []Job
var tasks []Task
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
    go udpServer(WORKER_PORT)
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
    passLen, _ := strconv.Atoi(jobParams[1])
    job := Job{
        len(jobs),
        conn,
        jobParams[0],
        passLen,
    }
    jobs = append(jobs, job)
    splitJob(job)
    distributeTask()
}

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
            newTask := Task{
                job.jobId,
                -1,
                0,
                start,
                prefix,
            }
            tasks = append(tasks, newTask)
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

func distributeTask() {
    for inTask, _ := range tasks {
        for inWorker, _ := range workers {
            if tasks[inTask].status == 0 && workers[inWorker].status == 0 {
                tasks[inTask].workerId = inWorker
                tasks[inTask].status = 1 //1 = assigned, 2 = completed
                workers[inWorker].status = 1 //1 = busy
                workers[inWorker].taskId = inTask
                sendWorkerTask(tasks[inTask], workers[inWorker])
            }
        }
    }
}

func removeJob(jobId int) {
    tempTasks := tasks[:0]
    for _, task := range tasks {
        if task.jobId != jobId {
            tempTasks = append(tempTasks, task)
        }
    }
    tasks = tempTasks
    jobs = append(jobs[:jobId], jobs[jobId + 1:]...)
}

func sendResultToClient(result string, udpAddr *net.UDPAddr) {
    for inWorker, _ := range workers {
        if workers[inWorker].workerAddr.String() == udpAddr.String() {
            job := jobs[tasks[workers[inWorker].taskId].jobId]
            _, err := job.reqConn.Write([]byte(fmt.Sprintf("Password Found : %s", result)))
            if err != nil {
                fmt.Println(err)
            }
            removeJob(job.jobId)
        }
    }
}