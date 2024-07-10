package setting

import "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"

type ServiceOptions struct {
	Get   func(args ...string) string
	List  func(args ...string) string
	Stubs api.RecordStubs
}
