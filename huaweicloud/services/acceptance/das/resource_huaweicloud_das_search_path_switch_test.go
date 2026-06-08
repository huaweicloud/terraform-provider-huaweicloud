package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccSearchPathSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_search_path_switch.test"

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSearchPathSwitch_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switch_on", "true"),
				),
			},
			{
				Config: testAccSearchPathSwitch_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switch_on", "false"),
				),
			},
		},
	})
}

func testAccSearchPathSwitch_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "random_password" "test" {
  length           = 12
  min_numeric      = 1
  min_upper        = 1
  min_lower        = 1
  special          = true
  min_special      = 1
  override_special = "!@"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "postgresql"
  db_version    = "16"
  instance_mode = "single"
  vcpus         = 1
  memory        = 2
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type     = "postgresql"
    version  = "16"
    port     = 8634
    password = random_password.test.result
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

resource "huaweicloud_das_database_instance_connection" "test" {
  instance_id      = huaweicloud_rds_instance.test.id
  engine_type      = "postgresql"
  network_type     = "rds"
  username         = "root"
  password         = random_password.test.result
  is_save_password = false
  description      = "Created by terraform script!"
}

`, common.TestBaseNetwork(name), name)
}

func testAccSearchPathSwitch_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_search_path_switch" "test" {
  connection_id = huaweicloud_das_database_instance_connection.test.id
  switch_on     = true
}
`, testAccSearchPathSwitch_base(name))
}

func testAccSearchPathSwitch_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_search_path_switch" "test" {
  connection_id = huaweicloud_das_database_instance_connection.test.id
  switch_on     = false
}
`, testAccSearchPathSwitch_base(name))
}
