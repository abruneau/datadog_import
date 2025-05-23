# Dynatrace

## Synthetics

| Source  | Support |
| ------- | ------- |
| browser | ✅      |
| api     | ✅      |

### HTTP tests

- Dynatrace allows checking if the status code is greater or less than a value. Datadog doesn't allow it and only allow is equal, not equal, matches, not match
- Datadog doesn't support pre or post processing scripts

## Configuration

```yaml
# Dynatrace Configuration
dynatrace:
  # Dynatrace URL
  url: ""
  # Dynatrace API Key
  api_key: ""
  # Input directory containing Dynatrace synthetics tests definitions
  input: ""

  # Filters to apply to the list of tests
  filters:
    # Filters the resulting set of monitors to those which are part of the specified management zone ID.
    management_zone: ""
    # Filters the resulting set of monitors by specified tags.
    tags: []
    # Filters the resulting set of monitors by specified location ID.
    location: ""
    # Filters the resulting set of monitors to those of the specified type: BROWSER or HTTP.
    type: ""
    # Filters the resulting set of monitors to those which are enabled (true) or disabled (false)
    enabled: ""
    # Filters the resulting set of monitors to those using the specified credential ID.
    credential_id: ""
    # Filters the resulting set of monitors to those using the specified credential owner.
    credential_owner: ""
    # Filters the resulting set of monitors to those assigned to the specified application IDs.
    assigned_apps: ""

  # List of test IDs to fetch
  id_list: []
  # List of custom tags to add to the tests
  custom_tags: []
```
