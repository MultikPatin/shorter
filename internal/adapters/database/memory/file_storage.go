package memory

import (
	"bufio"
	"encoding/json"
	"fmt"
	"main/internal/constants"
	"main/internal/models"
	"os"
	"path/filepath"
)

// FileProducer writes events to a file using buffered I/O.
type FileProducer struct {
	file   *os.File      // Underlying file handle.
	writer *bufio.Writer // Buffered writer for efficient writes.
}

// FileConsumer reads events from a file using buffered scanning.
type FileConsumer struct {
	file    *os.File       // Underlying file handle.
	scanner *bufio.Scanner // Scanner for reading lines efficiently.
}

// NewFileProducer creates a new file producer instance, ensuring proper directory creation and file access.
func NewFileProducer(path string) (*FileProducer, error) {
	file, err := newFile(path, constants.DefaultProducerFileFlags, constants.DefaultFilePermissions)
	if err != nil {
		return nil, err
	}
	fs := &FileProducer{
		file:   file,
		writer: bufio.NewWriterSize(file, 4096), // Buffer size of 4KB.
	}
	return fs, nil
}

// NewFileConsumer creates a new file consumer instance, preparing the scanner for reading events.
func NewFileConsumer(path string) (*FileConsumer, error) {
	file, err := newFile(path, constants.DefaultConsumerFileFlags, constants.DefaultFilePermissions)
	if err != nil {
		return nil, err
	}
	fs := &FileConsumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
	return fs, nil
}

// newFile opens or creates a file with appropriate permissions and handling directory creation if needed.
func newFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) { // Check if directory doesn't exist.
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("couldn't create directory: %w", err)
		}
	}
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// WriteEvent serializes an Event model to JSON and appends it to the underlying file buffer.
func (fs *FileProducer) WriteEvent(event *models.Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := fs.writer.Write(data); err != nil {
		return err
	}
	if err := fs.writer.WriteByte('\n'); err != nil { // Separate records with newline.
		return err
	}

	err = fs.writer.Flush() // Ensure immediate write to disk.
	return err
}

// Close flushes any remaining data and closes the file handle for the producer.
func (fs *FileProducer) Close() error {
	if err := fs.writer.Flush(); err != nil {
		return err
	}
	return fs.file.Close()
}

// ReadAllEvents scans the entire file and deserializes its contents into a slice of Events.
func (fs *FileConsumer) ReadAllEvents() ([]*models.Event, error) {
	info, err := os.Stat(fs.file.Name())
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		return nil, nil
	}

	var events []*models.Event
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

// Close releases resources associated with the file consumer.
func (fs *FileConsumer) Close() error {
	return fs.file.Close()
}

// parseLineToEvent decodes a byte array containing JSON-encoded event data into an Event model.
func parseLineToEvent(data []byte) (*models.Event, error) {
	event := &models.Event{}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, err
	}
	return event, nil
}
