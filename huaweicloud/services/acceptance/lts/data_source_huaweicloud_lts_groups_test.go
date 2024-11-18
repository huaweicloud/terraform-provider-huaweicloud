package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGroups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_groups.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.ttl_in_days"),
					resource.TestMatchResourceAttr(dataSource, "groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_exist_log_group", "true"),
					resource.TestCheckOutput("is_eps_return_and_matched", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lts_groups" "test" {
  depends_on = [
    huaweicloud_lts_group.test
  ]
}

output "is_exist_log_group" {
  value = contains(data.huaweicloud_lts_groups.test.groups[*].id, huaweicloud_lts_group.test.id)
}

locals {
  eps_filter_result = [for v in data.huaweicloud_lts_groups.test.groups :
  v.enterprise_project_id == huaweicloud_lts_group.test.enterprise_project_id if v.id == huaweicloud_lts_group.test.id]
}

output "is_eps_return_and_matched" {
  value = length(local.eps_filter_result) > 0 && alltrue(local.eps_filter_result)
}
`, testAccLtsGroup_basic(name, 30))
}
