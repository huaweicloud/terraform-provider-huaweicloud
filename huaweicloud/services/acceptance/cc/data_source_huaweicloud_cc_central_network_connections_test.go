package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkConnections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_connections.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCProjectID(t)
			acceptance.TestAccPreCheckCCRegionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCentralNetworkConnections_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.bandwidth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_connections.0.updated_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_type_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("cross_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcCentralNetworkConnections_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cc_central_network_connections" "test" {
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  central_network_id = huaweicloud_cc_central_network.test.id
}
  
locals {
  central_network_connections = data.huaweicloud_cc_central_network_connections.test.central_network_connections
  id                          = local.central_network_connections[0].id
  status                      = local.central_network_connections[0].status
  bandwidth_type              = local.central_network_connections[0].bandwidth_type
  type                        = local.central_network_connections[0].type
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_id" {
  central_network_id = huaweicloud_cc_central_network.test.id
  connection_id      = local.id
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_status" {
  central_network_id = huaweicloud_cc_central_network.test.id
  status             = local.status
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_bwtype" {
  central_network_id = huaweicloud_cc_central_network.test.id
  bandwidth_type     = local.bandwidth_type
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_type" {
  central_network_id = huaweicloud_cc_central_network.test.id
  type               = local.type
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_cross" {
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  central_network_id = huaweicloud_cc_central_network.test.id
  is_cross_region    = "true"
}
  
data "huaweicloud_cc_central_network_connections" "filter_by_not_cross" {
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  central_network_id = huaweicloud_cc_central_network.test.id
  is_cross_region    = "false"
}
  
locals {
  connsById       = data.huaweicloud_cc_central_network_connections.filter_by_id.central_network_connections
  connsBystatus   = data.huaweicloud_cc_central_network_connections.filter_by_status.central_network_connections
  connsByBwType   = data.huaweicloud_cc_central_network_connections.filter_by_bwtype.central_network_connections
  connsByType     = data.huaweicloud_cc_central_network_connections.filter_by_type.central_network_connections
  connsByCross    = data.huaweicloud_cc_central_network_connections.filter_by_cross.central_network_connections
  connsByNotCross = data.huaweicloud_cc_central_network_connections.filter_by_not_cross.central_network_connections
}
  
output "id_filter_is_useful" {
  value = length(local.connsById) > 0 && alltrue([for v in local.connsById[*].id : v == local.id])
}
  
output "status_filter_is_useful" {
  value = length(local.connsBystatus) > 0 && alltrue([for v in local.connsBystatus[*].status : v == local.status])
}
  
output "bandwidth_type_filter_is_useful" {
  value = length(local.connsByBwType) > 0 && alltrue(
    [for v in local.connsByBwType[*].bandwidth_type : v == local.bandwidth_type]
  )
}
  
output "type_filter_is_useful" {
  value = length(local.connsByType) > 0 && alltrue([for v in local.connsByType[*].type : v == local.type])
}
  
output "cross_filter_is_useful" {
  value = length(local.connsByCross) == 3 && length(local.connsByNotCross) == 0
}
`, testCentralNetworkConnection_dataBasic(name))
}

func testCentralNetworkConnection_dataBasic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "az1" {
  region = "%[1]s"
}

resource "huaweicloud_er_instance" "er1" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az1.names, 0, 1)

  region                         = "%[1]s"
  name                           = "%[4]s1"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

data "huaweicloud_er_availability_zones" "az2" {
  region = "%[2]s"
}
  
resource "huaweicloud_er_instance" "er2" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az2.names, 0, 1)

  region                         = "%[2]s"
  name                           = "%[4]s2"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

data "huaweicloud_er_availability_zones" "az3" {
  region = "%[3]s"
}
  
resource "huaweicloud_er_instance" "er3" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az3.names, 0, 1)

  region                         = "%[3]s"
  name                           = "%[4]s3"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

resource "huaweicloud_cc_central_network" "test" {
  name        = "%[4]s"
  description = "This is an accaptance test"
}
 
resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
 
  planes {
    associate_er_tables {
      project_id                 = "%[5]s"
      region_id                  = "%[1]s"
      enterprise_router_id       = huaweicloud_er_instance.er1.id
      enterprise_router_table_id = huaweicloud_er_instance.er1.default_association_route_table_id
    }

    associate_er_tables {
      project_id                 = "%[6]s"
      region_id                  = "%[2]s"
      enterprise_router_id       = huaweicloud_er_instance.er2.id
      enterprise_router_table_id = huaweicloud_er_instance.er2.default_association_route_table_id
    }

    associate_er_tables {
      project_id                 = "%[7]s"
      region_id                  = "%[3]s"
      enterprise_router_id       = huaweicloud_er_instance.er3.id
      enterprise_router_table_id = huaweicloud_er_instance.er3.default_association_route_table_id
    }
  }
 
  er_instances {
    project_id           = "%[5]s"
    region_id            = "%[1]s"
    enterprise_router_id = huaweicloud_er_instance.er1.id
  }

  er_instances {
    project_id           = "%[6]s"
    region_id            = "%[2]s"
    enterprise_router_id = huaweicloud_er_instance.er2.id
  }

  er_instances {
    project_id           = "%[7]s"
    region_id            = "%[3]s"
    enterprise_router_id = huaweicloud_er_instance.er3.id
  }
}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test.id
}
`, acceptance.HW_REGION_NAME_1, acceptance.HW_REGION_NAME_2, acceptance.HW_REGION_NAME_3,
		name, acceptance.HW_PROJECT_ID_1, acceptance.HW_PROJECT_ID_2, acceptance.HW_PROJECT_ID_3)
}
