package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccApigChannelMemberBatchAction_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()
		rc   = "huaweicloud_apig_channel_member_batch_action.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testApigChannelMemberBatchAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rc, "action", "enable"),
					resource.TestCheckResourceAttrSet(rc, "instance_id"),
					resource.TestCheckResourceAttrSet(rc, "vpc_channel_id"),
					resource.TestCheckResourceAttrSet(rc, "id"),
				),
			},
		},
	})
}

func testApigChannelMemberBatchAction_compute_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  id = "%[2]s"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                = format("%[1]s_%%d", count.index)
  description         = "terraform test"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test111"
  agent_list          = "hss"

  user_data = <<EOF
#! /bin/bash
echo user_test > /home/user.txt
EOF

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }
}`, name, acceptance.HW_SUBNET_ID)
}

func testApigChannelMemberBatchAction_base(name string) string {
	return fmt.Sprintf(`

	
resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  edition               = "BASIC"
  name                  = "%[2]s"
  description           = "created by acc test for channel member batch action"
}

resource "huaweicloud_apig_channel" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
  type        = "TCP"
  port        = 8080
  balance_strategy = "ROUND_ROBIN"
  member_type = "ECS"
  protocol   = "TCP"
  path       = "/"
  algorithm  = "WRR"
  vpc_health_config {
    protocol            = "TCP"
    delay               = 10
    timeout             = 3
    max_retries         = 3
    port                = 8080
    interval            = 5
  }
}
`, common.TestBaseNetwork(name), name)
}

func testApigChannelMemberBatchAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_channel_member_batch_action" "test" {
  instance_id     = huaweicloud_apig_instance.test.id
  vpc_channel_id  = huaweicloud_apig_channel.test.id
  action          = "enable"
  member_ids      = []
}
`, testApigChannelMemberBatchAction_base(name))
}
