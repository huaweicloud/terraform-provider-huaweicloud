package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCustomLines_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dns_custom_line.test"

		dataSource = "data.huaweicloud_dns_custom_lines.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLineId           = "data.huaweicloud_dns_custom_lines.filter_by_line_id"
		dcByLineId         = acceptance.InitDataSourceCheck(byLineId)
		byNotFoundLineId   = "data.huaweicloud_dns_custom_lines.filter_by_not_found_line_id"
		dcByNotFoundLineId = acceptance.InitDataSourceCheck(byNotFoundLineId)

		byName           = "data.huaweicloud_dns_custom_lines.filter_by_name"
		dcByName         = acceptance.InitDataSourceCheck(byName)
		byNameFuzzy      = "data.huaweicloud_dns_custom_lines.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNotFoundName   = "data.huaweicloud_dns_custom_lines.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byIp           = "data.huaweicloud_dns_custom_lines.filter_by_ip"
		dcByIp         = acceptance.InitDataSourceCheck(byIp)
		byNotFoundIp   = "data.huaweicloud_dns_custom_lines.filter_by_not_found_ip"
		dcByNotFoundIp = acceptance.InitDataSourceCheck(byNotFoundIp)

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
				Config: testAccDataCustomLines_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "lines.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by custom line ID.
					dcByLineId.CheckResourceExists(),
					resource.TestCheckOutput("is_line_id_filter_useful", "true"),
					dcByNotFoundLineId.CheckResourceExists(),
					resource.TestCheckOutput("line_id_not_found_validation_pass", "true"),
					// Exact search by custom line name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Fuzzy search by custom line name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by IP address.
					dcByIp.CheckResourceExists(),
					resource.TestCheckOutput("is_ip_filter_useful", "true"),
					dcByNotFoundIp.CheckResourceExists(),
					resource.TestCheckOutput("ip_not_found_validation_pass", "true"),
					// Filter by custom line status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byLineId, "lines.0.description", rName, "description"),
					resource.TestCheckResourceAttrPair(byLineId, "lines.0.ip_segments", rName, "ip_segments"),
					// The format is the time corresponding to the local time zone of the computer.
					resource.TestMatchResourceAttr(byLineId, "lines.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byLineId, "lines.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataCustomLines_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_custom_line" "test" {
  name        = "%s"
  description = "test description"
  # IP address ranges cannot overlap.
  ip_segments = ["1.0.0.1-2.1.1.2"]
}

data "huaweicloud_dns_custom_lines" "test" {
  depends_on = [huaweicloud_dns_custom_line.test]
}

# Filter by custom line ID.
locals {
  line_id = huaweicloud_dns_custom_line.test.id
}

data "huaweicloud_dns_custom_lines" "filter_by_line_id" {
  line_id = local.line_id
}

locals {
  line_id_filter_result = [for v in data.huaweicloud_dns_custom_lines.filter_by_line_id.lines[*].id : v == local.line_id]
}

output "is_line_id_filter_useful" {
  value = length(local.line_id_filter_result) > 0 && alltrue(local.line_id_filter_result)
}

data "huaweicloud_dns_custom_lines" "filter_by_not_found_line_id" {
  depends_on = [huaweicloud_dns_custom_line.test]
  line_id    = "not_found_line_id"
}

output "line_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_custom_lines.filter_by_not_found_line_id.lines) == 0
}

# Exact search by custom line name.
locals {
  line_name = huaweicloud_dns_custom_line.test.name
}

data "huaweicloud_dns_custom_lines" "filter_by_name" {
  depends_on = [huaweicloud_dns_custom_line.test]

  name = local.line_name
}

locals {
  name_filter_result = [for v in data.huaweicloud_dns_custom_lines.filter_by_name.lines[*].name : v == local.line_name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Fuzzy search by custom line name.
locals {
  name_prefix = "tf_test"
}

data "huaweicloud_dns_custom_lines" "filter_by_name_fuzzy" {
  depends_on = [huaweicloud_dns_custom_line.test]

  name = local.name_prefix
}

locals {
  name_fuzzy_filter_result = [for v in data.huaweicloud_dns_custom_lines.filter_by_name_fuzzy.lines[*].name :
    strcontains(v, local.name_prefix)
  ]
}

output "is_name_fuzzy_filter_useful" {
  value = length(local.name_fuzzy_filter_result) > 0 && alltrue(local.name_fuzzy_filter_result)
}

data "huaweicloud_dns_custom_lines" "filter_by_not_found_name" {
  depends_on = [huaweicloud_dns_custom_line.test]
  name       = "not_found_custom_name"
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_custom_lines.filter_by_not_found_name.lines) == 0
}

# Filter by IP address.
locals {
  ip = try(split("-", tolist(huaweicloud_dns_custom_line.test.ip_segments)[0])[0], "")
}

data "huaweicloud_dns_custom_lines" "filter_by_ip" {
  depends_on = [huaweicloud_dns_custom_line.test]

  ip = local.ip
}

locals {
  ip_filter_result = [for v in data.huaweicloud_dns_custom_lines.filter_by_ip.lines[*].ip_segments : strcontains(join(",", v), local.ip)]
}

output "is_ip_filter_useful" {
  value = length(local.ip_filter_result) > 0 && alltrue(local.ip_filter_result)
}

data "huaweicloud_dns_custom_lines" "filter_by_not_found_ip" {
  depends_on = [huaweicloud_dns_custom_line.test]

  ip = "0.0.0.0"
}

output "ip_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_custom_lines.filter_by_not_found_ip.lines) == 0
}

# Filter by custom line status.
locals {
  status = huaweicloud_dns_custom_line.test.status
}

data "huaweicloud_dns_custom_lines" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [for v in data.huaweicloud_dns_custom_lines.filter_by_status.lines[*].status : v == local.status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, name)
}
