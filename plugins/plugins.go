package plugins

import (
	"context"
	"datadog_import/converter"
)

type Plugin interface {
	GetReader(context.Context) (converter.Reader, error)
	GetTransformer(context.Context) converter.Transformer
	GetWriters(context.Context) ([]converter.Writer, error)
}
