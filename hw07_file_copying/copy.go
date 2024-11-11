package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	stat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	total := stat.Size()
	fmt.Println("File size", total)

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	_ = offset
	_ = limit
	buf := make([]byte, 8)
	var curr int64
	for {
		read, errRead := fromFile.ReadAt(buf, curr)
		if errRead != nil && errRead != io.EOF {
			return errRead
		}
		_, errWrite := toFile.Write(buf[:read])
		if errWrite != nil {
			return errWrite
		}
		curr += int64(read)
		fmt.Printf("\rCopying... %d%%", int64(float64(curr)/float64(total)*100))
		if errRead == io.EOF {
			break
		}
	}
	fmt.Println()
	return nil
}
