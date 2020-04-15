package utils

import (
	"io/ioutil"
	"strings"
	"math/rand"
	"time"
	"strconv"
	"fmt"
)

func ReadFile() []string{
	// read program from file to byte array
	c, _ := ioutil.ReadFile("addNumbers.txt")
	// make array of commands 
	return strings.Split(string(c), "\n")
}

func MakeArrayOfBlocks(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func AllocateVirtualMemory(blocksOfUserMemory [][]byte, vMemory *[10][]byte, isBlockAllocatedMap map[int]bool, pageTableMap map[int]int) {
	rand.Seed(time.Now().UnixNano())
    min := 0
    max := 29
	// allocate 10 block to vMemory
	index := 0
	var random int
	for index < 10 {
		random = rand.Intn(max - min + 1) + min
		if WordIsEmpty(blocksOfUserMemory[random]) && !isBlockAllocatedMap[random] {
			isBlockAllocatedMap[random] = true
			pageTableMap[index] = random
			vMemory[index] = blocksOfUserMemory[random]
			index++
		}
	}
}

func AllocatePageTable(blocksOfUserMemory [][]byte, rPTR []byte, isBlockAllocatedMap map[int]bool, pageTableMap map[int]int) {
	index := 0
	min := 0
    max := 29
	for index < 30 {
		random := rand.Intn(max - min + 1) + min
		if WordIsEmpty(blocksOfUserMemory[random]) && !isBlockAllocatedMap[random] {
			isBlockAllocatedMap[random] = true
			pageTableBlockArray := MakeArrayOfBlocks(blocksOfUserMemory[random], 4)
			for i := 0; i < 10; i++ {
				for index, num := range strconv.Itoa(pageTableMap[i]) {
					pageTableBlockArray[i][index] = byte(num)
				}
			}
			// set PTR register
			if CountDigits(random) == 1 {
				rPTR[3] = strconv.Itoa(random)[0]
			}
			if CountDigits(random) == 2 {
				rPTR[2] = strconv.Itoa(random)[0]
				rPTR[3] = strconv.Itoa(random)[1]
			}
			break
		}
	}
}

func WordIsEmpty(s []byte) bool {
    for _, v := range s {
        if v != 0 {
            return false
        }
    }
    return true
}

func CountDigits(i int) (count int) {
	for i != 0 {

		i /= 10
		count = count + 1
	}
	return count
}

func PrintvMemory(vMemory [10][]byte) {
	for i := 0; i < 10; i++ {
		fmt.Println("block number: ",i, " ", string(vMemory[i]))
		// if !WordIsEmpty(vMemory[i]) {
		// 	fmt.Println("block number: ",i, " ,,", vMemory[i])
		// 	fmt.Println("block number: ",i, " ,,", string(vMemory[i]))
		// }
	}
}

func PrintRegisters(R []byte, C byte, IC []byte) {
	// fmt.Println("R: ", R)
	// fmt.Println("C: ", C)
	// fmt.Println("IC: ", IC)
	fmt.Println("R: ", string(R))
	fmt.Println("C: ", string(C))
	fmt.Println("IC: ", string(IC))
}

func PrintUserMemory(blocksOfUserMemory [][]byte, blocksNumber int) {
	for i := 0; i < blocksNumber; i++ {
		// fmt.Println("block number: : ",i, " ,,", blocksOfUserMemory[i])
		// if !WordIsEmpty(blocksOfUserMemory[i]) {
		// 	fmt.Println("block number: ",i, " ,,", string(blocksOfUserMemory[i]))
		// 	fmt.Println("block number: ",i, " ,,", blocksOfUserMemory[i])
		// }		
		fmt.Println("block number: ",i, " ", string(blocksOfUserMemory[i]))
	}
}