package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

const (
	perm = 0666
)

type Event struct {
	ID       int    `json:"uuid"`
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}

type FileStorage struct {
	filename   string
	file       *os.File
	writer     *bufio.Writer
	scanner    *bufio.Scanner
	isProducer bool
}

func NewFileStorage(filename string, isProducer bool) (*FileStorage, error) {
	var fileMode int
	if isProducer {
		fileMode = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		fileMode = os.O_RDONLY | os.O_CREATE
	}

	file, err := os.OpenFile(filename, fileMode, perm)
	if err != nil {
		return nil, err
	}

	fs := &FileStorage{
		filename:   filename,
		file:       file,
		writer:     bufio.NewWriter(file),
		scanner:    bufio.NewScanner(file),
		isProducer: isProducer,
	}

	return fs, nil
}

func (fs *FileStorage) WriteEvent(event *Event) error {
	if !fs.isProducer {
		return errors.New("cannot write in consumer mode")
	}
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
	return fs.writer.Flush()
}

func (fs *FileStorage) ReadAllEvents() ([]*Event, error) {
	if fs.isProducer {
		return nil, errors.New("cannot read in producer mode")
	}

	var events []*Event

	for fs.scanner.Scan() {
		data := fs.scanner.Bytes()
		event, err := parseLineToEvent(data)
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

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}

func parseLineToEvent(data []byte) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}
