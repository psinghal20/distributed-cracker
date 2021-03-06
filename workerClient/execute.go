package main

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// Execute the query by generating each permutation and computing its hash.
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
				completed = true
			}
		}
		return
	}
	for i := 0; i < 26; i++ {
		newPrefix := prefix + string(dic[i])
		if !completed {
			permuteStrings(newPrefix, k-1)
		}
	}
}

// check if the hash of test string is same to the given hash string.
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
