package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rocketmq"
)

func getVolumeAutoExpandConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return rocketmq.GetVolumeAutoExpandConfiguration(client, state.Primary.ID)
}

func TestAccVolumeAutoExpandConfiguration_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_dms_rocketmq_volume_auto_expand_configuration.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getVolumeAutoExpandConfigurationResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeAutoExpandConfiguration_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_volume_expand_enable", "true"),
					resource.TestCheckResourceAttr(rName, "expand_threshold", "80"),
					resource.TestCheckResourceAttr(rName, "expand_increment", "10"),
					resource.TestCheckResourceAttrSet(rName, "max_volume_size"),
				),
			},
			{
				Config: testAccVolumeAutoExpandConfiguration_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "auto_volume_expand_enable", "false"),
				),
			},
		},
	})
}

func testAccVolumeAutoExpandConfiguration_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_instances" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_rocketmq_volume_auto_expand_configuration" "test" {
  instance_id               = "%[1]s"
  auto_volume_expand_enable = true
  expand_threshold          = 80
  expand_increment          = 10
  # Must be greater than the current instance disk capacity.
  max_volume_size           = try(data.huaweicloud_dms_rocketmq_instances.test.instances[0].storage_space, 0) + 100
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}

func testAccVolumeAutoExpandConfiguration_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rocketmq_volume_auto_expand_configuration" "test" {
  instance_id      = "%[1]s"
  enable_force_new = "true"
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}
