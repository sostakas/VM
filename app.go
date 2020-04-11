package main

import (
    "strings"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
	"strconv"
	"bufio"
	"os"
)

var BLOCKS_IN_USER_MEMORY int = 30
var WORDS_IN_BLOCK int = 10
var BYTES_IN_WORD int = 4
var BLOCKS_IN_VIRTUAL_MEMORY int = 10
var isBlockAllocatedMap map[int]bool
var pageTableMap map[int]int

func main()  {
	// virtual machine registers
	R := make([]byte, 4)
	C := byte('F')
	IC := make([]byte, 2)

	// allocate user memory
	userMemory := make([]byte, BLOCKS_IN_USER_MEMORY * WORDS_IN_BLOCK * BYTES_IN_WORD)
	// array of user memory blocks
	blocksOfUserMemory := makeArrayOfBlocks(userMemory, WORDS_IN_BLOCK * BYTES_IN_WORD)
	// allocate vMemory
	var vMemory [10][]byte
	// is user memory block already assigned to virtual memory
	isBlockAllocatedMap = make(map[int]bool)
	// key = virtual memory block index, value = user memory block index
	pageTableMap = make(map[int]int)

	// move to function
	// random allocation of user memory blocks, range 0-29
	rand.Seed(time.Now().UnixNano())
    min := 0
    max := 29
	// allocate 10 block to vMemory
	index := 0
	var random int
	for index < 10 {
		random = rand.Intn(max - min + 1) + min
		if wordIsEmpty(blocksOfUserMemory[random]) && !isBlockAllocatedMap[random] {
			isBlockAllocatedMap[random] = true
			pageTableMap[index] = random
			vMemory[index] = blocksOfUserMemory[random]
			index++
		}
	}
	// ______________________
	// move to function
	// allocate 1 block in user memory for page table
	index = 0
	for index < 30 {
		random = rand.Intn(max - min + 1) + min
		if wordIsEmpty(blocksOfUserMemory[random]) && !isBlockAllocatedMap[random] {
			isBlockAllocatedMap[random] = true
			// fmt.Println("PAGE TABLE BLOCK INDEX: ", random)	// delete later
			pageTableBlockArray := makeArrayOfBlocks(blocksOfUserMemory[random], 4)
			for i := 0; i < 10; i++ {
				for index, num := range strconv.Itoa(pageTableMap[i]) {
					pageTableBlockArray[i][index] = byte(num)
				}
			}
			break
		}
	}
	// _______________ 
	// move to function
	// read program from file to byte array
	// c, _ := ioutil.ReadFile("ha.txt")
	c, _ := ioutil.ReadFile("addNumbers.txt")

	// make array of commands 
	programArray := strings.Split(string(c), "\n")
	// load program to vMemory
	loadProgramToVirtualMemory(programArray, vMemory)
	// _______________
	// execute program
	// printRegisters(R, C, IC)

	executeProgram(programArray, vMemory, R, &C, IC)
	printRegisters(R, C, IC)
	fmt.Println("PAGE TABLE: ", pageTableMap)
	fmt.Println("VIRTUAL MEMORY:")
	printvMemory(vMemory)
	fmt.Println("USER MEMORY:")
	printUserMemory(blocksOfUserMemory, 29)
}

