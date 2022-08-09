package main

import (
	"log"
	"os"
	"strconv"
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
	offset := int64(512)
	whence := int(1)

	isWriting := false
	currFile := 0

	// loop through 512byte chunks
	for i := 0; i < int(fileLength); i += 512 {

		//read a 512 byte chunk
		chunk := make([]byte, 512)
		readData, err := rawData.Read(chunk)
		_ = readData
		if err != nil {
			log.Fatal("there was an error reading the chunk", err)
		}

		//isolate the chunk header or data in that location
		header := []byte{chunk[0], chunk[1], chunk[2], chunk[3]}

		// if header is a jpeg header
		if header[0] == 255 && header[1] == 216 && header[2] == 255 && header[3] >= 224 && header[3] <= 239 {
			// are we currently writing a file already?
			if isWriting == false {

				//create a new file and number it
				numName := i / 512
				numNameString := strconv.Itoa(numName)
				newJpg, err := os.Create(numNameString + ".jpg")
				if err != nil {
					log.Fatal("there was an error creating a new jpeg", err)
				}
				defer newJpg.Close()

				// write the first chunk to the file
				write, err := newJpg.Write(chunk)
				_ = write
				if err != nil {
					log.Fatal("there was something bad with the write!!!", err)
				}

				// seek forward a 512 byte step in file
				seek, err := rawData.Seek(offset, whence)
				_ = seek
				if err != nil {
					log.Fatal("something went wrong with the seek", err)
				}

				isWriting = true
				currFile = numName
				newJpg.Close()

			} else if isWriting == true { // if we're already writing a file

				// open the existing file - corresponding currFile
				openJpg, err := os.OpenFile(strconv.Itoa(currFile)+".jpg", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Fatal("failed to open existing file to continue writing", err)
				}
				defer openJpg.Close()

				// continue writing the chunk to the existing file
				write, err := openJpg.Write(chunk)
				_ = write
				if err != nil {
					log.Fatal("looks like something went wrong with the CONTINUED WRITE!!!", err)
				}
			}

		}

	}

}
