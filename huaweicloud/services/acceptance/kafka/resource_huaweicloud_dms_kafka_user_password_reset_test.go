package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, ensure that ciphertext access is enabled on the Kafka instance.
func TestAccUserPasswordReset_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// No need to check destroy because the resource is a one-time action resource.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserPasswordReset_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccUserPasswordReset_basic(),
			},
		},
	})
}

func testAccUserPasswordReset_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_user_password_reset" "test" {
  instance_id  = "%[1]s"
  user_name    = "instance_not_found"
  new_password = "%[2]s"
}
`, randomId, acceptance.RandomPassword())
}

func testAccUserPasswordReset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_user" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_user_password_reset" "test" {
  instance_id  = "%[1]s"
  user_name    = huaweicloud_dms_kafka_user.test.name
  new_password = "%[3]s"

  depends_on = [huaweicloud_dms_kafka_user.test]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID,
		acceptance.RandomAccResourceName(),
		acceptance.RandomPassword(),
		acceptance.RandomPassword(),
	)
}
