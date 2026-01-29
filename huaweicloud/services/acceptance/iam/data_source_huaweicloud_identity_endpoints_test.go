package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEndpoints_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_endpoints.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_identity_endpoints.filter_by_endpoint_id"
		dcById = acceptance.InitDataSourceCheck(byId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEndpoints_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "endpoints.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "endpoints.0.enabled", "true"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_endpoint_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataEndpoints_basic = `
# All
data "huaweicloud_identity_endpoints" "all" {}

# Filter by endpoint ID
locals {
  endpoint_id = try(data.huaweicloud_identity_endpoints.all.endpoints.0.id, "NOT_FOUND")
}

data "huaweicloud_identity_endpoints" "filter_by_endpoint_id" {
  endpoint_id = local.endpoint_id
}

locals {
  endpoint_id_filter_result = [
    for v in data.huaweicloud_identity_endpoints.filter_by_endpoint_id.endpoints[*].id : v == local.endpoint_id
  ]
}

output "is_endpoint_id_filter_useful" {
  value = length(local.endpoint_id_filter_result) > 0 && alltrue(local.endpoint_id_filter_result)
}
`
