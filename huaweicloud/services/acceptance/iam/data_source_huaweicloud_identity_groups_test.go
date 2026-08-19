package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3Groups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_groups.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_identity_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byDomainId   = "data.huaweicloud_identity_groups.filter_by_domain_id"
		dcByDomainId = acceptance.InitDataSourceCheck(byDomainId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3Groups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "groups.0.id", "huaweicloud_identity_group.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "groups.0.name", "huaweicloud_identity_group.test", "name"),
					resource.TestCheckResourceAttrPair(byName, "groups.0.description", "huaweicloud_identity_group.test", "description"),
					resource.TestCheckResourceAttrSet(byName, "groups.0.domain_id"),
					resource.TestMatchResourceAttr(byName, "groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by domain ID.
					dcByDomainId.CheckResourceExists(),
					resource.TestCheckOutput("is_domain_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV3Groups_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name        = "%[1]s"
  description = "created by terraform script"
}

data "huaweicloud_identity_groups" "test" {
  depends_on = [huaweicloud_identity_group.test]
}

# Filter by name parameter.
locals {
  group_name = huaweicloud_identity_group.test.name
}

data "huaweicloud_identity_groups" "filter_by_name" {
  depends_on = [huaweicloud_identity_group.test]

  name = local.group_name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_identity_groups.filter_by_name.groups) > 0 && alltrue(
    [for v in data.huaweicloud_identity_groups.filter_by_name.groups[*].name : v == local.group_name]
  )
}

# Filter by domain ID parameter.
locals {
  domain_id = "%[2]s"
}

data "huaweicloud_identity_groups" "filter_by_domain_id" {
  depends_on = [huaweicloud_identity_group.test]

  domain_id = local.domain_id
}

output "is_domain_id_filter_useful" {
  value = length(data.huaweicloud_identity_groups.filter_by_domain_id.groups) > 0 && alltrue(
    [for v in data.huaweicloud_identity_groups.filter_by_domain_id.groups[*].domain_id : v == local.domain_id]
  )
}
`, name, acceptance.HW_DOMAIN_ID)
}
