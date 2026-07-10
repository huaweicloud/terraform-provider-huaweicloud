package taurusdb

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBDedicatedResourceDetails_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBDedicatedResourceDetails_basic(),
				// No available dedicated resource, expect error with 'Server failure' using invalid id.
				ExpectError: regexp.MustCompile(`Server failure`),
			},
		},
	})
}

func testAccDataSourceTaurusDBDedicatedResourceDetails_basic() string {
	return `
data "huaweicloud_taurusdb_dedicated_resource" "test" {}

data "huaweicloud_taurusdb_dedicated_resource_details" "test" {
  dedicated_resource_id = data.huaweicloud_taurusdb_dedicated_resource.test.id
}
`
}
