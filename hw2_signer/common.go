package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"sort"
)

type job func(in, out chan interface{})

const (
	MaxInputDataLen = 100
)

var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)
	return dataHash
}

var SingleHash = func(in, out chan interface{}) {
	for data := range in {
		dataStr := fmt.Sprintf("%v", data)
		dataHash := DataSignerCrc32(dataStr) + "~" + DataSignerCrc32(DataSignerMd5(dataStr))
		fmt.Printf("SingleHash, value: %v\t hash: %v\n", data, dataHash)
		out <- dataHash
	}
}

var MultiHash = func(in, out chan interface{}) {
	for data := range in {
		var resHash string
		dataStr := fmt.Sprintf("%v", data)
		for i := 0; i < 6; i++ {
			datHash := DataSignerCrc32(fmt.Sprintf("%d%s", i, dataStr))
			resHash += datHash
		}
		fmt.Printf("MultiHash, value: %s\t hash: %v\n", data, resHash)
		out <- resHash
	}
}

var CombineResults = func(in, out chan interface{}) {
	var resHash []string
	for data := range in {
		dataStr := fmt.Sprintf("%v", data)
		resHash = append(resHash, dataStr)
	}
	sort.Strings(resHash)
	out <- strings.Join(resHash, "_")
}
