package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLineGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_line_groups.filter_by_line_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceLineGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.lines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.updated_at"),

					resource.TestCheckOutput("is_line_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_name", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceLineGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  line_id = huaweicloud_dns_line_group.test.id
}

data "huaweicloud_dns_line_groups" "filter_by_line_id" {
  line_id = local.line_id
}

output "is_line_id_filter_useful" {
  value = length(data.huaweicloud_dns_line_groups.filter_by_line_id.groups) > 0 && alltrue(
    [for v in data.huaweicloud_dns_line_groups.filter_by_line_id.groups[*].id : v == local.line_id]
  )
}

locals {
  name = huaweicloud_dns_line_group.test.name
}

data "huaweicloud_dns_line_groups" "filter_by_name" {
  depends_on = [huaweicloud_dns_line_group.test]
  name       = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_dns_line_groups.filter_by_name.groups) > 0 && alltrue(
    [for v in data.huaweicloud_dns_line_groups.filter_by_name.groups[*].name : v == local.name]
  )
}

data "huaweicloud_dns_line_groups" "filter_not_found_name" {
  depends_on = [huaweicloud_dns_line_group.test]
  name       = "not_found_name"
}

output "not_found_name" {
  value = length(data.huaweicloud_dns_line_groups.filter_not_found_name.groups) == 0
}
`, testDNSLineGroup_basic(name))
}
