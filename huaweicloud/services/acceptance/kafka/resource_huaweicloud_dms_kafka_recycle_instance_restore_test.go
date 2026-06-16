package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRecycleInstanceRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaRecycleBinInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceRecycleInstanceRestore_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccResourceRecycleInstanceRestore_basic(),
			},
		},
	})
}

func testAccResourceRecycleInstanceRestore_instanceNotFound() string {
	randomUUID, _ := uuid.NewRandom()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_recycle_instance_restore" "test" {
  instance_id = "%s"
}
`, randomUUID.String())
}

func testAccResourceRecycleInstanceRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_recycle_instance_restore" "test" {
  instance_id = "%s"
}
`, acceptance.HW_DMS_KAFKA_RECYCLE_BIN_INSTANCE_ID)
}
