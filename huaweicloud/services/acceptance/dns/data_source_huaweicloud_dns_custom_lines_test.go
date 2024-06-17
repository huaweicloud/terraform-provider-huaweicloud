package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCustomLines_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_dns_custom_lines.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byId   = "data.huaweicloud_dns_custom_lines.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_dns_custom_lines.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_dns_custom_lines.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDnsCustomLines_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "lines.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "lines.0.ip_segments.#"),
					// The format is the time corresponding to the local time zone of the computer.
					resource.TestMatchResourceAttr(dataSource, "lines.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "lines.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDnsCustomLines_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s"
  description = "test description"
  # IP address ranges cannot overlap.
  ip_segments = ["1.0.0.1-2.1.1.2"]
}

data "huaweicloud_dns_custom_lines" "test" {
  depends_on = [
    huaweicloud_dns_custom_line.test
  ]
}

locals {
  line_id = data.huaweicloud_dns_custom_lines.test.lines[0].id
}

data "huaweicloud_dns_custom_lines" "filter_by_id" {
  line_id = local.line_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_dns_custom_lines.filter_by_id.lines[*].id : v == local.line_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

locals {
  line_name = data.huaweicloud_dns_custom_lines.test.lines[0].name
}

data "huaweicloud_dns_custom_lines" "filter_by_name" {
  name = local.line_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dns_custom_lines.filter_by_name.lines[*].name : v == local.line_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

locals {
  not_found_name = "not_found"
}

data "huaweicloud_dns_custom_lines" "filter_by_not_found_name" {
  name = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_dns_custom_lines.filter_by_not_found_name.lines[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

locals {
  status = data.huaweicloud_dns_custom_lines.test.lines[0].status
}

data "huaweicloud_dns_custom_lines" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dns_custom_lines.filter_by_status.lines[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, name)
}
