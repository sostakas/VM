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
	var rPTR [4]byte	// page table address
	MODE := byte('S')	// mode USER/SUPERVISOR
	var TI int			// timer
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

	// random allocation of user memory blocks, range 0-29
	allocateVirtualMemory(blocksOfUserMemory, &vMemory)
	// allocate 1 block in user memory for page table
	allocatePageTable(blocksOfUserMemory, &rPTR)
	// read file
	programArray := readFile()
	// load program to vMemory
	loadProgramToVirtualMemory(programArray, vMemory, &TI)
	// simple shell
	cmd(blocksOfUserMemory, vMemory, R, C, IC, TI, MODE)
	// cmd(blocksOfUserMemory, vMemory, R, C, IC, TI, MODE)

}

func cmd(blocksOfUserMemory [][]byte, vMemory [10][]byte, R []byte, C byte, IC []byte, TI int, MODE byte ) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
  
	for {
	  fmt.Println("1.Full mode")
	  fmt.Println("2.Step mode")
	  fmt.Println("3.Print virtual memory")
	  fmt.Println("4.Print user memory")
	  fmt.Println("5.Exit")
	  fmt.Print("-> ")
	  text, _ := reader.ReadString('\n')
	  text = strings.Replace(text, "\n", "", -1)
  
	  if strings.Compare("1", text) == 0 {
		executeProgram(vMemory, R, &C, IC, TI, MODE, 'F')
		fmt.Println("---------------------")
	  }

	  if strings.Compare("2", text) == 0 {
		executeProgram(vMemory, R, &C, IC, TI, MODE, 'S')
		fmt.Println("---------------------")
	  }

	  if strings.Compare("3", text) == 0 {
		printvMemory(vMemory)
		fmt.Println("---------------------")
	  }

	  if strings.Compare("4", text) == 0 {
		  printUserMemory(blocksOfUserMemory, 29)
	  }

	  if strings.Compare("5", text) == 0 {
		break
	  }
	}

}

func readFile() []string{
	// read program from file to byte array
	c, _ := ioutil.ReadFile("addNumbers.txt")
	// make array of commands 
	return strings.Split(string(c), "\n")
}

func allocateVirtualMemory(blocksOfUserMemory [][]byte, vMemory *[10][]byte) {
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
}

