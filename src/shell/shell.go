package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vm"
	"utils"
)

func CMD(blocksOfUserMemory [][]byte, vMemory [10][]byte, R []byte, C byte, IC []byte, TI int, MODE byte ) {
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
		vm.ExecuteProgram(vMemory, R, &C, IC, TI, MODE, 'F')
		fmt.Println("---------------------")
	  }

	  if strings.Compare("2", text) == 0 {
		vm.ExecuteProgram(vMemory, R, &C, IC, TI, MODE, 'S')
		fmt.Println("---------------------")
	  }

	  if strings.Compare("3", text) == 0 {
		utils.PrintvMemory(vMemory)
		fmt.Println("---------------------")
	  }

	  if strings.Compare("4", text) == 0 {
		  utils.PrintUserMemory(blocksOfUserMemory, 29)
	  }

	  if strings.Compare("5", text) == 0 {
		break
	  }
	}

}