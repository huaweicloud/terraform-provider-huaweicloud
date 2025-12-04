package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataOus_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_ous.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByOuName   = "data.huaweicloud_workspace_ous.filter_by_ou_name"
		dcFilterByOuName = acceptance.InitDataSourceCheck(filterByOuName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// Without any filter parameter.
					resource.TestMatchResourceAttr(all, "ous.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "ous.0.id"),
					resource.TestCheckResourceAttrSet(all, "ous.0.name"),
					resource.TestCheckResourceAttrSet(all, "ous.0.domain_id"),
					resource.TestCheckResourceAttrSet(all, "ous.0.domain"),
					resource.TestCheckResourceAttrSet(all, "ous.0.ou_dn"),
					// Filter by 'ou_name' parameter.
					dcFilterByOuName.CheckResourceExists(),
					resource.TestCheckOutput("is_ou_name_filter_useful", "true"),
					// `description` may be empty, so we don't check it.
				),
			},
		},
	})
}

func testAccDataOus_basic() string {
	return fmt.Sprintf(`
# Without any filter parameter.
data "huaweicloud_workspace_ous" "all" {}

# Filter by 'ou_name' parameter.
data "huaweicloud_workspace_ous" "filter_by_ou_name" {
  ou_name = "%[1]s"
}

locals {
  ou_name_filter_result = [for v in data.huaweicloud_workspace_ous.filter_by_ou_name.ous[*].name :
    strcontains(v, "%[1]s")
  ]
}

output "is_ou_name_filter_useful" {
  value = length(local.ou_name_filter_result) > 0 && alltrue(local.ou_name_filter_result)
}
`, acceptance.HW_WORKSPACE_OU_NAME)
}
