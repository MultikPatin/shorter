package memory

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"main/internal/constants"
	"main/internal/models"
	"os"
	"testing"
)

func TestFileProducer(t *testing.T) {
	t.Run("create_file_with_invalid_path", func(t *testing.T) {
		invalidPath := ""
		_, err := NewFileProducer(invalidPath)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("write_event_successfully", func(t *testing.T) {
		path := "./test_write_event.jsonl"
		defer os.Remove(path)

		producer, err := NewFileProducer(path)
		assert.NoError(t, err)
		defer producer.Close()

		event := models.Event{ID: 1, Origin: "Origin", Short: "Short"}
		err = producer.WriteEvent(&event)
		assert.NoError(t, err)
	})
}

func TestFileConsumer(t *testing.T) {
	t.Run("create_file_with_invalid_path", func(t *testing.T) {
		invalidPath := ""
		_, err := NewFileConsumer(invalidPath)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("read_all_events_from_non_empty_file", func(t *testing.T) {
		path := "./test_read_events.jsonl"
		defer os.Remove(path)

		eventData := []string{
			"{\"ID\":1,\"Origin\":\"Origin1\",\"Short\":\"Short1\"}",
			"{\"ID\":2,\"Origin\":\"Origin2\",\"Short\":\"Short2\"}",
		}
		writeSampleEvents(path, eventData)

		consumer, err := NewFileConsumer(path)
		assert.NoError(t, err)
		defer consumer.Close()

		events, err := consumer.ReadAllEvents()
		assert.NoError(t, err)
		assert.Len(t, events, 2)
	})

	t.Run("read_all_events_from_empty_file", func(t *testing.T) {
		path := "./empty_file.jsonl"
		defer os.Remove(path)

		f, _ := os.Create(path)
		f.Close()

		consumer, err := NewFileConsumer(path)
		assert.NoError(t, err)
		defer consumer.Close()

		events, err := consumer.ReadAllEvents()
		assert.NoError(t, err)
		assert.Empty(t, events)
	})
}

func writeSampleEvents(path string, data []string) {
	file, err := os.OpenFile(path, constants.DefaultProducerFileFlags, constants.DefaultFilePermissions)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range data {
		if _, err := w.WriteString(line + "\n"); err != nil {
			panic(err)
		}
	}
	if err := w.Flush(); err != nil {
		panic(err)
	}
}

func TestCloseFiles(t *testing.T) {
	t.Run("close_producer_and_consumer_files", func(t *testing.T) {
		path := "./test_close.jsonl"
		defer os.Remove(path)

		producer, err := NewFileProducer(path)
		assert.NoError(t, err)
		defer producer.Close()

		consumer, err := NewFileConsumer(path)
		assert.NoError(t, err)
		defer consumer.Close()

		err = producer.Close()
		assert.NoError(t, err)

		err = consumer.Close()
		assert.NoError(t, err)
	})

}
