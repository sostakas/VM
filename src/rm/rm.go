package rm

import (  
	"strconv"
	"utils"
	"shell"
)

type realMachine struct {  
	MODE   byte
	PTR    []byte
	C 	   byte
	IC	   []byte
	Memory []byte 
	R 	   []byte
	TI	   int
}
var isBlockAllocatedMap map[int]bool
var pageTableMap map[int]int

func New(MODE byte, PTR []byte, C byte, IC []byte, Memory []byte, R []byte, TI int) realMachine {  
    rm := realMachine {MODE, PTR, C, IC, Memory, R, TI}
    return rm
}

func (rm realMachine) RunProgram(Memory []byte) {  
	programArray := utils.ReadFile()
	var vMemory [10][]byte
	isBlockAllocatedMap = make(map[int]bool)
	pageTableMap = make(map[int]int)
	blocksOfUserMemory := utils.MakeArrayOfBlocks(rm.Memory, 10 * 4)
	utils.AllocateVirtualMemory(blocksOfUserMemory, &vMemory, isBlockAllocatedMap, pageTableMap)
	utils.AllocatePageTable(blocksOfUserMemory, rm.PTR, isBlockAllocatedMap, pageTableMap)
	loadProgramToVirtualMemory(programArray, vMemory)
	shell.CMD(blocksOfUserMemory, vMemory, rm.R, rm.C, rm.IC, rm.TI, rm.MODE)
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