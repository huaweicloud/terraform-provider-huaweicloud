package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcGlobalGateways_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_dc_global_gateways.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		bySortAsc   = "data.huaweicloud_dc_global_gateways.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)

		bySortDesc   = "data.huaweicloud_dc_global_gateways.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)

		byFields   = "data.huaweicloud_dc_global_gateways.filter_by_fields"
		dcByFields = acceptance.InitDataSourceCheck(byFields)

		byGlobalGatewayIDs   = "data.huaweicloud_dc_global_gateways.filter_by_global_gateway_ids"
		dcByGlobalGatewayIDs = acceptance.InitDataSourceCheck(byGlobalGatewayIDs)

		byNames   = "data.huaweicloud_dc_global_gateways.filter_by_names"
		dcByNames = acceptance.InitDataSourceCheck(byNames)

		byEnterpriseProjectIDs   = "data.huaweicloud_dc_global_gateways.filter_by_enterprise_project_ids"
		dcByEnterpriseProjectIDs = acceptance.InitDataSourceCheck(byEnterpriseProjectIDs)

		byStatuses   = "data.huaweicloud_dc_global_gateways.filter_by_statuses"
		dcByStatuses = acceptance.InitDataSourceCheck(byStatuses)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcGlobalGateways_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.#"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.address_family"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.available_peer_link_count"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.bgp_asn"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.current_peer_link_count"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "gateways.0.updated_time"),

					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),

					dcByFields.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_fields_is_useful", "true"),

					dcByGlobalGatewayIDs.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_global_gateway_ids_is_useful", "true"),

					dcByNames.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_names_is_useful", "true"),

					dcByEnterpriseProjectIDs.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_enterprise_project_ids_is_useful", "true"),

					dcByStatuses.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_statuses_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcGlobalGateways_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dc_global_gateway" "test1" {
  name           = "%[1]s_1"
  description    = "test description"
  bgp_asn        = 10
  address_family = "ipv4"
}

resource "huaweicloud_dc_global_gateway" "test2" {
  name           = "%[1]s_2"
  description    = "test description"
  bgp_asn        = 10
  address_family = "ipv4"
}
`, name)
}

func testDataSourceDcGlobalGateways_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_global_gateways" "test" {
  depends_on = [
    huaweicloud_dc_global_gateway.test1,
    huaweicloud_dc_global_gateway.test2,
  ]
}

# filter by sort
data "huaweicloud_dc_global_gateways" "filter_by_sort_asc" {
  depends_on = [
    huaweicloud_dc_global_gateway.test1,
    huaweicloud_dc_global_gateway.test2,
  ]

  sort_key = "id"
  sort_dir = "asc"
}

data "huaweicloud_dc_global_gateways" "filter_by_sort_desc" {
  depends_on = [
    huaweicloud_dc_global_gateway.test1,
    huaweicloud_dc_global_gateway.test2,
  ]

  sort_key = "id"
  sort_dir = "desc"
}

locals {
  gateways_length = length(data.huaweicloud_dc_global_gateways.filter_by_sort_asc.gateways)
  asc_first_id = data.huaweicloud_dc_global_gateways.filter_by_sort_asc.gateways[0].id
  desc_last_id = data.huaweicloud_dc_global_gateways.filter_by_sort_desc.gateways[local.gateways_length - 1].id
}

output "sort_filter_is_useful" {
  value = local.asc_first_id == local.desc_last_id
}

# filter by fields
data "huaweicloud_dc_global_gateways" "filter_by_fields" {
  depends_on = [
    huaweicloud_dc_global_gateway.test1,
  ]

  fields = ["name", "id", "description"]
}

output "filter_by_fields_is_useful" {
  value = length(data.huaweicloud_dc_global_gateways.filter_by_fields.gateways) > 0
}

# filter by global_gateway_ids
locals {
  id = data.huaweicloud_dc_global_gateways.test.gateways[0].id
}

data "huaweicloud_dc_global_gateways" "filter_by_global_gateway_ids" {
  global_gateway_ids = [local.id]
}

locals {
  filter_by_global_gateway_ids_result = [
    for v in data.huaweicloud_dc_global_gateways.filter_by_global_gateway_ids.gateways[*].id : v == local.id
  ]
}

output "filter_by_global_gateway_ids_is_useful" {
  value = alltrue(local.filter_by_global_gateway_ids_result) && length(local.filter_by_global_gateway_ids_result) > 0
}

# filter by names
locals {
  name = data.huaweicloud_dc_global_gateways.test.gateways[0].name
}

data "huaweicloud_dc_global_gateways" "filter_by_names" {
  names = [local.name]
}

locals {
  filter_by_names_result = [
    for v in data.huaweicloud_dc_global_gateways.filter_by_names.gateways[*].name : v == local.name
  ]
}

output "filter_by_names_is_useful" {
  value = alltrue(local.filter_by_names_result) && length(local.filter_by_names_result) > 0
}

# filter by enterprise_project_ids
locals {
  enterprise_project_id = data.huaweicloud_dc_global_gateways.test.gateways[0].enterprise_project_id
}

data "huaweicloud_dc_global_gateways" "filter_by_enterprise_project_ids" {
  enterprise_project_ids = [local.enterprise_project_id]
}

locals {
  filter_by_enterprise_project_ids_result = [
    for v in data.huaweicloud_dc_global_gateways.filter_by_enterprise_project_ids.gateways[*].enterprise_project_id :
    v == local.enterprise_project_id
  ]
}

output "filter_by_enterprise_project_ids_is_useful" {
  value = alltrue(local.filter_by_enterprise_project_ids_result) && length(local.filter_by_enterprise_project_ids_result) > 0
}

# filter by statuses
locals {
  status = data.huaweicloud_dc_global_gateways.test.gateways[0].status
}

data "huaweicloud_dc_global_gateways" "filter_by_statuses" {
  statuses = [local.status]
}

locals {
  filter_by_statuses_result = [
    for v in data.huaweicloud_dc_global_gateways.filter_by_statuses.gateways[*].status : v == local.status
  ]
}

output "filter_by_statuses_is_useful" {
  value = alltrue(local.filter_by_statuses_result) && length(local.filter_by_statuses_result) > 0
}
`, testDataSourceDcGlobalGateways_base(name))
}
