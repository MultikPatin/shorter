package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultFilePermissions   = 0666
	defaultProducerFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultConsumerFileFlags = os.O_RDONLY | os.O_CREATE
)

type Event struct {
	ID       int    `json:"uuid"`
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

type FileProducer struct {
	file   *os.File
	writer *bufio.Writer
}
type FileConsumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewFileProducer(path string) (*FileProducer, error) {
	file, err := newFile(path, defaultProducerFileFlags, defaultFilePermissions)
	if err != nil {
		return nil, err
	}
	fs := &FileProducer{
		file:   file,
		writer: bufio.NewWriterSize(file, 4096),
	}
	return fs, nil
}

func NewFileConsumer(path string) (*FileConsumer, error) {
	file, err := newFile(path, defaultConsumerFileFlags, defaultFilePermissions)
	if err != nil {
		return nil, err
	}
	fs := &FileConsumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
	return fs, nil
}

func newFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию: %w", err)
		}
	}
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (fs *FileProducer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)

	if err != nil {
		return err
	}

	if _, err := fs.writer.Write(data); err != nil {
		return err
	}
	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}

	err = fs.writer.Flush()

	return err
}
func (fs *FileProducer) Close() error {
	if err := fs.writer.Flush(); err != nil {
		return err
	}
	return fs.file.Close()
}

func (fs *FileConsumer) ReadAllEvents() ([]*Event, error) {
	info, err := os.Stat(fs.file.Name())
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		return nil, nil
	}

	var events []*Event

	for fs.scanner.Scan() {
		event, err := parseLineToEvent(fs.scanner.Bytes())
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := fs.scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (fs *FileConsumer) Close() error {
	return fs.file.Close()
}

func parseLineToEvent(data []byte) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}
