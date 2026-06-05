package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile("^[1-9]([0-9]*)?$")),
					resource.TestCheckResourceAttrSet(all, "instances.0.id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.status"),
					resource.TestCheckResourceAttrSet(all, "instances.0.version"),
					resource.TestCheckResourceAttrSet(all, "instances.0.ip"),
					resource.TestCheckResourceAttrSet(all, "instances.0.port"),
					resource.TestCheckResourceAttrSet(all, "instances.0.cpu"),
					resource.TestCheckResourceAttrSet(all, "instances.0.mem"),
					resource.TestCheckResourceAttrSet(all, "instances.0.login_flag"),
					resource.TestCheckResourceAttrSet(all, "instances.0.slow_sql_flag"),
					resource.TestCheckResourceAttrSet(all, "instances.0.deadlock_flag"),
					resource.TestCheckResourceAttrSet(all, "instances.0.lock_blocking_flag"),
					resource.TestCheckResourceAttrSet(all, "instances.0.charge_flag"),
					resource.TestCheckResourceAttrSet(all, "instances.0.full_sql_flag"),
				),
			},
		},
	})
}

func testAccInstances_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "MySQL"
    version = "8.0"
    port    = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_instances" "test" {
  datastore_type = "MySQL"

  depends_on = [
    huaweicloud_rds_instance.test,
  ]
}
`, testAccInstances_base(name))
}
