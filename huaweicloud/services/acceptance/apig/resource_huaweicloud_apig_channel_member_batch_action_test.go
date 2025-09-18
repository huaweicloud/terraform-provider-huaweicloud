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

func testApigChannelMemberBatchAction_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

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
