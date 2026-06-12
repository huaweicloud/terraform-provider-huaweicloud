package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussDbPluginLicenseConfig_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()
	resourceName := "huaweicloud_gaussdb_instance_plugin_license_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDbPluginLicenseConfig_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_gaussdb_instance.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "license_str"),
				),
			},
		},
	})
}

func testAccGaussDbPluginLicenseConfig_base(name string, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0],
                      data.huaweicloud_availability_zones.test.names[1],
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[4]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name, password, acceptance.HW_ENTERPRISE_PROJECT_ID)
}

func testAccGaussDbPluginLicenseConfig_basic(name string, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance_plugin_license_config" "test" {
 instance_id      = huaweicloud_gaussdb_instance.test.id
 license_str      = "90ca926757144a2d98d78f727dc56664"
 enable_force_new = "true"
}
`, testAccGaussDbPluginLicenseConfig_base(name, password))
}
