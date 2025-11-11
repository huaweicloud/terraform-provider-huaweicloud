package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataVolumeAutoExpandConfiguration_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_volume_auto_expand_configuration.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataVolumeAutoExpandConfiguration_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccDataVolumeAutoExpandConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "auto_volume_expand_enable", "true"),
					resource.TestCheckResourceAttrPair(dataSource, "expand_threshold",
						"huaweicloud_dms_kafka_volume_auto_expand_configuration.test", "expand_threshold"),
					resource.TestCheckResourceAttrPair(dataSource, "max_volume_size",
						"huaweicloud_dms_kafka_volume_auto_expand_configuration.test", "max_volume_size"),
					resource.TestCheckResourceAttrPair(dataSource, "expand_increment",
						"huaweicloud_dms_kafka_volume_auto_expand_configuration.test", "expand_increment"),
				),
			},
		},
	})
}

func testAccDataVolumeAutoExpandConfiguration_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id = "%[1]s"
}
`, randomId)
}

func testAccDataVolumeAutoExpandConfiguration_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instances" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id               = "%[1]s"
  auto_volume_expand_enable = true
  expand_threshold          = 80
  expand_increment          = 10
  # Must be greater than the current instance disk capacity.
  max_volume_size = try(data.huaweicloud_dms_kafka_instances.test.instances[0].storage_space, 0) + 100
}

data "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id = "%[1]s"

  depends_on = [huaweicloud_dms_kafka_volume_auto_expand_configuration.test]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
