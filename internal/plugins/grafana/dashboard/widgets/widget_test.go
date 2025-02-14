package widgets

import (
	"testing"

	"context"
	"datadog_import/internal/plugins/grafana/dashboard/types"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

func TestConvertGrid(t *testing.T) {
	tests := []struct {
		name     string
		panel    types.Panel
		expected datadogV1.WidgetLayout
	}{
		{
			name: "Normal case",
			panel: types.Panel{
				GridPos: struct {
					H int `json:"h"`
					W int `json:"w"`
					X int `json:"x"`
					Y int `json:"y"`
				}{
					H: 4,
					W: 6,
					X: 2,
					Y: 2,
				},
			},
			expected: *datadogV1.NewWidgetLayout(2, 3, 1, 1),
		},
		{
			name: "Height zero case",
			panel: types.Panel{
				GridPos: struct {
					H int `json:"h"`
					W int `json:"w"`
					X int `json:"x"`
					Y int `json:"y"`
				}{
					H: 1,
					W: 6,
					X: 2,
					Y: 2,
				},
			},
			expected: *datadogV1.NewWidgetLayout(1, 3, 1, 1),
		},
		{
			name: "Odd dimensions case",
			panel: types.Panel{
				GridPos: struct {
					H int `json:"h"`
					W int `json:"w"`
					X int `json:"x"`
					Y int `json:"y"`
				}{
					H: 5,
					W: 7,
					X: 3,
					Y: 3,
				},
			},
			expected: *datadogV1.NewWidgetLayout(2, 3, 1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := NewPanelConvertor(context.Background(), "source", tt.panel)
			result := pc.convertGrid()

			if result.GetHeight() != tt.expected.GetHeight() ||
				result.GetWidth() != tt.expected.GetWidth() ||
				result.GetX() != tt.expected.GetX() ||
				result.GetY() != tt.expected.GetY() {
				t.Errorf("convertGrid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
