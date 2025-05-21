package memory

import (
	"context"
	"main/internal/adapters"
	"main/internal/models"
	"testing"
)

func BenchmarkInMemoryMethods(b *testing.B) {
	file := "test.json"
	logger := adapters.GetLogger()
	ctx := context.Background()
	short := "Short"

	db, _ := NewInMemoryDB(file, logger)
	logger.Info("Create InMemoryDB Connection")
	repo := NewLinksRepository(db)

	addedLink := models.AddedLink{
		CorrelationID: "CorrelationID",
		Short:         short,
		Origin:        "Origin",
	}

	var addedLinks []models.AddedLink
	for i := 0; i < 10; i++ {
		link := models.AddedLink{
			CorrelationID: "CorrelationID" + string(rune(i)),
			Short:         short + string(rune(i)),
			Origin:        "Origin" + string(rune(i)),
		}
		addedLinks = append(addedLinks, link)
	}

	b.ResetTimer()

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := repo.Add(ctx, addedLink)
			if err != nil {
				logger.Fatalw(err.Error(), "event", "Add")
			}
		}
	})
	b.Run("AddBatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := repo.AddBatch(ctx, addedLinks)
			if err != nil {
				logger.Fatalw(err.Error(), "event", "AddBatch")
			}
		}
	})
}
