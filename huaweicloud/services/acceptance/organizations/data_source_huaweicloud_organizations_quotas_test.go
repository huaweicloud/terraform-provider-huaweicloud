package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQuotas_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_quotas.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQuotas_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "quotas.0.resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "quotas.0.resources.0.type"),
					resource.TestCheckOutput("is_organizational_unit_used", "true"),
					resource.TestCheckOutput("is_quotas_max_set", "true"),
					resource.TestCheckOutput("is_quotas_min_set", "true"),
					resource.TestCheckOutput("is_quotas_quota_set", "true"),
					resource.TestCheckOutput("is_quotas_used_set", "true"),
				),
			},
		},
	})
}

func testAccDataQuotas_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}
`, name)
}

func testAccDataQuotas_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_quotas" "test" {
  depends_on = [huaweicloud_organizations_organizational_unit.test]
}

locals {
  all_used_quotas = [for v in flatten(data.huaweicloud_organizations_quotas.test.quotas[*].resources[*]) : v if v.used > 0]
}

output "is_organizational_unit_used" {
  value = length(local.all_used_quotas) > 0 && contains(local.all_used_quotas[*].type, "organizational_unit")
}

output "is_quotas_max_set" {
  value = length(local.all_used_quotas) > 0 && local.all_used_quotas[0].max > 0
}

output "is_quotas_min_set" {
  value = length(local.all_used_quotas) > 0 && local.all_used_quotas[0].min > 0
}

output "is_quotas_quota_set" {
  value = length(local.all_used_quotas) > 0 && local.all_used_quotas[0].quota > 0
}

output "is_quotas_used_set" {
  value = length(local.all_used_quotas) > 0 && local.all_used_quotas[0].used > 0
}
`, testAccDataQuotas_base(name))
}