func allocatePageTable(blocksOfUserMemory [][]byte, rPTR *[4]byte) {
	index := 0
	min := 0
    max := 29
	for index < 30 {
		random := rand.Intn(max - min + 1) + min
		if wordIsEmpty(blocksOfUserMemory[random]) && !isBlockAllocatedMap[random] {
			isBlockAllocatedMap[random] = true
			pageTableBlockArray := makeArrayOfBlocks(blocksOfUserMemory[random], 4)
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

func executeProgram(vMemory [10][]byte, R []byte, C *byte, IC []byte, TI int, MODE byte, run_mode byte) {
	MODE = 'U'
	endProgram := false
	
	for _, cmd := range vMemory {
		vMemoryToWords := makeArrayOfBlocks(cmd, 4)
		if endProgram == true {
			break
		}
		for index, command := range vMemoryToWords {
			if strings.Compare("HALT", string(command)) == 0 && !wordIsEmpty(command) {
				// SI = 3 , HALT sets to 3
				// TEST()
				endProgram = true
				break
			}
			if !wordIsEmpty(command) {
				x, _ := strconv.Atoi(string(command[2]))
				y, _ := strconv.Atoi(string(command[3]))
				// address check , if not legal RM_PI = 1

				if len(command) == 4 {
					if run_mode == 'S' {
						fmt.Println("---------------------")
						fmt.Println("COMMAND IS: ", string(command))
						for  {
							fmt.Println("NEXT COMMAND: ", string(vMemoryToWords[index+1]))
							fmt.Println("---------------------")
							fmt.Println("1.next")
							fmt.Println("2.registers")
							fmt.Println("3.vMemory")
							fmt.Print("-> ")
							reader := bufio.NewReader(os.Stdin)
							text, _ := reader.ReadString('\n')
							text = strings.Replace(text, "\n", "", -1)
							if strings.Compare("1", text) == 0 {
								fmt.Println("---------------------")
								fmt.Println("command executed")
								execCommand(string(command), R, vMemory, x, y, IC, *C)
								fmt.Println("---------------------")
								break
							} 
							if strings.Compare("2", text) == 0 {
								fmt.Println("---------------------")
								fmt.Println("registers")
								printRegisters(R, *C, IC)
								fmt.Println("---------------------")
							} 
							if strings.Compare("3", text) == 0 {
								// should be vMemory of running command
								fmt.Println("---------------------")
								fmt.Println("virtual memory at block: ", x)
								fmt.Println(string(vMemory[x]))
								fmt.Println(vMemory[x])
								fmt.Println("---------------------")
							} 
						}
					}
					if run_mode == 'F' {
						execCommand(string(command), R, vMemory, x, y, IC, *C)
					}
				}
			}
		}
	}
}

func loadProgramToVirtualMemory(program []string, vMemory [10][]byte, TI *int) {
	blockNumber := 0
	startWriting := false
	indexFromBlockStart := 0
	for index, value := range program {
		// skip $AMJ
		if index == 0 {
			continue
		}
		// assign TI 
		if index == 1 {
			ti, _ := strconv.Atoi(string(value))
			*TI = ti
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
		// fmt.Println("block number: : ",i, " ,,", blocksOfUserMemory[i])
		// fmt.Println("block number: : ",i, " ,,", string(blocksOfUserMemory[i]))
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
	fmt.Println("IC: ", string(IC))
}

func printvMemory(vMemory [10][]byte) {
	// for i := 0; i < 10; i++ {
	// 	if !wordIsEmpty(vMemory[i]) {
	// 		fmt.Println("block number: ",i, " ,,", vMemory[i])
	// 	}
	// }
	for i := 0; i < 10; i++ {
		if !wordIsEmpty(vMemory[i]) {
			fmt.Println("block number: ",i, " ,,", string(vMemory[i]))
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

func channelDevice() {
	// SB := make([]byte, 4)
	// DB := make([]byte, 4)
	// ST := 0
	// DT := 0
}


func test() {
	// if (PI+SI) > 0 || TI == 0 {
	// 	MODE = 'S'
	// 	interruptHandling()
	// 	MODE = 'U'
	// }
}

func interruptHandling() {
	// implement
}

func execCommand(command string, R []byte, vMemory [10][]byte, x int, y int, IC []byte, C byte) {
	// R []byte, C byte, IC []byte, vMemory [10][]byte
	// fmt.Println("C: ", C)
	if len(command) == 4 {
		// fmt.Println("CURRENT COMMAND: ", command)
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
					C = 'T'
				} else {
					C = 'F'
				}
			case "BT":
				// does it assigns 0:2 or 2:4?
				if C == 'T' {
					for i := 0; i < 2; i++ {
						IC[i] = vMemory[x][4*y + i + 2]
					}
				}
			case "GD":
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter text: ")
				text, _ := reader.ReadString('\n')
				for index, value := range text {
					if index == len(text) - 1 {
						break
					}
					vMemory[x][index] = byte(value)
				}
				// RM_SI = 1 ?
				// RM_TI -= 2 ?
			case "PD":
				fmt.Println(string(vMemory[x][0:10*4]))
				// RM_SI = 1 ? 
				// RM_TI -= 2 ? 
			case "NT":
				// is it correct?
				if C != 'T' {
					IC[0] = strconv.Itoa(x)[0]
					IC[1] = strconv.Itoa(y)[0]
				}
			case "GO":
				// is it correct?
				IC[0] = strconv.Itoa(x)[0]
				IC[1] = strconv.Itoa(y)[0]
			default:
				// RM_PI = 2 command code invalid
		} 
	}
	// RM_TI--
	// RM_TEST()
	// IC++
	// if run_mode = 'S', next = false
}