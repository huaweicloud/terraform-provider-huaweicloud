package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_policies.test"
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
				Config: testDataSourceDataSourceCcCentralNetworkPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policies.0.central_network_id"),
					resource.TestCheckResourceAttrSet(dataSource, "central_network_policies.0.created_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("is_applied_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcCentralNetworkPolicies_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_central_network_policies" "test" {
  depends_on = [
    huaweicloud_cc_central_network_policy.test2,
    huaweicloud_cc_central_network_policy_apply.test,
  ]
 
  central_network_id = huaweicloud_cc_central_network.test.id
}
	  
locals {
  central_network_policies = data.huaweicloud_cc_central_network_policies.test.central_network_policies
  id                       = local.central_network_policies[0].id
  status                   = local.central_network_policies[0].status
}

data "huaweicloud_cc_central_network_policies" "filter_by_id" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = local.id
}

data "huaweicloud_cc_central_network_policies" "filter_by_status" { 
  central_network_id = huaweicloud_cc_central_network.test.id
  status             = local.status
}

data "huaweicloud_cc_central_network_policies" "filter_by_applied" {
  depends_on = [
    huaweicloud_cc_central_network_policy.test2,
    huaweicloud_cc_central_network_policy_apply.test,
  ]

  central_network_id = huaweicloud_cc_central_network.test.id
  is_applied         = "true"
}

data "huaweicloud_cc_central_network_policies" "filter_by_not_applied" {
  depends_on = [
    huaweicloud_cc_central_network_policy.test2,
    huaweicloud_cc_central_network_policy_apply.test,
  ]

  central_network_id = huaweicloud_cc_central_network.test.id
  is_applied         = "false"
}

locals {
  policiesById         = data.huaweicloud_cc_central_network_policies.filter_by_id.central_network_policies
  policiesBystatus     = data.huaweicloud_cc_central_network_policies.filter_by_status.central_network_policies
  policiesByApplied    = data.huaweicloud_cc_central_network_policies.filter_by_applied.central_network_policies
  policiesByNotApplied = data.huaweicloud_cc_central_network_policies.filter_by_not_applied.central_network_policies
}

output "id_filter_is_useful" {
value = length(local.policiesById) > 0 && alltrue([for v in local.policiesById[*].id : v == local.id])
}

output "status_filter_is_useful" {
value = length(local.policiesBystatus) > 0 && alltrue([for v in local.policiesBystatus[*].status : v == local.status])
}

output "is_applied_filter_is_useful" {
value = length(local.policiesByApplied) == 1 && length(local.policiesByNotApplied) == 2
}
`, testCentralNetworkPolicies_dataBasic(name))
}

func testCentralNetworkPolicies_dataBasic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "az1" {
  region = "%[1]s"
}

resource "huaweicloud_er_instance" "er1" {
  availability_zones             = slice(data.huaweicloud_er_availability_zones.az1.names, 0, 1)
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
  availability_zones             = slice(data.huaweicloud_er_availability_zones.az2.names, 0, 1)
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
  availability_zones             = slice(data.huaweicloud_er_availability_zones.az3.names, 0, 1)
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
 
resource "huaweicloud_cc_central_network_policy" "test1" {
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
}

resource "huaweicloud_cc_central_network_policy" "test2" {
  central_network_id = huaweicloud_cc_central_network.test.id
   
  planes {
    associate_er_tables {
      project_id                 = "%[7]s"
      region_id                  = "%[3]s"
      enterprise_router_id       = huaweicloud_er_instance.er3.id
      enterprise_router_table_id = huaweicloud_er_instance.er3.default_association_route_table_id
    }
  }
   
  er_instances {
    project_id           = "%[7]s"
    region_id            = "%[3]s"
    enterprise_router_id = huaweicloud_er_instance.er3.id
  }
}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = huaweicloud_cc_central_network.test.id
  policy_id          = huaweicloud_cc_central_network_policy.test1.id
}
`, acceptance.HW_REGION_NAME_1, acceptance.HW_REGION_NAME_2, acceptance.HW_REGION_NAME_3,
		name, acceptance.HW_PROJECT_ID_1, acceptance.HW_PROJECT_ID_2, acceptance.HW_PROJECT_ID_3)
}
