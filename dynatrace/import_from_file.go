package dynatrace

// import (
// 	"dynatrace_to_datadog/common"
// 	"dynatrace_to_datadog/dynatrace/synthetic"
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
// 	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
// 	"github.com/sirupsen/logrus"
// )

// // Convert Dynatrace Synthetics to Datadog from json definition files
// func (conf *Config) convertFromFiles() {
// 	reader, err := common.NewFileReader(conf.Input)
// 	if err != nil {
// 		common.Check(err)
// 	}
// 	writer, err := common.NewFileWriter(conf.Output)
// 	if err != nil {
// 		common.Check(err)
// 	}

// 	for {
// 		// Read
// 		fileName, data, err := reader.Read()
// 		contextLogger := conf.Log.WithFields(logrus.Fields{
// 			"Synthetics": fileName,
// 		})
// 		if err != nil {
// 			if errors.Is(err, common.ErrNoMoreData) {
// 				break
// 			}
// 			contextLogger.Error(err)
// 			continue
// 		}

// 		// Convert
// 		test, err := conf.convertTest(data)
// 		if err != nil {
// 			contextLogger.Error(err)
// 			continue
// 		}

// 		// Write
// 		err = writer.Write(test, fileName)

// 		if err != nil {
// 			contextLogger.Error(err)
// 			continue
// 		}
// 	}
// }

// func (conf *Config) convertTest(data []byte) (interface {
// 	MarshalJSON() ([]byte, error)
// }, error) {
// 	test := &monitors.SyntheticMonitor{}
// 	if err := json.Unmarshal(data, test); err != nil {
// 		return nil, err
// 	}
// 	conf.Log.Debug(fmt.Sprintf("Start processing test %s", test.Name))
// 	if test.Type == monitors.Types.Browser {
// 		browserTest := &browser.SyntheticMonitor{}
// 		if err := json.Unmarshal(data, browserTest); err != nil {
// 			return nil, err
// 		}

// 		return synthetic.ConvertBrowserTest(browserTest)
// 	} else {
// 		return nil, fmt.Errorf("SYnthetic type not supported: %s", test.Type)
// 	}
// }
