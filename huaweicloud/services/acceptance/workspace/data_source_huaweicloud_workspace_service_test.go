package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataService_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_service.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAD(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataService_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_service_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_ad_domain_set_and_valid", "true"),
					resource.TestCheckOutput("is_auth_type_set_and_valid", "true"),
					resource.TestCheckOutput("is_access_mode_set_and_valid", "true"),
					resource.TestCheckOutput("is_vpc_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_network_ids_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataService_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_service" "test" {
  ad_domain {
    name               = try(element(regexall("\\w+\\.(.*)", element(split(",", "%[1]s"), 0))[0], 0), "")
    active_domain_name = element(split(",", "%[1]s"), 0)
    active_domain_ip   = element(split(",", "%[2]s"), 0)
    active_dns_ip      = element(split(",", "%[2]s"), 0)
    // Standby configuration
    standby_domain_name = element(split(",", "%[1]s"), 1)
    standby_domain_ip   = element(split(",", "%[2]s"), 1)
    standby_dns_ip      = element(split(",", "%[2]s"), 1)
    admin_account      = "Administrator"
    password           = "%[3]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[4]s"
  network_ids = ["%[5]s"]
}

data "huaweicloud_workspace_service" "test" {
  depends_on = [
    huaweicloud_workspace_service.test
  ]
}

locals {
  service_result = data.huaweicloud_workspace_service.test
}

output "is_service_id_set_and_valid" {
  value = local.service_result.id == huaweicloud_workspace_service.test.id
}

output "is_ad_domain_set_and_valid" {
  value = alltrue([
    local.service_result.ad_domain[0].name == huaweicloud_workspace_service.test.ad_domain[0].name,
    local.service_result.ad_domain[0].active_domain_name == huaweicloud_workspace_service.test.ad_domain[0].active_domain_name,
    local.service_result.ad_domain[0].active_domain_ip == huaweicloud_workspace_service.test.ad_domain[0].active_domain_ip,
    local.service_result.ad_domain[0].active_dns_ip == huaweicloud_workspace_service.test.ad_domain[0].active_dns_ip,
    local.service_result.ad_domain[0].standby_domain_name == huaweicloud_workspace_service.test.ad_domain[0].standby_domain_name,
    local.service_result.ad_domain[0].standby_domain_ip == huaweicloud_workspace_service.test.ad_domain[0].standby_domain_ip,
    local.service_result.ad_domain[0].standby_dns_ip == huaweicloud_workspace_service.test.ad_domain[0].standby_dns_ip,
    local.service_result.ad_domain[0].admin_account == huaweicloud_workspace_service.test.ad_domain[0].admin_account,
  ])
}

output "is_auth_type_set_and_valid" {
  value = local.service_result.auth_type == huaweicloud_workspace_service.test.auth_type
}

output "is_access_mode_set_and_valid" {
  value = local.service_result.access_mode == huaweicloud_workspace_service.test.access_mode
}

output "is_vpc_id_set_and_valid" {
  value = local.service_result.vpc_id == huaweicloud_workspace_service.test.vpc_id
}

output "is_network_ids_set_and_valid" {
  value = alltrue([
    length(local.service_result.network_ids) == 1,
    local.service_result.network_ids[0] == huaweicloud_workspace_service.test.network_ids[0],
  ])
}
`, acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_IPS,
		acceptance.HW_WORKSPACE_AD_SERVER_PWD,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID)
}
