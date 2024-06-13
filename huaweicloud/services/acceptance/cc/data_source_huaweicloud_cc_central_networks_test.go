package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_networks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcCentralNetworks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.updated_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_not_found", "true"),
					resource.TestCheckOutput("state_filter_is_useful", "true"),
					resource.TestCheckOutput("ep_id_filter_is_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcCentralNetworks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cc_central_networks" "test" {
  depends_on = [
    huaweicloud_cc_central_network.test1,
    huaweicloud_cc_central_network.test2,
    huaweicloud_cc_central_network.test3,
  ]
}

locals {
  id                    = data.huaweicloud_cc_central_networks.test.central_networks[0].id
  state                 = data.huaweicloud_cc_central_networks.test.central_networks[1].state
  enterprise_project_id = data.huaweicloud_cc_central_networks.test.central_networks[2].enterprise_project_id
  tags                  = huaweicloud_cc_central_network.test1.tags
}	

data "huaweicloud_cc_central_networks" "filter_by_id" {
  central_network_id = local.id
}

data "huaweicloud_cc_central_networks" "filter_by_name" {
  depends_on = [
    huaweicloud_cc_central_network.test1,
    huaweicloud_cc_central_network.test2,
    huaweicloud_cc_central_network.test3,
  ]

  name = "%[2]s"
}

data "huaweicloud_cc_central_networks" "filter_by_name_not_found" {
  depends_on = [
    huaweicloud_cc_central_network.test1,
    huaweicloud_cc_central_network.test2,
    huaweicloud_cc_central_network.test3,
  ]

  name = "%[2]s_not_found"
}

data "huaweicloud_cc_central_networks" "filter_by_state" {
  state = local.state
}

data "huaweicloud_cc_central_networks" "filter_by_ep_id" {
  enterprise_project_id = local.enterprise_project_id
}

data "huaweicloud_cc_central_networks" "filter_by_tags" {
  tags = local.tags

  depends_on = [
    huaweicloud_cc_central_network.test1,
    huaweicloud_cc_central_network.test2,
    huaweicloud_cc_central_network.test3,
  ]
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_cc_central_networks.filter_by_id.central_networks) > 0 && alltrue(
    [for v in data.huaweicloud_cc_central_networks.filter_by_id.central_networks[*].id : v == local.id]
  )
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cc_central_networks.filter_by_name.central_networks) == 3
}

output "name_filter_not_found" {
  value = length(data.huaweicloud_cc_central_networks.filter_by_name_not_found.central_networks) == 0
}

output "state_filter_is_useful" {
  value = length(data.huaweicloud_cc_central_networks.filter_by_state.central_networks) > 0 && alltrue(
    [for v in data.huaweicloud_cc_central_networks.filter_by_state.central_networks[*].state : v == local.state]
  )
}

locals {
  central_networks = data.huaweicloud_cc_central_networks.filter_by_ep_id.central_networks
}

output "ep_id_filter_is_useful" {
  value = length(local.central_networks) > 0 && alltrue(
    [for v in local.central_networks[*].enterprise_project_id : v == local.enterprise_project_id]
  )
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_cc_central_networks.filter_by_tags.central_networks) >= 1 && alltrue([
    for ct in data.huaweicloud_cc_central_networks.filter_by_tags.central_networks[*].tags : alltrue([
      for k, v in local.tags : ct[k] == v
    ])
  ])
}
`, testAccDatasourceCreateCentralNetwork(name), name)
}

func testAccDatasourceCreateCentralNetwork(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_central_network" "test1" {
  name        = "%[1]s"
  description = "This is an accaptance test"

  tags = {
    k1 = "v1"
  }
}

resource "huaweicloud_cc_central_network" "test2" {
  name        = "%[1]s_a"
  description = "This is an accaptance test"

  tags = {
    k2 = "v2"
  }
}

resource "huaweicloud_cc_central_network" "test3" {
  name        = "%[1]s_b"
  description = "This is an accaptance test"

  tags = {
    k1 = "v1"
    k2 = "v2"
  }
}
`, rName)
}
