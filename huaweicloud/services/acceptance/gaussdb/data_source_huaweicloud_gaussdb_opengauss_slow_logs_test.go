package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceOpenGaussSlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_slow_logs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOpenGaussSlowLogs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.workflow_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_size"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.file_link"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.status"),
				),
			},
		},
	})
}

func testDataSourceOpenGaussSlowLogs_base(rName string) string {
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

  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  flavor                = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name                  = "%[2]s"
  password              = "Huangwei!120521"
  enterprise_project_id = "%[3]s"

  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

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

func testDataSourceOpenGaussSlowLogs_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_slow_logs" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
`, testDataSourceOpenGaussSlowLogs_base(rName))
}
