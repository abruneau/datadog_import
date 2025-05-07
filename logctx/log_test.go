package logctx

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	entry := logrus.NewEntry(logrus.New())

	newCtx := New(ctx, entry)

	// Verify that the context contains the log entry
	retrievedEntry, ok := newCtx.Value(contextKey{}).(*logrus.Entry)
	assert.True(t, ok)
	assert.Equal(t, entry, retrievedEntry)
}

func TestFrom(t *testing.T) {
	t.Run("with log entry in context", func(t *testing.T) {
		ctx := context.Background()
		entry := logrus.NewEntry(logrus.New())
		ctx = New(ctx, entry)

		retrievedEntry := From(ctx)
		assert.Equal(t, entry, retrievedEntry)
	})

	t.Run("without log entry in context", func(t *testing.T) {
		ctx := context.Background()
		retrievedEntry := From(ctx)
		assert.Equal(t, Default, retrievedEntry)
	})
}

func TestAddField(t *testing.T) {
	ctx := context.Background()
	entry := logrus.NewEntry(logrus.New())
	ctx = New(ctx, entry)

	// Add a field to the context
	key := "test_key"
	value := "test_value"
	newCtx := AddField(ctx, key, value)

	// Verify that the field was added
	retrievedEntry := From(newCtx)
	assert.Equal(t, value, retrievedEntry.Data[key])

	// Verify that the original context was not modified
	originalEntry := From(ctx)
	assert.NotEqual(t, value, originalEntry.Data[key])
}
