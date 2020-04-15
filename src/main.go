package main

import (
	"rm"
)
var BLOCKS_IN_USER_MEMORY int = 30
var WORDS_IN_BLOCK int = 10
var BYTES_IN_WORD int = 4

func main()  {
	rm := rm.New('S', make([]byte, 4), 'F', make([]byte, 2), make([]byte, BLOCKS_IN_USER_MEMORY * WORDS_IN_BLOCK * BYTES_IN_WORD), make([]byte, 4), 10)
	rm.RunProgram(rm.Memory)
}