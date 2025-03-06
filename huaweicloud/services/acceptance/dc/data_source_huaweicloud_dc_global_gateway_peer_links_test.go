package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcGlobalGatewayPeerLinks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dc_global_gateway_peer_links.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		bySort   = "data.huaweicloud_dc_global_gateway_peer_links.filter_by_sort"
		dcBySort = acceptance.InitDataSourceCheck(bySort)

		byFields   = "data.huaweicloud_dc_global_gateway_peer_links.filter_by_fields"
		dcByFields = acceptance.InitDataSourceCheck(byFields)

		byPeerLinkIDs   = "data.huaweicloud_dc_global_gateway_peer_links.filter_by_peer_link_ids"
		dcByPeerLinkIDs = acceptance.InitDataSourceCheck(byPeerLinkIDs)

		byNames   = "data.huaweicloud_dc_global_gateway_peer_links.filter_by_names"
		dcByNames = acceptance.InitDataSourceCheck(byNames)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please configure the global gateway ID containing peer link.
			acceptance.TestAccPreCheckDcGlobalGatewayIDHasPeerLink(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcGlobalGatewayPeerLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.bandwidth_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.create_owner"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.global_dc_gateway_id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.#"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.gateway_id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.link_id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.peer_site.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "peer_links.0.updated_time"),

					dcBySort.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_sort_is_useful", "true"),

					dcByFields.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_fields_is_useful", "true"),

					dcByPeerLinkIDs.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_peer_link_ids_is_useful", "true"),

					dcByNames.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_names_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDcGlobalGatewayPeerLinks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dc_global_gateway_peer_links" "test" {
  global_dc_gateway_id = "%[1]s"
}

# filter by sort (Just determine whether the data exists)
data "huaweicloud_dc_global_gateway_peer_links" "filter_by_sort" {
  global_dc_gateway_id = "%[1]s"
  sort_key             = "id"
  sort_dir             = "desc"
}

output "filter_by_sort_is_useful" {
  value = length(data.huaweicloud_dc_global_gateway_peer_links.filter_by_sort.peer_links) > 0
}

# filter by fields (Just determine whether the data exists)
data "huaweicloud_dc_global_gateway_peer_links" "filter_by_fields" {
  global_dc_gateway_id = "%[1]s"
  fields               = ["name", "id"]
}

output "filter_by_fields_is_useful" {
  value = length(data.huaweicloud_dc_global_gateway_peer_links.filter_by_fields.peer_links) > 0
}

# filter by peer_link_ids
locals {
  id = data.huaweicloud_dc_global_gateway_peer_links.test.peer_links[0].id
}

data "huaweicloud_dc_global_gateway_peer_links" "filter_by_peer_link_ids" {
  global_dc_gateway_id = "%[1]s"
  peer_link_ids        = [local.id]
}

locals {
  filter_by_peer_link_ids_result = [
    for v in data.huaweicloud_dc_global_gateway_peer_links.filter_by_peer_link_ids.peer_links[*].id : v == local.id
  ]
}

output "filter_by_peer_link_ids_is_useful" {
  value = alltrue(local.filter_by_peer_link_ids_result) && length(local.filter_by_peer_link_ids_result) > 0
}

# filter by names
locals {
  name = data.huaweicloud_dc_global_gateway_peer_links.test.peer_links[0].name
}

data "huaweicloud_dc_global_gateway_peer_links" "filter_by_names" {
  global_dc_gateway_id = "%[1]s"
  names                = [local.name]
}

locals {
  filter_by_names_result = [
    for v in data.huaweicloud_dc_global_gateway_peer_links.filter_by_names.peer_links[*].name : v == local.name
  ]
}

output "filter_by_names_is_useful" {
  value = alltrue(local.filter_by_names_result) && length(local.filter_by_names_result) > 0
}
`, acceptance.HW_DC_GLOBAL_GATEWAY_ID_HAS_PEER_LINK)
}
