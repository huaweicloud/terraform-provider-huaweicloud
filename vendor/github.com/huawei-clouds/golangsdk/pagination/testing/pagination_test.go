package testing

import (
	"github.com/huawei-clouds/golangsdk"
	"github.com/huawei-clouds/golangsdk/testhelper"
)

func createClient() *golangsdk.ServiceClient {
	return &golangsdk.ServiceClient{
		ProviderClient: &golangsdk.ProviderClient{TokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}
