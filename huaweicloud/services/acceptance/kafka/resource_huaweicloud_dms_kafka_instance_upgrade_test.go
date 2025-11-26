package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This is a one-time action resource, so it does not need to be destroyed.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccInstanceUpgrade_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccInstanceUpgrade_basic(),
			},
		},
	})
}

func testAccInstanceUpgrade_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_instance_upgrade" "test" {
  instance_id = "%[1]s"
}`, randomId)
}

func testAccInstanceUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_instance_upgrade" "test" {
  instance_id = "%[1]s"
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
