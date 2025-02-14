package dynatrace

import (
	"testing"
)

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		name     string
		filters  Filters
		expected string
	}{
		{
			name:     "Empty filters",
			filters:  Filters{},
			expected: "",
		},
		{
			name: "ManagementZone filter",
			filters: Filters{
				ManagementZone: "zone1",
			},
			expected: "managementZone:zone1",
		},
		{
			name: "Tags filter",
			filters: Filters{
				Tags: []string{"tag1", "tag2"},
			},
			expected: "tag=tag1&tag=tag2",
		},
		{
			name: "Location filter",
			filters: Filters{
				Location: "location1",
			},
			expected: "location:location1",
		},
		{
			name: "Type filter",
			filters: Filters{
				Type: "BROWSER",
			},
			expected: "type:BROWSER",
		},
		{
			name: "Enabled filter",
			filters: Filters{
				Enabled: "true",
			},
			expected: "enabled:true",
		},
		{
			name: "CredentialId filter",
			filters: Filters{
				CredentialId: "cred1",
			},
			expected: "credentialId:cred1",
		},
		{
			name: "CredentialOwner filter",
			filters: Filters{
				CredentialOwner: "owner1",
			},
			expected: "credentialOwner:owner1",
		},
		{
			name: "AssignedApps filter",
			filters: Filters{
				AssignedApps: "app1",
			},
			expected: "assignedApps:app1",
		},
		{
			name: "Multiple filters",
			filters: Filters{
				ManagementZone:  "zone1",
				Tags:            []string{"tag1", "tag2"},
				Location:        "location1",
				Type:            "BROWSER",
				Enabled:         "true",
				CredentialId:    "cred1",
				CredentialOwner: "owner1",
				AssignedApps:    "app1",
			},
			expected: "managementZone:zone1&tag=tag1&tag=tag2&location:location1&type:BROWSER&enabled:true&credentialId:cred1&credentialOwner:owner1&assignedApps:app1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFilter(tt.filters)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
