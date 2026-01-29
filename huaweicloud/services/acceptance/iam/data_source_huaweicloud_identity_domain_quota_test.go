package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomainQuota_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_domain_quota.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_identity_domain_quota.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataDomainQuota_basic = `
# All
data "huaweicloud_identity_domain_quota" "all" {}

# Filter by type
data "huaweicloud_identity_domain_quota" "filter_by_type" {
  type = "user"
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_identity_domain_quota.filter_by_type.resources[*].type : v == "user"
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`
