package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	const BYTES_READ = 8;
	f, err := os.Open("messages.txt")	
	if err != nil {
		log.Fatal("Failed to read file")	
	}

	defer f.Close()

	buf := make([]byte, BYTES_READ)
	var offset = int64(0)
	n, err := f.ReadAt(buf, offset);
	
	if n == 0 {
		return
	}

	for {
		fmt.Println("read: " + string(buf))
		n, err = f.ReadAt(buf, offset)

		offset, _ = f.Seek(BYTES_READ, 1)

		if err == io.EOF {
			if n > 0 && n < BYTES_READ {
				fmt.Println("read: " + string(buf[:n-1]))
			}
			break;
		}
	}

}
