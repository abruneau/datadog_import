package dynatrace

// import (
// 	"fmt"
// 	"path"
// 	"strings"

// 	"dynatrace_to_datadog/common"
// 	"dynatrace_to_datadog/dynatrace/api"
// 	"dynatrace_to_datadog/dynatrace/synthetic"

// 	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
// 	"github.com/sirupsen/logrus"
// )

// func (conf *Config) Convert() {
// 	if conf.Input != "" {
// 		conf.convertFromFiles()
// 	}
// 	if conf.ApiKey != "" && conf.URL != "" {
// 		conf.importFromApi()
// 	}
// }

// func (conf *Config) importFromApi() {
// 	client := api.NewClient(conf.URL, conf.ApiKey)
// 	tests, err := client.List()
// 	common.Check(err)
// 	outputDirectory, err := common.CreateSubDirectories(conf.Output, "api/")
// 	for _, test := range tests {
// 		var res []byte
// 		browserTest := &browser.SyntheticMonitor{}
// 		client.Get(test.ID, browserTest)
// 		contextLogger := conf.Log.WithFields(logrus.Fields{
// 			"Synthetics": test.Name,
// 		})

// 		convertedTest, err := synthetic.ConvertBrowserTest(browserTest)
// 		if err != nil {
// 			contextLogger.Error(err)
// 			continue
// 		}

// 		res, err = convertedTest.MarshalJSON()
// 		if err != nil {
// 			contextLogger.Error(err)
// 			continue
// 		}

// 		testName := strings.ReplaceAll(test.Name, "/", "%2F")
// 		output := path.Join(outputDirectory, fmt.Sprintf("%s.json", testName))
// 		err = common.WriteFileToDisk(output, res)
// 		if err != nil {
// 			contextLogger.Error(err)
// 			continue
// 		}
// 	}
// }
