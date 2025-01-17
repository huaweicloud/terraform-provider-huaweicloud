package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussUpgradeVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_upgrade_versions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussUpgradeVersions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_type_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_type_list.0.enable"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_type_list.0.is_parallel_upgrade"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_type_list.0.upgrade_action_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_type_list.0.upgrade_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rollback_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "source_version"),
					resource.TestCheckResourceAttrSet(dataSource, "roll_upgrade_progress.#"),
					resource.TestCheckResourceAttrSet(dataSource, "roll_upgrade_progress.0.az_description_map.%"),
					resource.TestCheckResourceAttrSet(dataSource, "roll_upgrade_progress.0.not_fully_upgraded_az"),
					resource.TestCheckResourceAttrSet(dataSource, "upgrade_candidate_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hotfix_upgrade_candidate_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hotfix_rollback_candidate_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hotfix_upgrade_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hotfix_rollback_infos.#"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussUpgradeVersions_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
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

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
  password          = "Huangwei!120521"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

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
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussUpgradeVersions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_upgrade_versions" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
`, testDataSourceGaussdbOpengaussUpgradeVersions_base(name))
}
