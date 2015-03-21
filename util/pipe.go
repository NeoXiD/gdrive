package util

import (
	"errors"
	"fmt"
	"os"
)

const SEEKABLE_PIPE_DEBUG = false
const SEEKABLE_PIPE_BUFFER_SIZE int64 = 1024 * 1024

type SeekablePipe struct {
	os.File
	offset int64
}

func NewSeekablePipe(input *os.File) *SeekablePipe {
	self := &SeekablePipe{File: *input}
	self.offset = 0
	return self
}

func (self *SeekablePipe) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64

	if SEEKABLE_PIPE_DEBUG {
		switch whence {
		case 0:
			fmt.Printf("Trying to seek from [%d] to [%d], skipping [%d] bytes\n", self.offset, (self.offset + offset), offset)
		case 1:
			fmt.Printf("Trying to seek from [%d] to [%d], skipping [%d] bytes\n", self.offset, offset, offset-self.offset)
		}
	}

	switch true {
	case
		(whence == 0 && offset < 0),
		(whence == 1 && offset < self.offset),
		(whence == 2):

		return 0, errors.New("SeekablePipe is unable to seek backwards!")

	case whence == 0:
		newOffset = self.offset + offset
		self.SkipBytes(newOffset - self.offset)

	case whence == 1:
		newOffset = offset
		self.SkipBytes(newOffset - self.offset)

	default:
		return 0, errors.New("Invalid whence argument given to SeekablePipe.Seek()")
	}

	self.offset = newOffset
	return self.offset, nil
}

func (self *SeekablePipe) Read(buffer []byte) (length int, err error) {
	length, err = self.File.Read(buffer)
	self.offset += int64(length)
	return
}

func (self *SeekablePipe) ReadAt(buffer []byte, offset int64) (int, error) {
	if SEEKABLE_PIPE_DEBUG {
		fmt.Printf("Reading [%d] bytes from offset [%d]...\n", len(buffer), offset)
	}

	_, err := self.Seek(offset, 1)
	if err != nil {
		return 0, err
	}

	return self.Read(buffer)
}

func (self *SeekablePipe) SkipBytes(count int64) {
	buffer := make([]byte, SEEKABLE_PIPE_BUFFER_SIZE)

	for count > 0 {
		if count < SEEKABLE_PIPE_BUFFER_SIZE {
			buffer = make([]byte, count)
		}

		self.File.Read(buffer)
		count -= SEEKABLE_PIPE_BUFFER_SIZE
	}
}
