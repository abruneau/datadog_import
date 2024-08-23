package converter

import (
	"context"
	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/logctx"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Converter struct {
	Reader    Reader
	Transform Transformer
	Writers   []Writer
}

func (c *Converter) Convert(ctx context.Context) {
	read := 0
	readError := 0
	transformed := 0
	transformError := 0
	written := 0
	writeError := 0

	logger := logctx.From(ctx)
	for {
		// Read
		id, name, data, err := c.Reader.Read(ctx)
		contextLogger := logger.WithFields(logrus.Fields{
			"Object":   name,
			"OriginId": id,
		})
		if err != nil {
			if errors.Is(err, common.ErrNoMoreData) {
				break
			}
			contextLogger.WithField("Reader", reflect.TypeOf(c.Reader)).Error(err)
			readError += 1
			continue
		}
		read += 1

		// Convert
		contextLogger.Debug("converting object")
		newObj, err := c.Transform(logctx.New(ctx, contextLogger), data)
		if err != nil {
			contextLogger.WithField("Transform", reflect.TypeOf(c.Transform)).Error(err)
			transformError += 1
			continue
		}
		transformed += 1

		// Write
		contextLogger.Debug("writing object")
		for _, writer := range c.Writers {
			err = writer.Write(logctx.New(ctx, contextLogger), newObj, name)
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
