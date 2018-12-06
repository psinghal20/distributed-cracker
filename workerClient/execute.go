package main

import (
    "fmt"
    "crypto/sha1"
    "io"
)

func executeQuery() {
    permuteStrings("", len(receivedPacket.Start))
    if completed {
        if resultFound {
            foundResponse()
        } else {
            notFoundResponse()
        }
    }
}
var dic = "abcdefghijklmnopqrstuvwxyz"
var flag bool = false
var resultFound = false
var completed = false
var result string

func permuteStrings(prefix string, k int) {
    if k == 0 {
        if receivedPacket.Start == prefix {
            flag = true
        }
        if flag {
            checkString(prefix, receivedPacket.Hash)
            if prefix == receivedPacket.End {
                fmt.Println("TASK COMPLETED!")
                completed = true
            }
        }
        return
    }
    for i := 0; i < 26; i++ {
        newPrefix := prefix + string(dic[i])
        if !completed {
            permuteStrings(newPrefix, k - 1)
        }
    }
}

func checkString(test string, hash string) {
    h := sha1.New()
    io.WriteString(h, test)
    if fmt.Sprintf("%x", h.Sum(nil)) == hash {
        resultFound = true
        completed = true
        result = test
        fmt.Printf("RESULT FOUND: %s\n", test)
    }
}