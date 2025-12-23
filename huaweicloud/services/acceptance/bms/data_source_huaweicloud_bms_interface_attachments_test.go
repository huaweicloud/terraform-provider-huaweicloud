package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBmsInterfaceAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_bms_interface_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBmsInterfaceAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.port_state"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.fixed_ips.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.net_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "interface_attachments.0.mac_addr"),
				),
			},
		},
	})
}

func testDataSourceBmsInterfaceAttachments_base(rName string) string {
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

func testDataSourceBmsInterfaceAttachments_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_bms_interface_attachments" "test" {
  server_id = huaweicloud_bms_instance.test.id
}
`, testDataSourceBmsInterfaceAttachments_base(name))
}
