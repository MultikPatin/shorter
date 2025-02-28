package database

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultFilePermissions = 0666
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

func NewFileStorage(path string, isProducer bool) (*FileStorage, error) {
	var fileMode int
	if isProducer {
		fileMode = os.O_RDWR | os.O_CREATE | os.O_APPEND
	} else {
		fileMode = os.O_RDONLY | os.O_CREATE
	}

	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию: %w", err)
		}
	}

	file, err := os.OpenFile(path, fileMode, defaultFilePermissions)
	if err != nil {
		return nil, err
	}

	fs := &FileStorage{
		filename:   path,
		file:       file,
		writer:     bufio.NewWriterSize(file, 4096),
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

func (fs *FileStorage) Close() error {
	if err := fs.writer.Flush(); err != nil {
		return err
	}
	return fs.file.Close()
}

func parseLineToEvent(data []byte) (*Event, error) {
	event := &Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}
