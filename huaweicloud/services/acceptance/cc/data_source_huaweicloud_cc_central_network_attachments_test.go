package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_attachments.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCProjectID(t)
			acceptance.TestAccPreCheckCCRegionName(t)
			acceptance.TestAccPreCheckCCGlobalGateway(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCentralNetworkAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.central_network_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.enterprise_router_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.enterprise_router_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_attachments.0.enterprise_router_region_id"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcCentralNetworkAttachments_basic(name string) string {
	return fmt.Sprintf(`
%s

locals {
  id   = huaweicloud_cc_central_network_attachment.test.id
  name = huaweicloud_cc_central_network_attachment.test.name
  type = "GDGW"
}

data "huaweicloud_cc_central_network_attachments" "filter_by_id" {
  central_network_id = huaweicloud_cc_central_network.test.id
  attachment_id      = local.id

  depends_on = [
    huaweicloud_cc_central_network_attachment.test,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cc_central_network_attachments.filter_by_id.central_network_attachments) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_central_network_attachments.filter_by_id.central_network_attachments[*].id : v == local.id]
  )
}

data "huaweicloud_cc_central_network_attachments" "filter_by_name" {
  central_network_id = huaweicloud_cc_central_network.test.id
  name               = local.name

  depends_on = [
    huaweicloud_cc_central_network_attachment.test,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cc_central_network_attachments.filter_by_name.central_network_attachments) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_central_network_attachments.filter_by_name.central_network_attachments[*].name : v == local.name]
  )
}

data "huaweicloud_cc_central_network_attachments" "filter_by_type" {
  central_network_id       = huaweicloud_cc_central_network.test.id
  attachment_instance_type = local.type

  depends_on = [
    huaweicloud_cc_central_network_attachment.test,
  ]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_cc_central_network_attachments.filter_by_type.central_network_attachments) >= 1 && alltrue([
    for v in data.huaweicloud_cc_central_network_attachments.filter_by_type.central_network_attachments[*].attachment_instance_type : 
    v == local.type
  ])
}
`, testDataSourceCcCentralNetworkAttachments_base(name))
}

func testDataSourceCcCentralNetworkAttachments_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "az1" {
  region = "%[1]s"
}

resource "huaweicloud_er_instance" "er1" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.az1.names, 0, 1)

  region                         = "%[1]s"
  name                           = "%[3]s1"
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
  name                           = "%[3]s2"
  asn                            = 64512
  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true
}

resource "huaweicloud_cc_central_network" "test" {
  name        = "%[3]s"
  description = "This is an accaptance test"
}
 
resource "huaweicloud_cc_central_network_policy" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
 
  planes {
    associate_er_tables {
      project_id                 = "%[4]s"
      region_id                  = "%[1]s"
      enterprise_router_id       = huaweicloud_er_instance.er1.id
      enterprise_router_table_id = huaweicloud_er_instance.er1.default_association_route_table_id
    }

    associate_er_tables {
      project_id                 = "%[5]s"
      region_id                  = "%[2]s"
      enterprise_router_id       = huaweicloud_er_instance.er2.id
      enterprise_router_table_id = huaweicloud_er_instance.er2.default_association_route_table_id
    }
  }
 
  er_instances {
    project_id           = "%[4]s"
    region_id            = "%[1]s"
    enterprise_router_id = huaweicloud_er_instance.er1.id
  }

  er_instances {
    project_id           = "%[5]s"
    region_id            = "%[2]s"
    enterprise_router_id = huaweicloud_er_instance.er2.id
  }
}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test.id
}

resource "huaweicloud_cc_central_network_attachment" "test" {
  name                         = "%[3]s"
  description                  = "This is a demo"
  central_network_id           = huaweicloud_cc_central_network.test.id
  enterprise_router_id         = huaweicloud_er_instance.er1.id
  enterprise_router_project_id = "%[4]s"
  enterprise_router_region_id  = "%[1]s"
  global_dc_gateway_id         = "%[6]s"
  global_dc_gateway_project_id = "%[4]s"
  global_dc_gateway_region_id  = "%[1]s"

  depends_on = [
    huaweicloud_cc_central_network_policy_apply.test,
  ]
}
`, acceptance.HW_REGION_NAME_1, acceptance.HW_REGION_NAME_2, name, acceptance.HW_PROJECT_ID_1,
		acceptance.HW_PROJECT_ID_2, acceptance.HW_CC_GLOBAL_GATEWAY_ID)
}
