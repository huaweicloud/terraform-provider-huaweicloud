package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccConfigurationApply_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationApply_apply(rName),
			},
		},
	})
}

func testAccConfigurationApply_apply_base(name string) string {
	return fmt.Sprintf(`
%[1]s

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
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccConfigurationApply_apply(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_parametergroup" "test" {
  name        = "%[2]s"
  description = "description_1"

  values = {
    auto_increment_increment     = "2"
    binlog_rows_query_log_events = "ON"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}

resource "huaweicloud_rds_parametergroup_apply" "test" {
  config_id   = huaweicloud_rds_parametergroup.test.id
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccConfigurationApply_apply_base(rName), rName)
}