func executeProgram(programArray []string, vMemory [10][]byte, R []byte, C *byte, IC []byte) {
	for _, command := range programArray {
		x, _ := strconv.Atoi(string(command[2]))
		y, _ := strconv.Atoi(string(command[3]))
		// fmt.Println("command is: ", command)

		if command == "HALT" {
			break
		}

		if len(command) == 4 {
			switch string(command[:2]) {
				case "LR":
					for i := 0; i < 4; i++ {
						R[i] = vMemory[x][4*y + i]
					}
				case "SR":
					for i := 0; i < 4; i++ {
						vMemory[x][4*y + i] = R[i]
					}
				case "AD":
					a, _ := strconv.Atoi(string(R))
					b, _ := strconv.Atoi(string(makeArrayOfBlocks(vMemory[x], 4)[y]))
					sum := a + b
					if CountDigits(sum) == 1 {
						R[3] = strconv.Itoa(sum)[0]
					}
					if CountDigits(sum) == 2 {
						R[2] = strconv.Itoa(sum)[0]
						R[3] = strconv.Itoa(sum)[1]
					}
					if CountDigits(sum) == 3 {
						R[1] = strconv.Itoa(sum)[0]
						R[2] = strconv.Itoa(sum)[1]
						R[3] = strconv.Itoa(sum)[2]
					}
					if CountDigits(sum) == 4 {
						R[0] = strconv.Itoa(sum)[0]
						R[1] = strconv.Itoa(sum)[1]
						R[2] = strconv.Itoa(sum)[2]
						R[3] = strconv.Itoa(sum)[3]
					}

				case "CR":
					if string(makeArrayOfBlocks(vMemory[x], 4)[y]) == string(R) {
						*C = 'T'
					} else {
						*C = 'F'
					}
				case "BT":
					// does it assigns 0:2 or 2:4?
					if *C == 'T' {
						for i := 0; i < 2; i++ {
							IC[i] = vMemory[x][4*y + i + 2]
						}
					}
				case "GD":
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Waiting for input: ")
					text, _ := reader.ReadString('\n')
					for index, value := range text {
						if index == len(text) - 1 {
							break
						}
						vMemory[x][index] = byte(value)
					}
				case "PD":
					fmt.Println(string(vMemory[x][0:10*4]))
				case "NT":
					// is it correct?
					if *C != 'T' {
						IC[0] = strconv.Itoa(x)[0]
						IC[1] = strconv.Itoa(y)[0]
					}
				case "GO":
					// is it correct?
					IC[0] = strconv.Itoa(x)[0]
					IC[1] = strconv.Itoa(y)[0]
			} 
		}
	}
}

func loadProgramToVirtualMemory(program []string, vMemory [10][]byte) {
	blockNumber := 0
	startWriting := false
	indexFromBlockStart := 0
	for index, value := range program {
		// skip $AMJ
		if index == 0 {
			continue
		}
		// end program loading
		if value == "$END" {
			break
		}

		// if starts with $ and is size of command set start writing to true
		if value[0] == 36 && len(value) == 4 {
			i, _ := strconv.Atoi(string(value[2]))
			blockNumber = i
			startWriting = true
			indexFromBlockStart = 0
			continue
		}
		if startWriting {
			for i := 0; i < len(value); i++ {
				if indexFromBlockStart == 40 {
					blockNumber++
					indexFromBlockStart = 0
				}
				vMemory[blockNumber][indexFromBlockStart] = value[i]
				indexFromBlockStart++
			}
		}
	 }
}
func printUserMemory(blocksOfUserMemory [][]byte, blocksNumber int) {
	for i := 0; i < blocksNumber; i++ {
		// fmt.Println("i: ",i, " ,,", blocksOfUserMemory[i])
		// fmt.Println("i: ",i, " ,,", string(blocksOfUserMemory[i]))
		if !wordIsEmpty(blocksOfUserMemory[i]) {
			fmt.Println("block number: ",i, " ,,", string(blocksOfUserMemory[i]))

		}		
	}
}

func makeArrayOfBlocks(buf []byte, lim int) [][]byte {
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

func wordIsEmpty(s []byte) bool {
    for _, v := range s {
        if v != 0 {
            return false
        }
    }
    return true
}

func printRegisters(R []byte, C byte, IC []byte) {
	// fmt.Println("R: ", R)
	// fmt.Println("C: ", C)
	// fmt.Println("IC: ", IC)
	fmt.Println("R: ", string(R))
	fmt.Println("C: ", string(C))
	// fmt.Println("IC: ", string(IC))
}

func printvMemory(vMemory [10][]byte) {
	// for i := 0; i < 10; i++ {
	// 	if !wordIsEmpty(vMemory[i]) {
	// 		fmt.Println("block number: ",i, " ,,", vMemory[i])
	// 	}
	// }
	for i := 0; i < 10; i++ {
		if !wordIsEmpty(vMemory[i]) {
			fmt.Println("i: ",i, " ,,", string(vMemory[i]))
		}
	}
}

func CountDigits(i int) (count int) {
	for i != 0 {

		i /= 10
		count = count + 1
	}
	return count
}