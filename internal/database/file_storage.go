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
	encoder    *json.Encoder
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
		encoder:    json.NewEncoder(file),
		scanner:    bufio.NewScanner(file),
		isProducer: isProducer,
	}

	return fs, nil
}

func (fs *FileStorage) WriteEvent(event *Event) error {
	if !fs.isProducer {
		return errors.New("cannot write in consumer mode")
	}
	return fs.encoder.Encode(event)
}

func (fs *FileStorage) ReadEvent() ([]*Event, error) {
	if fs.isProducer {
		return nil, errors.New("cannot read in producer mode")
	}

	var events []*Event

	for fs.scanner.Scan() {
		line := fs.scanner.Text()
		event, err := parseLineToEvent(line)
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

func parseLineToEvent(line string) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal([]byte(line), event); err != nil {
		return nil, err
	}
	return event, nil
}
