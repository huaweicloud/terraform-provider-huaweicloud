package sfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTurbosDataSource_basic(t *testing.T) {
	var (
		rName         = acceptance.RandomAccResourceNameWithDash()
		dcByName      = acceptance.InitDataSourceCheck("data.huaweicloud_sfs_turbos.by_name")
		dcBySize      = acceptance.InitDataSourceCheck("data.huaweicloud_sfs_turbos.by_size")
		dcByShareType = acceptance.InitDataSourceCheck("data.huaweicloud_sfs_turbos.by_share_type")
		dcByEpsId     = acceptance.InitDataSourceCheck("data.huaweicloud_sfs_turbos.by_eps_id")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTurbosDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_query_result_validation", "true"),
					dcBySize.CheckResourceExists(),
					resource.TestCheckOutput("size_query_result_validation", "true"),
					dcByShareType.CheckResourceExists(),
					resource.TestCheckOutput("share_type_query_result_validation", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("eps_id_query_result_validation", "true"),
				),
			},
		},
	})
}

func testAccTurbosDataSource_basic(rName string) string {
	return fmt.Sprintf(`
variable "turbo_configuration" {
  type = list(object({
    size        = number
    share_type  = string
    eps_enabled = bool
  }))

  default = [
    {size = 100, share_type = "PERFORMANCE", eps_enabled = false},
    {size = 200, share_type = "STANDARD", eps_enabled = false},
    {size = 200, share_type = "PERFORMANCE", eps_enabled = true},
  ]
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_sfs_turbo" "test" {
  count = length(var.turbo_configuration)

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  name                  = "%[1]s-${count.index}"
  size                  = var.turbo_configuration[count.index]["size"]
  share_proto           = "NFS"
  share_type            = var.turbo_configuration[count.index]["share_type"]
  enterprise_project_id = var.turbo_configuration[count.index]["eps_enabled"] ? "%[2]s" : "0"
}

data "huaweicloud_sfs_turbos" "by_name" {
  depends_on = [huaweicloud_sfs_turbo.test]

  name = huaweicloud_sfs_turbo.test[0].name
}

data "huaweicloud_sfs_turbos" "by_size" {
  depends_on = [huaweicloud_sfs_turbo.test]

  size = var.turbo_configuration[0]["size"]
}

data "huaweicloud_sfs_turbos" "by_share_type" {
  depends_on = [huaweicloud_sfs_turbo.test]

  share_type = var.turbo_configuration[1]["share_type"]
}

data "huaweicloud_sfs_turbos" "by_eps_id" {
  depends_on = [huaweicloud_sfs_turbo.test]

  enterprise_project_id = "%[2]s"
}

output "name_query_result_validation" {
  value = contains(data.huaweicloud_sfs_turbos.by_name.turbos[*].id,
  huaweicloud_sfs_turbo.test[0].id) && !contains(data.huaweicloud_sfs_turbos.by_name.turbos[*].id,
  huaweicloud_sfs_turbo.test[1].id) && !contains(data.huaweicloud_sfs_turbos.by_name.turbos[*].id,
  huaweicloud_sfs_turbo.test[2].id)
}

output "size_query_result_validation" {
  value = contains(data.huaweicloud_sfs_turbos.by_size.turbos[*].id,
  huaweicloud_sfs_turbo.test[0].id) && !contains(data.huaweicloud_sfs_turbos.by_size.turbos[*].id,
  huaweicloud_sfs_turbo.test[1].id) && !contains(data.huaweicloud_sfs_turbos.by_size.turbos[*].id,
  huaweicloud_sfs_turbo.test[2].id)
}

output "share_type_query_result_validation" {
  value = contains(data.huaweicloud_sfs_turbos.by_share_type.turbos[*].id,
  huaweicloud_sfs_turbo.test[1].id) && !contains(data.huaweicloud_sfs_turbos.by_share_type.turbos[*].id,
  huaweicloud_sfs_turbo.test[0].id) && !contains(data.huaweicloud_sfs_turbos.by_share_type.turbos[*].id,
  huaweicloud_sfs_turbo.test[2].id)
}

output "eps_id_query_result_validation" {
  value = contains(data.huaweicloud_sfs_turbos.by_eps_id.turbos[*].id,
  huaweicloud_sfs_turbo.test[2].id) && !contains(data.huaweicloud_sfs_turbos.by_eps_id.turbos[*].id,
  huaweicloud_sfs_turbo.test[0].id) && !contains(data.huaweicloud_sfs_turbos.by_eps_id.turbos[*].id,
  huaweicloud_sfs_turbo.test[1].id)
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
