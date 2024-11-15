package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrUnsupportedFile is returned when file size is 0.
	ErrUnsupportedFile = errors.New("unsupported file")
	// ErrOffsetExceedsFileSize is returned when offset exceeds file size.
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy copies fromPath to toPath from offset to offset+limit or to the end of the file if limit is 0.
func Copy(fromPath, toPath string, offset, limit int64) error {
	stat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	total := stat.Size()
	fmt.Println("File size", total)
	if total <= 0 {
		return ErrUnsupportedFile
	}

	if offset > total {
		return ErrOffsetExceedsFileSize
	}

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

	if limit > 0 && offset+limit < total {
		total = offset + limit
	}
	buf := make([]byte, 8)
	curr := offset
	for {
		read, errRead := fromFile.ReadAt(buf, curr)
		if errRead != nil && errRead != io.EOF {
			return errRead
		}
		if curr+int64(read) > total {
			read = int(total - curr)
		}
		_, errWrite := toFile.Write(buf[:read])
		if errWrite != nil {
			return errWrite
		}
		curr += int64(read)
		fmt.Printf("\rCopying... %d%%", int64(float64(curr-offset)/float64(total-offset)*100))
		if errRead == io.EOF || curr >= total {
			break
		}
	}
	fmt.Println()
	return nil
}
