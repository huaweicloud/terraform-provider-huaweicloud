package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataConnections_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_eg_connections.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_eg_connections.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byFuzzyName   = "data.huaweicloud_eg_connections.filter_by_fuzzy_name"
		dcByFuzzyName = acceptance.InitDataSourceCheck(byFuzzyName)

		bySortDesc   = "data.huaweicloud_eg_connections.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)

		bySortAsc   = "data.huaweicloud_eg_connections.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgConnectionIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.id"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.name"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.vpc_id"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.subnet_id"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.type"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.status"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.created_time"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.updated_time"),
					resource.TestCheckResourceAttr(byName, "connections.0.flavor.#", "1"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.flavor.0.name"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.flavor.0.concurrency_type"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.flavor.0.concurrency"),
					resource.TestCheckResourceAttrSet(byName, "connections.0.flavor.0.bandwidth_type"),
					dcByFuzzyName.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_name_filter_useful", "true"),
					dcBySortDesc.CheckResourceExists(),
					dcBySortAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
					resource.TestCheckOutput("is_kafka_detail_and_valid", "true"),
					// `description` and `error_info` may be empty, so we don't check them.
				),
			},
		},
	})
}

func testAccDataConnections_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_eg_connections" "test" {}

locals {
  connection_ids = "%[1]s"
  filter_result  = try([for v in data.huaweicloud_eg_connections.test.connections : v if v.id == split(",", local.connection_ids)[0]][0], {})
  name           = lookup(local.filter_result, "name", "")
  kafka_connection_filter_result = try(
    [for v in data.huaweicloud_eg_connections.test.connections : v if v.id == split(",", local.connection_ids)[1]][0].kafka_detail[0],
    {}
  )
}

data "huaweicloud_eg_connections" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [for v in data.huaweicloud_eg_connections.filter_by_name.connections[*].name : v == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

data "huaweicloud_eg_connections" "filter_by_fuzzy_name" {
  fuzzy_name = local.name
}

locals {
  fuzzy_name_filter_result = [for v in data.huaweicloud_eg_connections.filter_by_fuzzy_name.connections[*].name :
  strcontains(v, local.name)]
}

output "is_fuzzy_name_filter_useful" {
  value = length(local.fuzzy_name_filter_result) > 0 && alltrue(local.fuzzy_name_filter_result)
}

data "huaweicloud_eg_connections" "filter_by_sort_desc" {
  sort = "name:DESC"
}

data "huaweicloud_eg_connections" "filter_by_sort_asc" {
  sort = "name:ASC"
}

locals {
  sort_desc_filter_result = data.huaweicloud_eg_connections.filter_by_sort_desc.connections[*].name
  sort_asc_filter_result  = data.huaweicloud_eg_connections.filter_by_sort_asc.connections[*].name
}

output "is_sort_filter_useful" {
  value = (
    length(local.sort_desc_filter_result) == length(local.sort_asc_filter_result) &&
    local.sort_desc_filter_result == reverse(local.sort_asc_filter_result)
  )
}

output "is_kafka_detail_and_valid" {
  value = alltrue([
    lookup(local.kafka_connection_filter_result, "instance_id", "") != "",
    lookup(local.kafka_connection_filter_result, "connect_address", "") != "",
    lookup(local.kafka_connection_filter_result, "security_protocol", "") != "",
    lookup(local.kafka_connection_filter_result, "enable_sasl_ssl", "") == true,
    lookup(local.kafka_connection_filter_result, "user_name", "") != "",
    lookup(local.kafka_connection_filter_result, "acks", "") != "",
    lookup(local.kafka_connection_filter_result, "address", "") != "",
  ])
}
`, acceptance.HW_EG_CONNECTION_IDS)
}
