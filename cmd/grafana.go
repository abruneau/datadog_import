package cmd

import (
	"datadog_import/plugins/grafana/dashboard/types"
	"datadog_import/plugins/grafana/dashboard/widgets/converter"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var grafanaCmd = &cobra.Command{
	Use:   "grafana",
	Short: "Convert Grafana queries to Datadog queries",
	Long:  `Convert queries from various Grafana datasources (Azure Monitor, CloudWatch, Stackdriver, Prometheus, Loki) to Datadog queries.`,
	RunE:  runGrafana,
}

var (
	queryType string
	queries   []string
)

func init() {
	rootCmd.AddCommand(grafanaCmd)
	grafanaCmd.Flags().StringVarP(&queryType, "type", "t", "", "Query type (grafana-azure-monitor-datasource, cloudwatch, stackdriver, prometheus, loki)")
	grafanaCmd.Flags().StringArrayVarP(&queries, "query", "q", []string{}, "Queries to convert (can be specified multiple times)")
	grafanaCmd.MarkFlagRequired("type")
	grafanaCmd.MarkFlagRequired("query")
}

func runGrafana(cmd *cobra.Command, args []string) error {
	conv, err := converter.NewConverter(queryType)
	if err != nil {
		return fmt.Errorf("error creating converter: %v", err)
	}

	for _, query := range queries {

		// Create a panel with our query
		panel := types.Panel{
			Title: "Query Conversion",
			Targets: []map[string]interface{}{
				{
					"expr":      query,
					"refId":     "A",
					"queryType": "metrics",
				},
			},
		}

		// For Azure Monitor, we need to set additional fields
		if queryType == "grafana-azure-monitor-datasource" {
			panel.Targets[0]["queryType"] = "Azure Monitor"
		}

		// Convert the query
		queries, formulas, err := conv.Parse(panel, true, true)
		if err != nil {
			return fmt.Errorf("error converting query: %v", err)
		}

		// Create a result structure for JSON output
		result := struct {
			SourceQuery string        `json:"source_query"`
			Queries     []interface{} `json:"queries"`
			Formulas    []interface{} `json:"formulas"`
		}{
			SourceQuery: query,
			Queries:     make([]interface{}, len(queries)),
			Formulas:    make([]interface{}, len(formulas)),
		}

		// Convert queries and formulas to interface{} for marshaling
		for i := range queries {
			result.Queries[i] = queries[i]
		}
		for i := range formulas {
			result.Formulas[i] = formulas[i]
		}

		// Marshal and print the result
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling result: %v", err)
		}
		fmt.Println(string(output))
	}
	return nil
}
