package eps

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectQuotas_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_quotas.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_eps_quota_configured", "true"),
					resource.TestCheckOutput("is_type_configured", "true"),
					resource.TestCheckOutput("is_quota_used_number_configured", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectQuotas_basic = `
data "huaweicloud_enterprise_project_quotas" "test" {}

locals {
  enterprise_project_quota = try(element([for o in data.huaweicloud_enterprise_project_quotas.test.quotas: o if
    o.type == "enterprise_project"], 0), null)
}

output "is_eps_quota_configured" {
  value = local.enterprise_project_quota != null
}

output "is_total_quota_number_configured" {
  value = lookup(local.enterprise_project_quota, "quota") > 0
}

output "is_type_configured" {
  value = lookup(local.enterprise_project_quota, "type") != ""
}

output "is_quota_used_number_configured" {
  value = lookup(local.enterprise_project_quota, "used") > 0
}
`
