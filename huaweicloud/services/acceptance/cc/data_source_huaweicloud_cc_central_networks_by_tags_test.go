package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworksByTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_cc_central_networks_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCConnectionRouteRegionName(t)
			acceptance.TestAccPreCheckCCConnectionRouteProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCentralNetworksByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.default_plane_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.auto_associate_route_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.auto_propagate_route_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.0.associate_er_tables.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.associate_er_tables.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.associate_er_tables.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.associate_er_tables.0.enterprise_router_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.associate_er_tables.0.enterprise_router_table_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.0.exclude_er_connections.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.exclude_er_connections.0.exclude_er_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.exclude_er_connections.0.exclude_er_instances.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.exclude_er_connections.0.exclude_er_instances.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource,
						"central_networks.0.planes.0.exclude_er_connections.0.exclude_er_instances.0.enterprise_router_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.planes.0.is_full_mesh"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.0.enterprise_router_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.0.asn"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.er_instances.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "central_networks.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceCcCentralNetworksByTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "az1" {
 region = "%[2]s"
}

resource "huaweicloud_er_instance" "er1" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az1.names, 0, 1)

  region                         = "%[2]s"
  name                           = "%[1]s_1"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

data "huaweicloud_er_availability_zones" "az2" {
  region = "%[3]s"
}

resource "huaweicloud_er_instance" "er2" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az2.names, 0, 1)

  region                         = "%[3]s"
  name                           = "%[1]s_2"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

resource "huaweicloud_cc_central_network" "test" {
  name        = "%[1]s"
  description = "This is an accaptance test"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id

  planes {
    associate_er_tables {
      project_id                 = "%[4]s"
      region_id                  = "%[2]s"
      enterprise_router_id       = huaweicloud_er_instance.er1.id
      enterprise_router_table_id = huaweicloud_er_instance.er1.default_association_route_table_id
    }
    associate_er_tables {
      project_id                 = "%[5]s"
      region_id                  = "%[3]s"
      enterprise_router_id       = huaweicloud_er_instance.er2.id
      enterprise_router_table_id = huaweicloud_er_instance.er2.default_association_route_table_id
    }

   exclude_er_connections {
     exclude_er_instances {
       project_id           = "%[4]s"
       region_id            = "%[2]s"
       enterprise_router_id = huaweicloud_er_instance.er1.id
     }
     exclude_er_instances {
       project_id           = "%[5]s"
       region_id            = "%[3]s"
       enterprise_router_id = huaweicloud_er_instance.er2.id
     }
  }
}

  er_instances {
    project_id           = "%[4]s"
    region_id            = "%[2]s"
    enterprise_router_id = huaweicloud_er_instance.er1.id
  }
  er_instances {
    project_id           = "%[5]s"
    region_id            = "%[3]s"
    enterprise_router_id = huaweicloud_er_instance.er2.id
  }
}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test.id
}
`, name, acceptance.HW_REGION_NAME_1, acceptance.HW_REGION_NAME_2, acceptance.HW_PROJECT_ID_1, acceptance.HW_PROJECT_ID_2)
}

func testDataSourceCcCentralNetworksByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_central_networks_by_tags" "test"{
  depends_on = [huaweicloud_cc_central_network_policy_apply.test]

  tags {
    key    = "key"
    values = ["value"]
  }
}
`, testDataSourceCcCentralNetworksByTags_base(name))
}
