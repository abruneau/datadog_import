package converter

import (
	"dynatrace_to_datadog/common"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Converter struct {
	*common.Config
	Reader    Reader
	Transform Transformer
	Writers   []Writer
}

func (c *Converter) Convert() {
	read := 0
	readError := 0
	transformed := 0
	transformError := 0
	written := 0
	writeError := 0
	for {
		// Read
		name, data, err := c.Reader.Read()
		contextLogger := c.Log.WithFields(logrus.Fields{
			"Object": name,
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
		newObj, err := c.Transform(data)
		if err != nil {
			contextLogger.WithField("Transform", reflect.TypeOf(c.Transform)).Error(err)
			transformError += 1
			continue
		}
		transformed += 1

		// Write
		for _, writer := range c.Writers {
			err = writer.Write(newObj, name)
			if err != nil {
				contextLogger.WithField("Writer", reflect.TypeOf(writer)).Error(err)
				writeError += 1
				continue
			}
			written += 1
		}
	}
	c.Log.Infof("%d objects read, %d transformed and %d written.", read, transformed, written)
	c.Log.Infof("%d read errors, %d transform errors and %d write errors.", readError, transformError, writeError)
}
