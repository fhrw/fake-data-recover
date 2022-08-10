package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	// open raw "memory"
	rawData, err := os.OpenFile("card.raw", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal("something went wrong reading the file", err)
	}
	defer rawData.Close()
	rdInfo, err := rawData.Stat()
	if err != nil {
		log.Fatal("error with stat", err)
	}
	fileLength := rdInfo.Size()

	isWriting := false
	currFile := 0

	recoveredJpg := [][]byte{}

	fmt.Println(int(fileLength) / 512)

	// loop through 512byte chunks
	for i := 0; i < int(fileLength)/512; i++ {

		chunk := make([]byte, 512)
		readChunk, err := rawData.ReadAt(chunk, int64(i)*512)
		_ = readChunk
		if err != nil {
			panic(err)
		}

		header := []byte{chunk[0], chunk[1], chunk[2], chunk[3]}

		if header[0] == 255 && header[1] == 216 && header[2] == 255 && header[3] >= 224 && header[3] <= 239 && isWriting == false {
			isWriting = true
			newFile := []byte{}
			recoveredJpg = append(recoveredJpg, newFile)
		} else if header[0] == 255 && header[1] == 216 && header[2] == 255 && header[3] >= 224 && header[3] <= 239 && isWriting == true {
			currFile += 1
			newFile := []byte{}
			newFile = append(newFile, chunk...)
			recoveredJpg = append(recoveredJpg, newFile)
		}

		if header[0] != 255 && isWriting == true {
			recoveredJpg[currFile] = append(recoveredJpg[currFile], chunk...)
		}

	}

	fmt.Println(len(recoveredJpg))
}
