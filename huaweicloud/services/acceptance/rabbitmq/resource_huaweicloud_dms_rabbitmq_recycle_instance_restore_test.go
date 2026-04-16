package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRecycleInstanceRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQRecycleBinInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccRecycleInstanceRestore_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccRecycleInstanceRestore_basic(),
			},
		},
	})
}

func testAccRecycleInstanceRestore_instanceNotFound() string {
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_rabbitmq_recycle_instance_restore" "test" {
  instance_id = "%s"
}
`, randomUUID)
}

func testAccRecycleInstanceRestore_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rabbitmq_recycle_instance_restore" "test" {
  instance_id = "%s"
}
`, acceptance.HW_DMS_RABBITMQ_RECYCLE_BIN_INSTANCE_ID)
}
