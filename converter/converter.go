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
			continue
		}

		// Convert
		newObj, err := c.Transform(data)
		if err != nil {
			contextLogger.WithField("Transform", reflect.TypeOf(c.Transform)).Error(err)
			continue
		}

		// Write
		for _, writer := range c.Writers {
			err = writer.Write(newObj, name)
			if err != nil {
				contextLogger.WithField("Writer", reflect.TypeOf(writer)).Error(err)
				continue
			}
		}
	}
}
