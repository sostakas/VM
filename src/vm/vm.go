package vm

import (
	"utils"
	"strings"
	"strconv"
	"fmt"
	"bufio"
	"os"
)

func ExecuteProgram(vMemory [10][]byte, R []byte, C *byte, IC []byte, TI int, MODE byte, run_mode byte) {
	MODE = 'U'
	endProgram := false
	
	for _, cmd := range vMemory {
		vMemoryToWords := utils.MakeArrayOfBlocks(cmd, 4)
		if endProgram == true {
			break
		}
		for index, command := range vMemoryToWords {
			if strings.Compare("HALT", string(command)) == 0 && !utils.WordIsEmpty(command) {
				// SI = 3 , HALT sets to 3
				// TEST()
				endProgram = true
				break
			}
			if !utils.WordIsEmpty(command) {
				x, _ := strconv.Atoi(string(command[2]))
				y, _ := strconv.Atoi(string(command[3]))
		
				// address check , if not legal RM_PI = 1

				if len(command) == 4 {
					if run_mode == 'S' {
						fmt.Println("COMMAND IS: ", string(command))
						for  {
							fmt.Println("NEXT COMMAND: ", string(vMemoryToWords[index+1]))
							fmt.Println("1.next")
							fmt.Println("2.registers")
							fmt.Println("3.vMemory")
							fmt.Print("-> ")
							reader := bufio.NewReader(os.Stdin)
							text, _ := reader.ReadString('\n')
							text = strings.Replace(text, "\n", "", -1)
							if strings.Compare("1", text) == 0 {
								execCommand(string(command), R, vMemory, x, y, IC, *C)
								break
							} 
							if strings.Compare("2", text) == 0 {
								utils.PrintRegisters(R, *C, IC)
								fmt.Println("---------------------")
							} 
							if strings.Compare("3", text) == 0 {
								// should be vMemory of running command
								utils.PrintvMemory(vMemory)
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

func execCommand(command string, R []byte, vMemory [10][]byte, x int, y int, IC []byte, C byte) {
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
				b, _ := strconv.Atoi(string(utils.MakeArrayOfBlocks(vMemory[x], 4)[y]))
				sum := a + b
				if utils.CountDigits(sum) == 1 {
					R[3] = strconv.Itoa(sum)[0]
				}
				if utils.CountDigits(sum) == 2 {
					R[2] = strconv.Itoa(sum)[0]
					R[3] = strconv.Itoa(sum)[1]
				}
				if utils.CountDigits(sum) == 3 {
					R[1] = strconv.Itoa(sum)[0]
					R[2] = strconv.Itoa(sum)[1]
					R[3] = strconv.Itoa(sum)[2]
				}
				if utils.CountDigits(sum) == 4 {
					R[0] = strconv.Itoa(sum)[0]
					R[1] = strconv.Itoa(sum)[1]
					R[2] = strconv.Itoa(sum)[2]
					R[3] = strconv.Itoa(sum)[3]
				}

			case "CR":
				if string(utils.MakeArrayOfBlocks(vMemory[x], 4)[y]) == string(R) {
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