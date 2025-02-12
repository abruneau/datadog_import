# Grafana

## Compatibility

| Source      | Support |
| ----------- | ------- |
| Stackdriver | ✅      |
| Cloudwatch  | ✅      |
| Azure       | ✅      |
| Prometheus  | ✅      |
| Loki        | ✅      |

## Configuration

```yaml
# Grafana Configuration
grafana:
  # Input directory containing Grafana dashboards definitions
  input: ""

  # API configuration
  api:
    # Host is the domain name or IP address of the host that serves the API.
    host: ""
    # BasePath is the URL prefix for all API paths, relative to the host root.
    base_path: ""
    # Schemes are the transfer protocols used by the API (http or https).
    schemes: []
    # APIKey is an optional API key or service account token.
    api_key: ""
    basic_auth:
      # Username for basic auth credentials.
      username: ""
      # Password for basic auth credentials.
      password: ""

    # OrgID provides an optional organization ID.
    org_id: 0

  # Filters to apply to the list of dashboards
  filters:
    # List of dashboard uid’s to search for
    dashboard_uids: []
    # Flag indicating if only soft deleted Dashboards should be returned
    deleted: false
    # List of folder UID’s to search in for dashboards
    folder_uids: []
    # Limit the number of returned dashboards (max 5000)
    limit: 0
    # Search Query
    query: ""
    # Flag indicating if only starred Dashboards should be returned
    starred: false
    # List of tags to search for
    tag: []
    # Type of dashboards to search for, dash-folder or dash-db
    type: ""
```
