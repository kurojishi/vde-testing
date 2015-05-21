package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

//CreateFile check if the file specified in file in path exist and is the correct size, if it isn't it create a new one(writing over the old one if necessary)
//then it open it and return the os.File object
func CreateFile(size int, path string) (*bytes.Buffer, error) {
	stats, err := os.Stat(path)
	if os.IsNotExist(err) || stats.Size() != int64(size) {
		dd := exec.Command("dd", "if=/dev/zero of="+path+" bs="+string(size)+"M count=1")
		err := dd.Run()
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.ReadFrom(file)
	log.Print("file is big: ", buf.Len()/(1024*1024))

	return &buf, nil
}
