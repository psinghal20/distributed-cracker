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
    workerAssigned Worker
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
}

var dic ="abcdefghijklmnopqrstuvwxyz"
var start string
var counter = 0
var flag = false

func splitJob(job Job) {
    start = strings.Repeat("a", job.len)
    fmt.Printf("jobs: %v", job.len)
    permuteStrings("", job.len, job)
    fmt.Printf("\nTASK : %v\n", tasks)
}

func permuteStrings(prefix string, k int, job Job) {
    if k == 0 {
        counter++
        if counter == 5000 || prefix == strings.Repeat("z", job.len) {
            newTask := Task{
                job.jobId,
                Worker{},
                0,
                start,
                prefix,
            }
            fmt.Printf("\nTASK: %v\n", newTask)
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
