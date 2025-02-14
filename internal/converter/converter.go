package converter

import (
	"context"
	"datadog_import/internal/common"
	"datadog_import/internal/logctx"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Converter interface {
	Convert()
}

type converter struct {
	ctx context.Context
	r   Reader
	t   Transformer
	ws  []Writer
}

func NewConverter(ctx context.Context, reader Reader, transform Transformer, writers []Writer) Converter {
	return &converter{
		ctx: ctx,
		r:   reader,
		t:   transform,
		ws:  writers,
	}
}

func (c *converter) Convert() {
	read := 0
	readError := 0
	transformed := 0
	transformError := 0
	written := 0
	writeError := 0

	logger := logctx.From(c.ctx)
	for {
		// Read
		id, name, data, err := c.r.Read(c.ctx)
		contextLogger := logger.WithFields(logrus.Fields{
			"Object":   name,
			"OriginId": id,
		})
		if err != nil {
			if errors.Is(err, common.ErrNoMoreData) {
				break
			}
			contextLogger.WithField("Reader", reflect.TypeOf(c.r)).Error(err)
			readError += 1
			continue
		}
		read += 1

		// Convert
		contextLogger.Debug("converting object")
		newObj, err := c.t(logctx.New(c.ctx, contextLogger), data)
		if err != nil {
			contextLogger.WithField("Transform", reflect.TypeOf(c.t)).Error(err)
			transformError += 1
			continue
		}
		transformed += 1

		// Write
		contextLogger.Debug("writing object")
		for _, writer := range c.ws {
			err = writer.Write(logctx.New(c.ctx, contextLogger), newObj, name)
			if err != nil {
				contextLogger.WithField("Writer", reflect.TypeOf(writer)).Error(err)
				writeError += 1
				continue
			}
			written += 1
		}
	}
	logger.Infof("%d objects read, %d transformed and %d written.", read, transformed, written)
	logger.Infof("%d read errors, %d transform errors and %d write errors.", readError, transformError, writeError)
}
