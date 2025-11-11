package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVolumeAutoExpandConfiguration_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeAutoExpandConfiguration_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccVolumeAutoExpandConfiguration_basic_step1(),
			},
			{
				Config: testAccVolumeAutoExpandConfiguration_basic_step2(),
			},
		},
	})
}

func testAccVolumeAutoExpandConfiguration_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id               = "%[1]s"
  auto_volume_expand_enable = true
  expand_threshold          = 80
  expand_increment          = 10
  max_volume_size           = 400
}`, randomId)
}

func testAccVolumeAutoExpandConfiguration_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instances" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "enable" {
  instance_id               = "%[1]s"
  auto_volume_expand_enable = true
  expand_threshold          = 80
  expand_increment          = 10
  # Must be greater than the current instance disk capacity.
  max_volume_size = try(data.huaweicloud_dms_kafka_instances.test.instances[0].storage_space, 0) + 100
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccVolumeAutoExpandConfiguration_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "disable" {
  instance_id               = "%[1]s"
  auto_volume_expand_enable = false
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
