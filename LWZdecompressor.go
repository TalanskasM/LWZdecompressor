package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {

	path := flag.String("path", "compressedfile4.z", "specify file path")
	flag.Parse()

	decompressFile(*path)

}

func decompressFile(path string) {

	r, err := os.Open(path)
	errCheck(err)
	defer r.Close()

	w, err := os.Create("output.txt")
	errCheck(err)
	defer w.Close()

	dict := make(map[int]string)     //creates initial dictionary with 256 entries
	currCode := intitaliseDict(dict) //and sets code counter to 256

	var out string
	var firstCycle bool = true

	//main (while) loop, runs until no more bytes to read
	for {
		//////
		////// Read in file 3 bytes at a time
		////// if no bytes read decompress is done, if 2 are read then depad
		wr := bufio.NewWriter(w)
		b := make([]byte, 3)
		n, err := r.Read(b) //n = number of read bytes

		if err != nil {
			break //if no more bytes exit loop
		} else if n == 2 {
			//if 2 bytes left => remove padding and decode last code
			last := fmt.Sprintf("%08b", b[0])[4:] + fmt.Sprintf("%08b", b[1]) //remove padding
			lastInt, err := strconv.ParseInt(last, 2, 13)                     //convert to int
			errCheck(err)
			out, currCode = processbyte(int(lastInt), dict, out, currCode) //decode
			wr.WriteString(out)
			wr.Flush()
			break //finished decoding
		}
		//////
		//////

		//////
		//////separates the 3 bytes into 2 12-bit ints

		mid := fmt.Sprintf("%08b", b[1]) //middle byte => string

		first := fmt.Sprintf("%08b", b[0]) + mid[:4]    //first byte + first half of second byte
		firstInt, err := strconv.ParseInt(first, 2, 13) //convert to int
		errCheck(err)
		second := mid[4:] + fmt.Sprintf("%08b", b[2])     //second half of second byte + third byte
		secondInt, err := strconv.ParseInt(second, 2, 13) //convert to int
		errCheck(err)

		//////
		//////

		if firstCycle {
			out = dict[int(firstInt)]
			wr.WriteString(out)
		}

		//////

		//processes first byte, doesnt run during first cycle
		if !firstCycle {
			out, currCode = processbyte(int(firstInt), dict, out, currCode)
			wr.WriteString(out)
		}
		//checks if dictionary is full
		if currCode == 4096 {
			currCode = intitaliseDict(dict)
		}

		//processes second byte
		out, currCode = processbyte(int(secondInt), dict, out, currCode)
		wr.WriteString(out)
		//checks if dictionary is full again
		if currCode == 4096 {
			currCode = intitaliseDict(dict)
		}

		firstCycle = false

		wr.Flush() //writes to file

	}
}

func intitaliseDict(dict map[int]string) int {
	//sets initial entries in dictionary and returns next available code key
	for i := 0; i < 256; i++ {
		dict[i] = string(rune(i))
	}

	return 256

}

func processbyte(b int, dict map[int]string, prevString string, currCode int) (string, int) {

	var outString string

	if value, ok := dict[b]; ok {
		outString = value
		dict[currCode] = prevString + outString[:1]
	} else {
		outString = prevString + prevString[:1]
		dict[currCode] = prevString + prevString[:1]
	}

	newCode := currCode + 1

	return outString, newCode

}

func errCheck(e error) {
	if e != nil {
		panic(e)
	}
}
