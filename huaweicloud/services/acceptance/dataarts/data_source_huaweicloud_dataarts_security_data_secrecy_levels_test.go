package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityDataSecrecyLevels_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_data_secrecy_levels.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byOrderByAsc    = "data.huaweicloud_dataarts_security_data_secrecy_levels.filter_by_order_by_asc"
		dcByOrderByAsc  = acceptance.InitDataSourceCheck(byOrderByAsc)
		byOrderByDesc   = "data.huaweicloud_dataarts_security_data_secrecy_levels.filter_by_order_by_desc"
		dcByOrderByDesc = acceptance.InitDataSourceCheck(byOrderByDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityDataSecrecyLevels_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestMatchResourceAttr(all, "secrecy_levels.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.id"),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.name"),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.level_number"),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.created_by"),
					resource.TestMatchResourceAttr(all, "secrecy_levels.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "secrecy_levels.0.updated_by"),
					resource.TestMatchResourceAttr(all, "secrecy_levels.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					// Filter by 'order_by' parameter.
					dcByOrderByAsc.CheckResourceExists(),
					dcByOrderByDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_order_by_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityDataSecrecyLevels_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = "Created by terraform"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataSecurityDataSecrecyLevels_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_security_data_secrecy_levels" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_security_data_secrecy_level.test]
}

# Filter by 'order_by' parameter to get the data secrecy levels in ascending order
data "huaweicloud_dataarts_security_data_secrecy_levels" "filter_by_order_by_asc" {
  workspace_id = "%[2]s"
  order_by     = "name"

  depends_on = [huaweicloud_dataarts_security_data_secrecy_level.test]
}


# Filter by 'order_by' parameter to get the data secrecy levels in descending order
data "huaweicloud_dataarts_security_data_secrecy_levels" "filter_by_order_by_desc" {
  workspace_id = "%[2]s"
  order_by     = "name"
  desc         = true

  depends_on = [huaweicloud_dataarts_security_data_secrecy_level.test]
}

locals {
  secrecy_level_names_desc = data.huaweicloud_dataarts_security_data_secrecy_levels.filter_by_order_by_desc.secrecy_levels[*].name
  secrecy_level_names_asc  = data.huaweicloud_dataarts_security_data_secrecy_levels.filter_by_order_by_asc.secrecy_levels[*].name
}

output "is_order_by_filter_useful" {
  value = local.secrecy_level_names_desc == reverse(local.secrecy_level_names_asc)
}
`, testAccDataSecurityDataSecrecyLevels_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
