package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceRestart_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceRestart_basic(rName),
			},
			{
				Config: testAccInstanceRestart_delay(rName),
			},
		},
	})
}

func testAccInstanceRestart_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2022_SE"
  instance_mode = "single"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  tde_enabled       = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), rName)
}

func testAccInstanceRestart_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance_restart" "test" {
  instance_id    = huaweicloud_rds_instance.test.id
  restart_server = true
  forcible       = true
}
`, testAccInstanceRestart_base(rName))
}

func testAccInstanceRestart_delay(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance_restart" "test" {
  instance_id    = huaweicloud_rds_instance.test.id
  restart_server = true
  forcible       = true
}

resource "huaweicloud_rds_instance_restart" "test_delay" {
  depends_on = [huaweicloud_rds_instance_restart.test]

  instance_id    = huaweicloud_rds_instance.test.id
  restart_server = false
  forcible       = false
  delay          = true
}
`, testAccInstanceRestart_base(rName))
}
