# Distributed Hash Cracker
Distributed password cracker written in Golang, inspired by one of the CMU's Distributed systems course assignment. It was mainly a project to gain experience with golang.

The whole app is distributed in 3 modules:
- Request Client: It makes hashing requests to central server.
- Central Server: Central server handles all the main tasks, managing workers and distributing jobs.
- Worker Client: Worker nodes perform the main computation and execute jobs received from central server.

NOTE: The project is still not complete and has some rough edges. I will try to smoothen them out as soon as possible!

# How to run
Generate Binaries for all three modules by running either(in respective directory):

    go build
or

    go install

### Running central server
    ./server [TCP(or Request Server) PORT] [UDP(or Worker Server) PORT]
### Running request client
    ./requestClient [Server IP Address] [Request Server PORT] [HASH] [Length of Password]
### Running worker node
    ./workerClient [Server IP Address] [Worker Server PORT]

# Contributing
Contributions are always welcome! I will try to look at them as soon as I can. Ping me once I forget to reply. Thanks! :smile:

PS: If someone has experience with distributed computing or testing apps for concurrency issues and have a minute to spare, please shoot me an email! I would love some advice!

# License
This project is licensed under MIT license. View [License](https://github.com/psinghal20/distributed-cracker/blob/master/LICENSE)

# References
- [CMU's Distributed Computing 15-440, Assignment Handout](https://www.cs.cmu.edu/afs/cs.cmu.edu/academic/class/15440-f11/P1/writeup.pdf)
- [gocrack](https://github.com/henrykhadass/gocrack)