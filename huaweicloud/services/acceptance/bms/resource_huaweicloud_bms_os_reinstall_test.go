package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBmsOsReinstall_Basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsOsReinstall_basic(rName),
			},
		},
	})
}

func testAccBmsOsReinstall_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")
  name              = "%[2]s"
  user_id           = "%[3]s"
  system_disk_type  = "GPSSD"
  system_disk_size  = 150
  power_action      = "OFF"

  user_data = <<EOF
#!/bin/bash 
sudo mkdir /example
EOF

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  metadata = {
    foo1 = "bar1"
    key1 = "value1"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"
}
`, testAccBmsInstance_base(rName), rName, acceptance.HW_USER_ID)
}

func testAccBmsOsReinstall_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_bms_os_reinstall" "test" {
  server_id = huaweicloud_bms_instance.test.id

  os_reinstall {
    key_name = huaweicloud_kps_keypair.test.name
    user_id  = "%[2]s"

    metadata {
      user_data = <<EOF
#!/bin/bash 
sudo mkdir /example
EOF
      __system__encrypted = "0"
    }
  }
}
`, testAccBmsOsReinstall_base(rName), acceptance.HW_USER_ID)
}
